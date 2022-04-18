package aoi

import (
	"melee_game_server/api/client/proto"
	"melee_game_server/framework/entity"
	"melee_game_server/internal/normal_game/aoi/object"
	"melee_game_server/internal/normal_game/codec"
	"melee_game_server/internal/normal_game/game_net"
	"melee_game_server/plugins/logger"
	"time"
)

/**
*@Author Sly
*@Date 2022/4/12
*@Version 1.0
*@Description:1.0版本 暂时将地图宽度和高度,格子宽度高度设置为整形,方便进行各种取模运算
大概思想来自:https://blog.csdn.net/qq_36748278/article/details/102787225
为了效率,很多固定的运算提前的进行了运算,例如每个格子的四周格子,每个格子的邻接格子等等
*/

type MapInfo struct {
	MX    int              //地图宽度
	MY    int              //地图高度
	GX    int              //每个格子的宽度
	GY    int              //每个格子的高度
	NGX   int              //X轴上有多少个格子
	NGY   int              //Y轴上有多少个格子
	NGrid int              //一共有多少个格子
	Grids [][]*object.Grid //被上面的博客带跑偏了= =感觉还是二维数组比较好写 而且性能也一样的= =
}

type HeroInfo struct {
	HeroNum        int
	Speed          float32 //每ms移动多少m
	UpdateDuration time.Duration

	Heroes map[int32]*Hero
}

type MsgChan struct {
	Ticker *time.Ticker      //定期发包的定时器
	Quit   chan *HeroQuitMsg //英雄退出的chan
	Move   chan *HeroMoveMsg //英雄移动的chan
	Finish chan struct{}
}

type AOI struct {
	*MapInfo
	*HeroInfo
	*MsgChan
	gn *game_net.NormalGameNetServer
}

//GetGridByPos 获取当前位置应该属于哪个grid,并从中取出grid
func (mp *MapInfo) GetGridByPos(pos *entity.Vector2) *object.Grid {
	if !mp.checkPositionLegal(pos) {
		logger.Errorf("position:%+v is illegal", pos)
		return nil
	}
	return mp.Grids[int(pos.Y)/mp.GY][(int(pos.X) / mp.GX)]
}

func (mp *MapInfo) GetGridByIdx(x, y int) *object.Grid {
	if !mp.checkIdxLegal(x, y) {
		//logger.Errorf("grid[%d][%d] is illegal",y,x)
		return nil
	}
	return mp.Grids[y][x]
}

func (mp *MapInfo) checkIdxLegal(x, y int) bool {
	return !(x < 0 || y < 0 || x >= mp.NGX || y >= mp.NGY)
}

//checkPositionLegal 判断当前位置是否合法,注意例如我们的地图x最大值是40,那么40.2 or 40是违法的位置
func (mp *MapInfo) checkPositionLegal(pos *entity.Vector2) bool {
	return !(pos.X < 0 || int(pos.X) >= mp.MX || pos.Y < 0 || int(pos.Y) >= mp.MY)
}

//NewAOI 创建一个AOI模块
func NewAOI(heroesInitInfo *HeroesInitInfo, mx, my, gx, gy int, updateDuration time.Duration, gn *game_net.NormalGameNetServer) *AOI {
	if mx < 0 || my < 0 || gx < 0 || gy < 0 {
		return nil
	}

	aoi := new(AOI)
	//初始化map有关字段
	mapInfo := &MapInfo{
		MX:    mx,
		MY:    my,
		GX:    gx,
		GY:    gy,
		NGX:   mx/gx + 1, //例如mx为10,gx为6,那么x方向应有两个格子
		NGY:   my/gy + 1,
		NGrid: (mx/gx + 1) * (my/gy + 1),
	}

	logger.Infof("aoi完成地图的初始化:%+v", *mapInfo)
	aoi.MapInfo = mapInfo
	//初始化grids,创建NGrid个数的格子
	id := 0
	aoi.Grids = make([][]*object.Grid, mapInfo.NGY)
	for i := 0; i < mapInfo.NGY; i++ {
		aoi.Grids[i] = make([]*object.Grid, mapInfo.NGX)
	}
	for y := 0; y < mapInfo.NGY; y++ {
		for x := 0; x < mapInfo.NGX; x++ {
			aoi.Grids[y][x] = object.NewGrid(id, x, y, aoi.Grids)
		}
	}

	//初始化英雄
	heroInfo := &HeroInfo{
		HeroNum: len(heroesInitInfo.heroes),
		Speed:   heroesInitInfo.Speed,
		Heroes:  make(map[int32]*Hero, len(heroesInitInfo.heroes)),
	}
	for i := 0; i < heroInfo.HeroNum; i++ {
		singleInfo := heroesInitInfo.heroes[i]
		if hero := NewHero(singleInfo.Id, singleInfo.Position, singleInfo.Direction, heroInfo.Speed, aoi); hero == nil {
			logger.Errorf("向aoi插入英雄的时候出现错误,请检查英雄是否在地图边界处/超出地图位置.hero position:%+v,id:%d", singleInfo.Position, singleInfo.Id)
		} else {
			heroInfo.Heroes[hero.Id] = hero
		}
	}
	heroInfo.UpdateDuration = updateDuration
	aoi.HeroInfo = heroInfo

	aoi.MsgChan = new(MsgChan)
	aoi.Quit = make(chan *HeroQuitMsg, 512)
	aoi.Move = make(chan *HeroMoveMsg, 65536)
	aoi.Finish = make(chan struct{})
	aoi.gn = gn
	return aoi
}

