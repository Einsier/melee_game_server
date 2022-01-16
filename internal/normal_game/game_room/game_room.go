package game_room

import (
	gt "melee_game_server/internal/normal_game/game_type"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:用于管理一次游戏的对局,包括玩家管理,对局中的prop(道具)管理以及子弹管理
 */

type GameRoom struct {
	heroes map[int32]*gt.Hero
	lock   sync.Mutex //todo:优化go程互斥访问map
}

func NewGameRoomManager() *GameRoom {
	return &GameRoom{
		heroes: make(map[int32]*gt.Hero),
	}
}

func (grm *GameRoom) AddHero(h *gt.Hero) {
	grm.lock.Lock()
	defer grm.lock.Unlock()
	grm.heroes[h.Id] = h
}

func (grm *GameRoom) DeleteHero(h *gt.Hero) {
	delete(grm.heroes, h.Id)
}

func (grm *GameRoom) GetHero(id int32) *gt.Hero {
	return grm.heroes[id]
}
