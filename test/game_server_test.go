package test

import (
	"fmt"
	"log"
	"melee_game_server/api/hall"
	framework "melee_game_server/framework/game_room"
	"net/rpc"
	"testing"
)

/**
*@Author Sly
*@Date 2022/2/16
*@Version 1.0
*@Description:
 */

func TestGameServerArrangeRoom(t *testing.T) {
	for i := 0; i < 10; i++ {
		createReq := new(hall.CreateNormalGameRequest)
		createReq.PlayerInfo = []*framework.PlayerInfo{
			{PlayerId: 1},
			{PlayerId: 2},
			{PlayerId: 3},
			{PlayerId: 4},
			{PlayerId: 5},
			{PlayerId: 6},
			{PlayerId: 7},
			{PlayerId: 8},
			{PlayerId: 9},
			{PlayerId: 10},
		}
		createResp := new(hall.CreateNormalGameResponse)
		callRpc("GameServer.CreateNormalGameRoom", createReq, createResp)

		startReq := new(hall.StartNormalGameRequest)
		startReq.RoomId = createResp.ConnectionInfo.Id
		startResp := new(hall.StartNormalGameResponse)
		callRpc("GameServer.StartNormalGame", startReq, startResp)
	}
}

func TestSimuHall(t *testing.T) {
	createReq := new(hall.CreateNormalGameRequest)
	createReq.PlayerInfo = []*framework.PlayerInfo{
		{PlayerId: 1},
		{PlayerId: 2},
		{PlayerId: 3},
		{PlayerId: 4},
		{PlayerId: 5},
		{PlayerId: 6},
		{PlayerId: 7},
		{PlayerId: 8},
		{PlayerId: 9},
		{PlayerId: 10},
	}
	createResp := new(hall.CreateNormalGameResponse)
	callRpc("GameServer.CreateNormalGameRoom", createReq, createResp)

	startReq := new(hall.StartNormalGameRequest)
	startReq.RoomId = createResp.ConnectionInfo.Id
	startResp := new(hall.StartNormalGameResponse)
	callRpc("GameServer.StartNormalGame", startReq, startResp)
}

func callRpc(rpcName string, args interface{}, reply interface{}) {
	c, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("Register尚未启动...")
	}

	err = c.Call(rpcName, args, reply)
	if err != nil {
		log.Fatal("rpc调用时出现:", err)
	}
	fmt.Printf("reply:%v\n", reply)
}