func (aoi *AOI) UpdateHeroPosition(info *HeroMoveMsg) {
	hero := aoi.Heroes[info.Id]

	if hero == nil {
		logger.Errorf("aoi收到了不存在的英雄位置更新信息:hero num:%d", info.Id)
		return
	}
	hero.UpdateMovement(info)
}

//Work aoi模块开始工作,每秒12帧的发送位置信息,注意一定是初始化了所有玩家的net.Conn才可以执行
func (aoi *AOI) Work() {
	ticker := time.NewTicker(aoi.UpdateDuration)
	go func() {
		var quitMsg *HeroQuitMsg
		var moveMsg *HeroMoveMsg
		var hero *Hero
		var ok bool
		for {
			select {
			case <-ticker.C:
				//因为发送的时候不会改变每个玩家的位置,所以拷贝一手全体玩家的当前位置,用于发送
				m := make(map[int32]*proto.HeroMovementChangeBroadcast, len(aoi.Heroes))
				for _, hero = range aoi.Heroes {
					//定时更新英雄的位置信息,更新hero的位置信息
					hero.UpdateMovement(nil)
					//把当前英雄的位置信息放到m中
					m[hero.Id] = &proto.HeroMovementChangeBroadcast{
						HeroId:           hero.Id,
						HeroMovementType: entity.V2toToHeroMovementType[hero.direction],
						HeroPosition:     hero.position.ToProto(),
						Time:             hero.updateTime.UnixMilli(),
					}
				}
				for _, me := range aoi.Heroes {
					//将当前hero视野中的全部英雄的topMsg的指针放到view中,把view传给网络模块进行发送
					meMap := make(map[int32]*proto.HeroMovementChangeBroadcast)
					for otherId := range hero.View {
						if meMap[otherId] = m[otherId]; meMap[otherId] == nil {
							panic("!!!!!")
						}
					}
					//logger.Infof("hero:%d 视野中的玩家有:%v",hero.Id,view)

					aoi.gn.SendByHeroId([]int32{me.Id}, codec.EncodeUnicast(&proto.HeroFrameSyncUnicast{Movement: meMap}))
				}
			case moveMsg = <-aoi.Move:
				hero, ok = aoi.Heroes[moveMsg.Id]
				if !ok {
					logger.Errorf("收到了已经退出/不存在的heroId", moveMsg.Id)
					continue
				} else {
					hero.UpdateMovement(moveMsg)
				}
			case quitMsg = <-aoi.Quit:
				quitHeroId := quitMsg.id
				//英雄退出,从map中删除英雄,并且从能看到被删玩家的玩家的View中删除被删玩家
				for otherHeroId := range aoi.Heroes[quitHeroId].View {
					delete(aoi.Heroes[otherHeroId].View, quitHeroId)
				}
				delete(aoi.Heroes, quitHeroId)
			case <-aoi.Finish:
				//结束
				logger.Infof("aoi模块结束工作")
				return
			}
		}
	}()
	aoi.Ticker = ticker
}

//Stop 停止aoi模块的运作
func (aoi *AOI) Stop() {
	aoi.Ticker.Stop()
	close(aoi.Finish)
}

func (aoi *AOI) PutMove(msg *HeroMoveMsg) {
	aoi.Move <- msg
}

func (aoi *AOI) RemoveHero(id int32) {
	aoi.Quit <- &HeroQuitMsg{id: id}
}
