package aoi

import (
	"melee_game_server/framework/entity"
	"melee_game_server/plugins/aoi/object"
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
	Speed          float64 //每ms移动多少m
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
func NewAOI(heroesInitInfo *HeroesInitInfo, mx, my, gx, gy int, updateDuration time.Duration) *AOI {
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
	aoi.Quit = make(chan *HeroQuitMsg)
	aoi.Move = make(chan *HeroMoveMsg)
	aoi.Finish = make(chan struct{})
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

//Work aoi模块开始工作,每秒12帧的发送位置信息
func (aoi *AOI) Work() {
	ticker := time.NewTicker(aoi.UpdateDuration)
	var quitMsg *HeroQuitMsg
	var moveMsg *HeroMoveMsg
	var hero *Hero
	var ok bool
	go func() {
		for {
			select {
			case <-ticker.C:
				//定时更新英雄的位置信息,并且广播给其它的玩家
				for _, hero = range aoi.Heroes {
					hero.UpdateMovement(nil)
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
				//英雄退出
				delete(aoi.Heroes, int32(*quitMsg))
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
