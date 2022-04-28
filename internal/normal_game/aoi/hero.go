package aoi

import (
	"melee_game_server/api/client/proto"
	"melee_game_server/framework/entity"
	"melee_game_server/internal/normal_game/aoi/collision"
	"melee_game_server/internal/normal_game/codec"
	"melee_game_server/internal/normal_game/game_net"
	"melee_game_server/plugins/logger"
	"strconv"
	"time"
)

/**
*@Author Sly
*@Date 2022/4/12
*@Version 1.0
*@Description:规定Hero的mu锁比Grid的锁更顶层,所以例如一个英雄的位置跨越了格子,想要更改一个英雄的所在的Grid,那么应该
Hero mu lock->更改英雄位置->英雄离开的Grid mu lock->删除英雄->英雄离开的Grid mu unlock->英雄进入的Grid mu lock->
将英雄加入格子 ->英雄进入的Grid mu unlock ->Hero mu lock
这样例如有一个英雄在两个格子的边缘来回走动,那么因为给英雄加了锁再给格子加锁,所以对格子的访问是互斥的,所以同一个英雄肯定同一时刻出现
且仅出现在一个格子中.
*/

type Hero struct {
	Id         int32 //英雄id
	Name       string
	direction  entity.Vector2 //当前运动状态,是Vector2Up,Down,Left,Right,Zero的枚举
	position   entity.Vector2 //英雄当前位置,也就是Ruby的中心的位置
	updateTime time.Time      //上次位置更新时间
	at         *Grid          //当前处于哪个格子
	speed      float32
	aoi        *AOI
	View       map[int32]struct{} //能看到的英雄
}

func (h *Hero) VisibleHeroes() []int32 {
	return nil
}

//NewHero 初始化英雄,需要初始化英雄当前所在的格子,所以需要将aoi传入
func NewHero(id int32, position, direction entity.Vector2, speed float32, aoi *AOI) *Hero {
	var grid *Grid
	if grid = aoi.GetGridByPos(&position); grid == nil {
		return nil
	} else {
		if _, ok := grid.Objs[id]; ok {
			//如果已经有相同id的英雄了，报错
			panic(ok)
		} else {
			//在格子中插入本英雄
			grid.Objs[id] = struct{}{}
		}
	}

	return &Hero{
		Id:         id,
		Name:       "Hero" + strconv.Itoa(int(id)),
		position:   position,
		direction:  direction,
		speed:      speed,
		at:         grid,
		aoi:        aoi,
		View:       map[int32]struct{}{},
		updateTime: time.Now(),
	}
}

