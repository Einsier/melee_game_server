package main

import (
	"flag"
	"melee_game_server/plugins/logger"
	"melee_game_server/server"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/24
*@Version 1.0
*@Description:mvp1版本的game_room的暂时的启动方式
 */

var hallRpcPortFlag = flag.String("hallRpcPort", ":8000", "set the port of rpc in order to communicate with hall")
var clientPortFlag = flag.String("clientPort", ":8001", "set the port of kcp in order to communicate with client")

/*func SimuHall() {
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
	c, err := rpc.DialHTTP("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal("Register尚未启动...")
	}

	err = c.Call(rpcName, args, reply)
	if err != nil {
		log.Fatal("rpc调用时出现:", err)
	}
}*/

func main() {
	flag.Parse()
	server.GS.HallRpcPort = *hallRpcPortFlag
	server.GS.ClientPort = *clientPortFlag
	server.GS.Run()

	logger.Infof("开启game server:大厅服务器rpc端口[%s],客户端kcp端口[%s]\n", *hallRpcPortFlag, *clientPortFlag)
	time.Sleep(10 * time.Millisecond)

	//SimuHall()
	time.Sleep(100 * time.Minute)
}
