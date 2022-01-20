package game_room

import (
	gt "melee_game_server/internal/normal_game/game_type"
	"sync"
	"sync/atomic"
)

/**
*@Author chenjiajia
*@Date 2022/1/19
*@Version 1.0
*@Description: 用于管理游戏中的 player
 */

type PlayersManager struct {
	players       map[int32]*gt.Player
	currentHeroId int32        //当前可分配的heroId
	lock          sync.RWMutex //用于对 Players增删加锁
}

func NewPlayersManager() *PlayersManager {
	return &PlayersManager{
		players:       make(map[int32]*gt.Player),
		currentHeroId: 0,
		lock:          sync.RWMutex{},
	}
}

//增加一个玩家
func (pm *PlayersManager) AddPlayer(p *gt.Player) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.players[p.Id] = p
}

//删除一个玩家
func (pm *PlayersManager) DeletePlayer(p *gt.Player) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	delete(pm.players, p.Id)
}

//获得一个玩家
func (pm *PlayersManager) GetPlayer(id int32) *gt.Player {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	return pm.players[id]
}

//分配heroId
func (pm *PlayersManager) CreateHeroId() int32 {
	return atomic.AddInt32(&pm.currentHeroId, 1)
}
