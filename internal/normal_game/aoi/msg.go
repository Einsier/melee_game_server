package aoi

import (
	"melee_game_server/api/client/proto"
	"melee_game_server/framework/entity"
	"melee_game_server/internal/normal_game/aoi/collision"
	"melee_game_server/utils"
	"strconv"
	"time"
)

/**
*@Author Sly
*@Date 2022/4/15
*@Version 1.0
*@Description:
 */

type HeroQuitMsg struct {
	id int32
}

//MapResourcesInfo 地图资源
type MapResourcesInfo []*collision.Rectangle

//SingleHeroInitInfo 用于初始化aoi模块中的每个英雄
type SingleHeroInitInfo struct {
	*HeroMoveMsg
}

type BulletLaunchMsg struct {
	HeroId    int32
	Position  entity.Vector2
	Direction entity.Vector2
	Time      time.Time
}

//HeroMoveMsg 表示英雄的当前状态
type HeroMoveMsg struct {
	Id        int32
	Position  entity.Vector2 //当前位置
	Direction entity.Vector2 //面朝方向
	Time      time.Time      //发生的时间
}

//HeroesInitInfo 用于初始化
type HeroesInitInfo struct {
	Speed  float32 //每ms走多少
	heroes []*HeroMoveMsg
}

//NewRandomHeroesInitInfo 根据传入的地图(使用qt表示传入的地图),随机的生成heroNum个不和地图障碍物重合的HeroMoveMsg.
func NewRandomHeroesInitInfo(heroNum int32, speed float32, qt *collision.Quadtree) *HeroesInitInfo {
	heroes := make([]*HeroMoveMsg, heroNum)
	var position entity.Vector2
	//英雄之间彼此也不应该重合,所以使用一个四叉树避免重合
	heroQt := collision.NewQuadtree(collision.NewRectangle("heroQt", entity.NewVector2(0, 0), qt.Self.Width, qt.Self.Height), 1)

	for i := 0; i < len(heroes); i++ {
		heroes[i] = new(HeroMoveMsg)
		heroes[i].Id = int32(i + 1) //normal game中所有玩家的id应该是1~heroNum
		heroes[i].Direction = entity.Vector2Zero
		heroes[i].Time = time.Now()
		for {
			position.X = utils.RandomFloat32(1, qt.Self.UR.X-1)
			position.Y = utils.RandomFloat32(1, qt.Self.UR.Y-1)
			heroMid := collision.NewRubyRectangleByMid(position, "insert-"+strconv.Itoa(i+1))
			if !qt.CheckCollision(heroMid) && !heroQt.CheckCollision(heroMid) {
				//如果和当前的地图障碍物或者其他英雄重合,则重新随机X和Y
				heroQt.Insert(heroMid)
				break
			}
		}
		heroes[i].Position = position
	}
	return &HeroesInitInfo{
		Speed:  speed,
		heroes: heroes,
	}
}

func NewHeroMoveMsgFromProto(req *proto.HeroMovementChangeRequest) *HeroMoveMsg {
	return &HeroMoveMsg{
		Id:        req.HeroId,
		Position:  entity.NewVector2(req.Position.X, req.Position.Y),
		Direction: entity.HeroMovementTypeToV2[req.HeroMovementType],
		Time:      time.UnixMilli(req.Time),
	}
}
