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

func main() {
	flag.Parse()
	server.GS.HallRpcPort = *hallRpcPortFlag
	server.GS.ClientPort = *clientPortFlag
	server.GS.Run()
	logger.Infof("开启game server:大厅服务器rpc端口[%s],客户端kcp端口[%s]\n", *hallRpcPortFlag, *clientPortFlag)

	time.Sleep(100 * time.Minute)
}
