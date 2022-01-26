package game_room

import (
	configs "melee_game_server/configs/normal_game_type_configs"
	gt "melee_game_server/internal/normal_game/game_type"
	"sync"
)

/**
*@Author chenjiajia
*@Date 2022/1/19
*@Version 1.0
*@Description: 用于管理游戏中的 player
 */

type PlayersManager struct {
	players      map[int32]*gt.Player //key为玩家id
	pToH         map[int32]int32      //key为int32的玩家id,value为int32的英雄id
	hToP         map[int32]int32      //key为int32的英雄id,value为int32的玩家id
	lock         sync.RWMutex         //用于对Players增删加锁
	RegisterLock sync.RWMutex         //用于PlayerEnterGameRequestCallback中加的锁
	LeaveLock    sync.RWMutex         //用于PlayerQuitGameRequestCallback中加的锁
}

//GetHeroId 通过PlayerId查找hid
func (pm *PlayersManager) GetHeroId(pid int32) int32 {
	pm.RegisterLock.RLocker()
	defer pm.RegisterLock.Unlock()
	return pm.pToH[pid]
}

//GetPlayerId 通过HeroId查找pid
func (pm *PlayersManager) GetPlayerId(hid int32) int32 {
	pm.RegisterLock.RLocker()
	defer pm.RegisterLock.Unlock()
	return pm.hToP[hid]
}

//IsPlayerRegistered 检查Player有没有注册过,返回false表示没注册过,应该注册
func (pm *PlayersManager) IsPlayerRegistered(pId int32) bool {
	p, ok := pm.players[pId]
	if !ok {
		//如果这个PlayerId不是本局应该加入游戏的PlayerId,返回true,以避免非法的用户的注册
		return true
	}
	if p.Status == configs.PlayerNotRegisteredStatus {
		//让合法并且没有注册的用户可以注册,返回false
		return false
	}
	//如果用户合法,但是已经注册过(注册过p.Conn不为nil),返回false,避免合法用户的重复注册
	return true
}

//IsPlayerInRoom 检查Player是不是合法的本房间的Player,也就是排队服务器里发过来的要参加本局游戏的Player
func (pm *PlayersManager) IsPlayerInRoom(pId int32) bool {
	_, ok := pm.players[pId]
	return ok
}

func NewPlayersManager() *PlayersManager {
	pm := PlayersManager{
		players:      make(map[int32]*gt.Player),
		lock:         sync.RWMutex{},
		RegisterLock: sync.RWMutex{},
		hToP:         make(map[int32]int32),
		pToH:         make(map[int32]int32),
	}
	return &pm
}

//GetPlayer 获得一个玩家
func (pm *PlayersManager) GetPlayer(id int32) *gt.Player {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	return pm.players[id]
}

//AddPlayer 增加一个玩家
func (pm *PlayersManager) AddPlayer(p *gt.Player) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.players[p.Id] = p
}

//DeletePlayer 删除一个玩家
func (pm *PlayersManager) DeletePlayer(p *gt.Player) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	delete(pm.players, p.Id)
}

//UpdateHeartBeatTime 更新id为pid的玩家的心跳时间,t为unixNano的int64
func (pm *PlayersManager) UpdateHeartBeatTime(pid int32, t int64) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	//ti := time.Unix()
}
