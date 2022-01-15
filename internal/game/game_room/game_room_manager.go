package game_room

import (
	gt "melee_game_server/internal/game/game_type"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:用于管理一次游戏的对局
 */

type GameRoomManager struct {
	heroes map[int32]*gt.Hero
	lock   sync.Mutex //todo:优化go程互斥访问map
}

func NewGameRoomManager() *GameRoomManager {
	return &GameRoomManager{
		heroes: make(map[int32]*gt.Hero),
	}
}

func (grm *GameRoomManager) AddHero(h *gt.Hero) {
	grm.lock.Lock()
	defer grm.lock.Unlock()
	grm.heroes[h.Id] = h
}

func (grm *GameRoomManager) DeleteHero(h *gt.Hero) {
	delete(grm.heroes, h.Id)
}

func (grm *GameRoomManager) GetHero(id int32) *gt.Hero {
	return grm.heroes[id]
}
