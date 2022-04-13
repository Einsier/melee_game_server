package aoi

import (
	"melee_game_server/framework/entity"
	"melee_game_server/plugins/aoi/object"
	"sync"
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
	Id			int32					//英雄id

	mu			sync.RWMutex			//保护下面的变量
	direction	entity.Vector2			//当前运动状态,是Vector2Up,Down,Left,Right,Zero的枚举
	position	entity.Vector2			//英雄当前位置
	updateTime	time.Time				//上次位置更新时间
	grid		*object.Grid			//当前处于哪个网格
	speed		float64
	mapInfo		*MapInfo
}

//HeroMovementInfo 表示英雄的当前状态
type HeroMovementInfo struct {
	Id			int32
	Position	entity.Vector2		//当前位置
	Direction	entity.Vector2		//面朝方向
	Time		time.Time			//发生的时间
}

//HeroesInitInfo 用于初始化
type HeroesInitInfo struct {
	Speed	float64
	heroes	[]*HeroMovementInfo
}



func (h *Hero)VisibleHeroes()[]int32{
	h.mu.RLock()
	grid := h.grid
	h.mu.RUnlock()

	//如果英雄在本函数执行期间改变了位置,由于九宫格算法对于英雄的当前位置有容错的空间,并且给很多格子加锁容易死锁,所以不加锁
	ret := make([]int32,0)
	for i := 0; i < 8; i++ {
		ret = append(ret, grid.Surround[i].GetObjs()...)
	}
	return ret
}

//NewHero 初始化英雄,需要初始化英雄当前所在的格子,所以需要将aoi传入
func NewHero(id int32,position,direction entity.Vector2,speed float64,aoi *AOI)*Hero{
	return &Hero{
		Id:       id,
		mu:       sync.RWMutex{},
		position: position,
		direction: direction,
		speed: 	  speed,
		grid:     aoi.GetGridByPos(position),
		mapInfo:  aoi.MapInfo,
	}
}

//UpdateMovement 更改玩家位置的唯一方式,如果传入的info不为nil,按照给定的info更新玩家位置,如果传入的info为nil,那么按照上次更新的时间更新
func (h *Hero)UpdateMovement(info *HeroMovementInfo){
	h.mu.Lock()
	defer h.mu.Unlock()
	if info == nil{
		//当前是定时更新,自己计算更新的info
		info = new(HeroMovementInfo)
		info.Direction = h.direction
		info.Time = time.Now()
		timeSpace := info.Time.Sub(h.updateTime).Milliseconds()
		info.Position = entity.Vector2{
		//运算逻辑:
		//新位置X= 旧位置X	+ 更新间隔(单位ms)	*英雄移动速度 * 英雄移动方向的X(取值为1/0/-1)
			X: h.position.X + float64(timeSpace)*h.speed*info.Direction.X,
			Y: h.position.Y + float64(timeSpace)*h.speed*info.Direction.Y,
		}
	}

	//更新位置
	to := h.mapInfo.GetGridByPos(info.Position)
	if h.grid != to{
		h.grid.ObjLeave(h.Id,to)
	}
}
