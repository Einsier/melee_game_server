package game_type

import "sync"

/**
*@Author Sly
*@Date 2022/3/19
*@Version 1.0
*@Description:存放玩家荣誉,例如击杀数,存活时间等
 */

type PlayerHonor struct {
	mu sync.Mutex

	PlayerId  int32 //玩家id
	KillNum   int32 //击杀数
	AliveTime int64 //生存时间
}

func NewPlayerHonor(pid int32) *PlayerHonor {
	return &PlayerHonor{
		PlayerId:  pid,
		KillNum:   0,
		AliveTime: 0,
	}
}

func (ph *PlayerHonor) SetAliveTime(at int64) {
	ph.mu.Lock()
	defer ph.mu.Unlock()

	ph.AliveTime = at
}
