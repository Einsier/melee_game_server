package game_room

import (
	"melee_game_server/api/hall"
	gt "melee_game_server/internal/normal_game/game_type"
	"sync"
)

/**
*@Author Sly
*@Date 2022/3/19
*@Version 1.0
*@Description:用于管理玩家本次对局中的荣誉,例如击杀数,生存时间等
 */

type HonorManager struct {
	mu     sync.RWMutex
	honors map[int32]*gt.PlayerHonor
}

func NewHonorManager() *HonorManager {
	return &HonorManager{
		honors: make(map[int32]*gt.PlayerHonor),
		mu:     sync.RWMutex{},
	}
}

//AddPlayerHonor 添加某个玩家的荣誉
func (hm *HonorManager) AddPlayerHonor(pid int32) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hm.honors[pid] = gt.NewPlayerHonor(pid)
}

//GetPlayerHonor 获取某个玩家的荣誉
func (hm *HonorManager) GetPlayerHonor(pid int32) *gt.PlayerHonor {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	return hm.honors[pid]
}

func (hm *HonorManager) GetAllPlayerHonor() map[int32]*hall.PlayerAccountInfo {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	m := make(map[int32]*hall.PlayerAccountInfo)
	for pid, honor := range hm.honors {
		pa := new(hall.PlayerAccountInfo)
		pa.Id = pid
		pa.KillNum = honor.KillNum
		pa.AliveTime = honor.AliveTime

		m[pid] = pa
	}
	return m
}