//UpdateMovement 更改玩家位置的唯一方式,如果传入的info不为nil,按照给定的info更新玩家位置,如果传入的info为nil,那么按照上次更新的时间更新
func (h *Hero) UpdateMovement(info *HeroMoveMsg, gn *game_net.NormalGameNetServer) {
	if info == nil {
		//当前是定时更新,自己计算更新的info
		info = new(HeroMoveMsg)
		info.Direction = h.direction
		info.Time = time.Now()
		timeSpace := info.Time.Sub(h.updateTime).Milliseconds()
		info.Position = entity.Vector2{
			//运算逻辑:
			//新位置X= 旧位置X	+ 更新间隔(单位ms)	*英雄移动速度 * 英雄移动方向的X(取值为1/0/-1)
			X: h.position.X + float32(timeSpace)*h.speed*info.Direction.X,
			Y: h.position.Y + float32(timeSpace)*h.speed*info.Direction.Y,
		}
	} else if info.Position == entity.Vector2Zero {
		//todo:测试用,只改变方向,让服务器计算当前的位置(本来该由前端做,服务器校验)
		info.Position = entity.Vector2{
			//运算逻辑:
			//新位置X= 旧位置X	+ 更新间隔(单位ms)	*英雄移动速度 * 英雄移动方向的X(取值为1/0/-1)
			X: h.position.X + float32(info.Time.Sub(h.updateTime).Milliseconds())*h.speed*h.direction.X,
			Y: h.position.Y + float32(info.Time.Sub(h.updateTime).Milliseconds())*h.speed*h.direction.Y,
		}
	}
	//现在info中存放的是希望走的位置,h.position是英雄原来的位置,应该做碰撞校验,判断英雄到现在的位置是否合法,如果合法,那么进行到下一步,更新九宫格
	//如果不合法,那么只是更新updateTime和,不改变h.position,这样相当于将英雄退回到上一帧的位置.
	if h.aoi.qt.CheckCollision(collision.NewRubyCollisionCheckRectangle(h.Name, &h.position, &info.Position)) {
		//发生碰撞让玩家停下
		h.direction = entity.Vector2Zero
		h.updateTime = info.Time
		return
	}

	//更新位置
	if to := h.aoi.GetGridByPos(&info.Position); to != nil {
		//如果更新的位置是正确的,那么更新位置,否则判断玩家位置不合法,让玩家位于上次更新的位置
		var joinGrid, leaveGrid []*Grid
		refreshAll := func() {
			leaveGrid = []*Grid{
				h.aoi.GetGridByIdx(h.at.XIdx-1, h.at.YIdx+1),
				h.aoi.GetGridByIdx(h.at.XIdx-1, h.at.YIdx),
				h.aoi.GetGridByIdx(h.at.XIdx-1, h.at.YIdx-1),
				h.aoi.GetGridByIdx(h.at.XIdx+1, h.at.YIdx+1),
				h.aoi.GetGridByIdx(h.at.XIdx+1, h.at.YIdx),
				h.aoi.GetGridByIdx(h.at.XIdx+1, h.at.YIdx-1),
				h.aoi.GetGridByIdx(h.at.XIdx, h.at.YIdx+1),
				h.aoi.GetGridByIdx(h.at.XIdx, h.at.YIdx-1),
			}
			joinGrid = []*Grid{
				h.aoi.GetGridByIdx(to.XIdx-1, to.YIdx+1),
				h.aoi.GetGridByIdx(to.XIdx-1, to.YIdx),
				h.aoi.GetGridByIdx(to.XIdx-1, to.YIdx-1),
				h.aoi.GetGridByIdx(to.XIdx+1, to.YIdx+1),
				h.aoi.GetGridByIdx(to.XIdx+1, to.YIdx),
				h.aoi.GetGridByIdx(to.XIdx+1, to.YIdx-1),
				h.aoi.GetGridByIdx(to.XIdx, to.YIdx+1),
				h.aoi.GetGridByIdx(to.XIdx, to.YIdx-1),
			}
		}

		if h.at != to {
			//如果离开了at,那么应该:
			//1.把玩家从at中删除
			//2.把玩家加入到to中
			logger.Infof("[hero:%d]position:%+v->%+v,grid:[%d][%d]->[%d][%d]", h.Id, &h.position, &info.Position, h.at.YIdx, h.at.XIdx, to.YIdx, to.XIdx)
			delete(h.at.Objs, h.Id)
			to.Objs[h.Id] = struct{}{}
			//3.判断to和at的关系,有可能有以下情况:
			//感觉用switch底层效率太低了...还是用if吧= =.注意能看见玩家的玩家 = 玩家能看见的玩家,所以双向更新一手
			if to.YIdx-h.at.YIdx == 0 {
				//如果Y没有变化,说明玩家左右移动
				if to.XIdx-h.at.XIdx == 1 {
					//玩家向右移动:	to在at的右边->删除at所在的格子的左边三个格子中的英雄的View中的本英雄id+向to所在的格子的右边三个格子中的英雄的View添加本英雄id
					leaveGrid = []*Grid{
						h.aoi.GetGridByIdx(h.at.XIdx-1, h.at.YIdx+1),
						h.aoi.GetGridByIdx(h.at.XIdx-1, h.at.YIdx),
						h.aoi.GetGridByIdx(h.at.XIdx-1, h.at.YIdx-1),
					}
					joinGrid = []*Grid{
						h.aoi.GetGridByIdx(to.XIdx+1, to.YIdx+1),
						h.aoi.GetGridByIdx(to.XIdx+1, to.YIdx),
						h.aoi.GetGridByIdx(to.XIdx+1, to.YIdx-1),
					}
				} else if to.XIdx-h.at.XIdx == -1 {
					//玩家向左移动:	to在at的左边->删除at所在的格子的右边三个格子中的英雄的View中的本英雄id+向to所在的格子的左边三个格子中的英雄的View添加本英雄id
					leaveGrid = []*Grid{
						h.aoi.GetGridByIdx(h.at.XIdx+1, h.at.YIdx+1),
						h.aoi.GetGridByIdx(h.at.XIdx+1, h.at.YIdx),
						h.aoi.GetGridByIdx(h.at.XIdx+1, h.at.YIdx-1),
					}
					joinGrid = []*Grid{
						h.aoi.GetGridByIdx(to.XIdx-1, to.YIdx+1),
						h.aoi.GetGridByIdx(to.XIdx-1, to.YIdx),
						h.aoi.GetGridByIdx(to.XIdx-1, to.YIdx-1),
					}
				} else {
					//其他情况	:	可能玩家网络不好等原因,一下跨越了很多个格子->删除at所在的周围八个格子中的英雄的View中的本英雄id+向to所在的格子的周围八个格子中的英雄的View添加本英雄id
					refreshAll()
				}
			} else if to.XIdx-h.at.XIdx == 0 {
				//如果X没有变化,说明玩家上下移动
				if to.YIdx-h.at.YIdx == 1 {
					//玩家向上移动:	to在at的上边->删除at所在的格子的下边三个格子中的英雄的View中的本英雄id+向to所在的格子的上边三个格子中的英雄的View添加本英雄id
					leaveGrid = []*Grid{
						h.aoi.GetGridByIdx(h.at.XIdx-1, h.at.YIdx-1),
						h.aoi.GetGridByIdx(h.at.XIdx+1, h.at.YIdx-1),
						h.aoi.GetGridByIdx(h.at.XIdx, h.at.YIdx-1),
					}
					joinGrid = []*Grid{
						h.aoi.GetGridByIdx(to.XIdx+1, to.YIdx+1),
						h.aoi.GetGridByIdx(to.XIdx-1, to.YIdx+1),
						h.aoi.GetGridByIdx(to.XIdx, to.YIdx+1),
					}
				} else if to.YIdx-h.at.YIdx == -1 {
					//玩家向下移动:	to在at的下边->删除at所在的格子的上边三个格子中的英雄的View中的本英雄id+向to所在的格子的下边三个格子中的英雄的View添加本英雄id
					leaveGrid = []*Grid{
						h.aoi.GetGridByIdx(h.at.XIdx+1, h.at.YIdx+1),
						h.aoi.GetGridByIdx(h.at.XIdx, h.at.YIdx+1),
						h.aoi.GetGridByIdx(h.at.XIdx-1, h.at.YIdx+1),
					}
					joinGrid = []*Grid{
						h.aoi.GetGridByIdx(to.XIdx-1, to.YIdx-1),
						h.aoi.GetGridByIdx(to.XIdx, to.YIdx-1),
						h.aoi.GetGridByIdx(to.XIdx+1, to.YIdx-1),
					}
				} else {
					//其他情况	:	可能玩家网络不好等原因,一下跨越了很多个格子->删除at所在的周围八个格子中的英雄的View中的本英雄id+向to所在的格子的周围八个格子中的英雄的View添加本英雄id
					refreshAll()
				}
			} else {
				//to和at的X和Y都不相等
				//其他情况	:	可能玩家网络不好等原因,一下跨越了很多个格子->删除at所在的周围八个格子中的英雄的View中的本英雄id+向to所在的格子的周围八个格子中的英雄的View添加本英雄id
				refreshAll()
			}

			//将玩家的id从leave中删除,加入到join中
			for i := 0; i < len(leaveGrid); i++ {
				if leaveGrid[i] != nil {
					//这里需要判断是否越界
					for id := range leaveGrid[i].Objs {
						//将本英雄从leaveGrid中的英雄的可见英雄集合中删除
						delete(h.aoi.Heroes[id].View, h.Id)
						//将leaveGrid中的英雄从本英雄的可见英雄集合中删除
						delete(h.View, id)

						//分别给自己和其他英雄发送离开视野的报文,让前端取消加载英雄
						msg1 := codec.EncodeUnicast(&proto.HeroLeaveSightUnicast{HeroId: h.Id})
						msg2 := codec.EncodeUnicast(&proto.HeroLeaveSightUnicast{HeroId: id})
						gn.SendByHeroId([]int32{id}, msg1)
						gn.SendByHeroId([]int32{h.Id}, msg2)
						logger.Infof("[%d][%d]离开了彼此的视野", id, h.Id)
					}
				}
				if joinGrid[i] != nil {
					//这里需要判断是否越界
					for id := range joinGrid[i].Objs {
						if id != h.Id {
							//将本英雄加入到joinGrid中的英雄的可见英雄集合中
							//bug fix:让自己不进入自己的可见英雄之内
							h.aoi.Heroes[id].View[h.Id] = struct{}{}
							//将joinGrid中的英雄的可见英雄中加入本英雄
							h.View[id] = struct{}{}
							logger.Infof("[%d][%d]进入了彼此的视野", id, h.Id)
						}
					}
				}
			}
		}
		h.direction = info.Direction
		h.position = info.Position //更新玩家位置
		h.at = to
	} else {
		//如果越界,那么同样将玩家的direction改成zero
		h.direction = entity.Vector2Zero
	}
	h.updateTime = info.Time
}
