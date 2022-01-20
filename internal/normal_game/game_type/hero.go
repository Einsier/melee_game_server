package game_type

import (
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/utils"
	"sync"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:存放对战中英雄有关内容,初步设计是从Hero这里加读写锁,让读写同一个英雄的go程互斥,同时外部的hero_manager加互斥锁用于增删待管理的英雄
 */

type Hero struct {
	Id            int32        //英雄在本局游戏中的id
	position      Vector2      //当前位置
	lookDirection Vector2      //看的方向
	Status        int32        //当前状态
	health        int32        //当前血量
	healthLock    sync.RWMutex //血量锁
	positionLock  sync.RWMutex //更改position的锁
	directionLock sync.RWMutex //更改direction的锁
	UpdateTime    int64        //上次更新时间
}

//NewHero 通过赋予id和位置创造一个角色
func NewHero(id int32, position Vector2) *Hero {
	return &Hero{
		Id:            id,
		position:      position,
		lookDirection: Vector2Down,
		Status:        configs.HeroAlive,
		health:        configs.HeroInitHealth,
		UpdateTime:    time.Now().UnixNano(),
	}
}

//GetPosition 获取当前英雄的position
func (h *Hero) GetPosition() Vector2 {
	h.positionLock.RLock()
	defer h.positionLock.RUnlock()
	return h.position
}

//GetDirection 获取当前英雄的朝向
func (h *Hero) GetDirection() Vector2 {
	h.directionLock.RLock()
	defer h.directionLock.Unlock()
	return h.lookDirection
}

//GetHealth 获取当前Hero的health
func (h *Hero) GetHealth() int32 {
	h.healthLock.RLock()
	defer h.healthLock.Unlock()
	return h.health
}

//ChangeHeath 英雄加血/掉血,采用互斥防止多个go程同时修改血量
func (h *Hero) ChangeHeath(heath int32) bool {
	/*
		*  @Author:sly
		*  @name:ChangeHeath
		*  @Description: 英雄加血/掉血,采用互斥防止多个go程同时修改血量
		*  @receiver h 英雄指针
		*  @param heath 变更的血量
		*  @return bool 变换之后有没有狗带,如果狗带返回true,并且将Status设置为狗带状态.
						只有得到health lock的go程才可以更改英雄状态,所以这里Status的改变不用互斥
		*  @date 2022-01-14 20:28:43
	*/
	if h.Status == configs.HeroDead {
		return true
	}
	h.healthLock.Lock()
	defer h.healthLock.Unlock()
	newHealth := utils.MaxInt32(h.health+heath, configs.HeroMaxHealth)
	if newHealth <= 0 {
		newHealth = 0
		h.Status = configs.HeroDead
		return true
	}
	return false
}

//ChangePosition 互斥的更改position
func (h *Hero) ChangePosition(position Vector2) {
	h.positionLock.Lock()
	defer h.positionLock.Unlock()
	h.position = position
}

//ChangeDirection 互斥更改direction
func (h *Hero) ChangeDirection(direction Vector2) {
	h.directionLock.Lock()
	defer h.directionLock.Unlock()
	h.lookDirection = direction
}
