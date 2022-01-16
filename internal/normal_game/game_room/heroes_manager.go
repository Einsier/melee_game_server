package game_room

import (
	gt "melee_game_server/internal/normal_game/game_type"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:用于管理游戏中的hero
 */

type HeroesManager struct {
	heroes map[int32]*gt.Hero
	lock   sync.RWMutex //用于对heroes增删加锁 todo:优化go程互斥访问map
}

func NewHeroesManager() *HeroesManager {
	return &HeroesManager{
		heroes: make(map[int32]*gt.Hero),
		lock:   sync.RWMutex{},
	}
}

//AddHero 向heroes中增加一个英雄
func (hm *HeroesManager) AddHero(h *gt.Hero) {
	hm.lock.Lock()
	defer hm.lock.Unlock()
	hm.heroes[h.Id] = h
}

//GetHeroPosition 获取heroId的英雄的当前位置
func (hm *HeroesManager) GetHeroPosition(heroId int32) gt.Vector2 {
	hm.lock.RLock()
	defer hm.lock.RUnlock()
	return hm.heroes[heroId].GetPosition()
}

//DeleteHero 从heroes中删除一个英雄
func (hm *HeroesManager) DeleteHero(h *gt.Hero) {
	hm.lock.Lock()
	defer hm.lock.Unlock()
	delete(hm.heroes, h.Id)
}

//GetHero 从heroes中获得一个英雄
func (hm *HeroesManager) GetHero(id int32) *gt.Hero {
	hm.lock.RLock()
	defer hm.lock.RUnlock()
	return hm.GetHero(id)
}

//MoveHeroPosition 更改hero的position
func (hm *HeroesManager) MoveHeroPosition(heroId int32, position gt.Vector2) {
	hm.lock.RLock()
	defer hm.lock.RUnlock()
	hm.heroes[heroId].ChangePosition(position)
}
