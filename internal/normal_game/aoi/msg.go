package aoi

import (
	"melee_game_server/framework/entity"
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

//SingleHeroInitInfo 用于初始化aoi模块中的每个英雄
type SingleHeroInitInfo struct {
	*HeroMoveMsg
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
