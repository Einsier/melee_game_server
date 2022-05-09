package game_type

import (
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/entity"
	"melee_game_server/plugins/logger"
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
	Id           int32          //英雄在本局游戏中的id
	position     entity.Vector2 //当前位置
	moveStatus   entity.Vector2 //运动状态,相当于枚举,静止->Vector2Zero,向上->Vector2Up,向下->Vector2Down,向左->Vector2Left,向右->Vector2Right
	status       int32          //当前状态,枚举类型见normal_game_type_configs/hero_configs
	Health       int32          //当前血量
	updateTime   int64          //发出当前动作的前端发过来的时间
	healthLock   sync.RWMutex   //血量锁
	movementLock sync.RWMutex   //更改moveType以及position的锁
	statusLock   sync.RWMutex   //更改status的锁
}

//NewHero 通过赋予id和位置创造一个角色
func NewHero(id int32, position entity.Vector2) *Hero {
	return &Hero{
		Id:         id,
		position:   position,
		moveStatus: entity.Vector2Zero,
		status:     configs.HeroAlive,
		Health:     configs.HeroInitHealth,
		updateTime: time.Now().UnixNano(),
	}
}

//GetPosition 获取当前英雄的position
func (h *Hero) GetPosition() entity.Vector2 {
	h.movementLock.RLock()
	defer h.movementLock.RUnlock()

	ct := time.Now().UnixNano()
	//运算逻辑: hero的客户端上次向服务器汇报的位置 + (当前运动状态(例如向上运动为(0,1)) * 上次汇报距离现在过了多长时间 * 英雄运动速度)
	p := h.position.Add(h.moveStatus.MultiplyScalar(float32(h.updateTime-ct) * configs.HeroMoveSpeed))
	//utils.TransNaN(&p.X)
	//utils.TransNaN(&p.Y)
	return p
}

//GetMoveStatus 获取当前英雄的运动状态
func (h *Hero) GetMoveStatus() entity.Vector2 {
	h.movementLock.RLock()
	defer h.movementLock.RUnlock()
	return h.moveStatus
}

//GetHealth 获取当前Hero的health
func (h *Hero) GetHealth() int32 {
	h.healthLock.RLock()
	defer h.healthLock.RUnlock()
	return h.Health
}

//ChangeHeath 英雄加血/掉血,采用互斥防止多个go程同时修改血量
func (h *Hero) ChangeHeath(heath int32) (isChange, isDead bool, newHealth int32) {
	/*
		*  @Author:sly
		*  @name:ChangeHeath
		*  @Description: 分为以下几种情况:
		1.Hero在调用本函数之前就已经狗带了->isChange = false,isDead = true -> 调用本函数的函数应该只是简单的什么也不做
		2.Hero在调用本次函数之后狗带了->isChange = true,isDead = true -> 调用本函数的函数应该负责处理英雄的死亡事件
		3.Hero在调用本函数之后没有狗带->isChange = true/false,isDead = false,newHealth根据isChange的结果,可能为新血量,也可能没有改变
		*  @receiver h 待改变的英雄
		*  @param heath	为正表示加血,为负表示掉血
		*  @return isChange	是否和原来的血量不同,如果不同需要进行广播
		*  @return isDead	是否因为本次掉血而狗带
		*  @return newHealth 返回新的血量
		*  @date 2022-01-22 11:35:01
	*/

	if h.status == configs.HeroDead {
		logger.Errorf("[ChangeHeath]收到了更新已经狗带的英雄的血量的请求!heroId:&d", h.Id)
		return false, true, 0
	}
	h.healthLock.Lock()
	defer h.healthLock.Unlock()
	oldHealth := h.Health
	newHealth = utils.MinInt32(h.Health+heath, configs.HeroMaxHealth)
	if oldHealth != newHealth {
		isChange = true
	}
	if newHealth <= 0 {
		newHealth = 0
		h.status = configs.HeroDead
		isDead = true
		return
	}
	h.Health = newHealth
	return
}

func (h *Hero) SetStatus(s int32) {
	h.statusLock.Lock()
	defer h.statusLock.Unlock()
	h.status = s
}

func (h *Hero) GetStatus() int32 {
	h.statusLock.RLock()
	defer h.statusLock.RUnlock()
	return h.status
}

//SetPositionInfo 互斥的更改当前的位置position,当前的运动状态moveStatus,更新时间t,注意这个时间应该是客户端生成的
func (h *Hero) SetPositionInfo(position entity.Vector2, moveStatus entity.Vector2, t int64) {
	h.movementLock.Lock()
	defer h.movementLock.Unlock()
	h.position = position
	h.moveStatus = moveStatus
	h.updateTime = t
}

//SetMoveStatus 互斥更改moveType
func (h *Hero) SetMoveStatus(moveStatus entity.Vector2) {
	h.movementLock.Lock()
	defer h.movementLock.Unlock()
	h.moveStatus = moveStatus
}
