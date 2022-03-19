package server

import (
	"errors"
	"fmt"
	"math/rand"
	"melee_game_server/api/hall"
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

	/*	//todo 测试用,只固定开放房间1的话把下面的注释去掉,把上面的for注释掉
		grm.mu.Lock()
		_, ok := grm.gameRooms[1]
		if ok {
			grm.mu.Unlock()
			panic("room1 已经被占用")
		}
		grm.gameRooms[1] = *new(framework.GameRoom)
		grm.mu.Unlock()
		ret = 1*/

	return ret
}

//AddNormalGameRoom 创建一个 NormalGameRoom 并且把它加入到gameRooms中
func (grm *GameRoomManger) AddNormalGameRoom(playerInfo []*framework.PlayerInfo, gameId string) (*framework.RoomConnectionInfo, error) {
	room := new(ngr.NormalGameRoom)

	info := new(framework.RoomInitInfo)
	info.Over = make(chan struct{})
	info.JoinPlayers = playerInfo
	info.Id = grm.GetFreeId()

	go func(roomId int32, room framework.GameRoom) {
		<-info.Over
		//todo 结束事件,数据库持久化对局信息等
		fmt.Printf("对局:%v已结束\n", roomId)

		accountInfo := room.GetGameAccount().(*hall.GameAccountInfo)
		SetAccountToEtcd(gameId, accountInfo)
		logger.Infof("gameId:%s 的对局结算信息已放入etcd", gameId)
		grm.DeleteGameRoom(roomId)
	}(info.Id, room)

	grm.mu.Lock()
	grm.gameRooms[info.Id] = room
	grm.mu.Unlock()

	room.Init(info)

	connInfo := new(framework.RoomConnectionInfo)
	connInfo.Id = info.Id
	return connInfo, nil
}

//DeleteGameRoom 只是从注册列表中删除GameRoom,不会对GameRoom内部逻辑产生影响
func (grm *GameRoomManger) DeleteGameRoom(roomId int32) {
	grm.mu.Lock()
	defer grm.mu.Unlock()

	delete(grm.gameRooms, roomId)
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
