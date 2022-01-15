package game_type

import (
	config "melee_game_server/configs/game_type_config"
	"melee_game_server/utils"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:存放对战中英雄有关内容
 */

type Hero struct {
	Id            int32        //英雄在本局游戏中的id
	Position      Vector2      //当前位置
	LookDirection Vector2      //看的方向
	Status        int32        //当前状态
	Health        int32        //当前血量
	healthLock    sync.RWMutex //血量锁
	positionLock  sync.RWMutex //更改position的锁
}

//NewHero 通过赋予id和位置创造一个角色
func NewHero(id int32, position Vector2) *Hero {
	return &Hero{
		Id:            id,
		Position:      position,
		LookDirection: Vector2Down,
		Status:        config.HeroAlive,
		Health:        config.HeroInitHealth,
	}
}

func (h *Hero) ChangeHeath(heath int32) bool {
	/*
	*  @Author:sly
	*  @name:ChangeHeath
	*  @Description: 英雄加血/掉血,采用互斥防止多个go程同时修改血量
	*  @receiver h 英雄指针
	*  @param heath 变更的血量
	*  @return bool 变换之后有没有狗带,如果狗带返回true,并且将Status设置为狗带状态
	*  @date 2022-01-14 20:28:43
	 */
	if h.Status == config.HeroDead {
		return true
	}
	h.healthLock.Lock()
	defer h.healthLock.Unlock()
	newHealth := utils.MaxInt32(h.Health+heath, config.HeroMaxHealth)
	if newHealth <= 0 {
		newHealth = 0
		h.Status = config.HeroDead
		return true
	}
	return false
}

//ChangePosition 互斥的更改position
func (h *Hero) ChangePosition(position Vector2) {
	h.positionLock.Lock()
	defer h.positionLock.Unlock()
	h.Position = position
}
