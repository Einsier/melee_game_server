package server

import (
	"errors"
	"fmt"
	"math/rand"
	framework "melee_game_server/framework/game_room"
	ngr "melee_game_server/internal/normal_game/game_room"
	"melee_game_server/plugins/logger"
	"sync"
)

/**
*@Author Sly
*@Date 2022/2/17
*@Version 1.0
*@Description:
 */

type GameRoomManger struct {
	mu        sync.RWMutex
	gameRooms map[int32]framework.GameRoom
}

//GetFreeId 使用随机数的方式,返回并占用一个空闲id.
func (grm *GameRoomManger) GetFreeId() (ret int32) {
	for {
		ret = rand.Int31()
		grm.mu.Lock()
		_, ok := grm.gameRooms[ret]
		if ok {
			grm.mu.Unlock()
			continue
		}
		//填上一个指针,防止其他GetFreeId检测为空
		grm.gameRooms[ret] = *new(framework.GameRoom)
		grm.mu.Unlock()
		break
	}
	return ret
}

//AddNormalGameRoom 创建一个 NormalGameRoom 并且把它加入到gameRooms中
func (grm *GameRoomManger) AddNormalGameRoom(playerInfo []*framework.PlayerInfo) (*framework.RoomConnectionInfo, error) {
	room := new(ngr.NormalGameRoom)

	info := new(framework.RoomInitInfo)
	info.Over = make(chan struct{})
	info.JoinPlayers = playerInfo
	info.Id = grm.GetFreeId()

	go func(roomId int32) {
		<-info.Over
		//todo 结束事件,数据库持久化对局信息等
		fmt.Printf("对局:%v已结束\n", roomId)
	}(info.Id)

	grm.mu.Lock()
	grm.gameRooms[info.Id] = room
	grm.mu.Unlock()

	room.Init(info)

	connInfo := new(framework.RoomConnectionInfo)
	connInfo.Id = info.Id
	return connInfo, nil
}

//StartNormalGame 开始游戏
func (grm *GameRoomManger) StartNormalGame(id int32) error {
	room, ok := grm.gameRooms[id]
	if !ok {
		return errors.New("room do not exist")
	}
	go room.Start()
	logger.Testf("%d号房间开始游戏", id)
	return nil
}

func (grm *GameRoomManger) GetRoom(id int32) (framework.GameRoom, bool) {
	grm.mu.RLock()
	defer grm.mu.RUnlock()
	room, ok := grm.gameRooms[id]
	return room, ok
}
