package main

import (
	"flag"
	"melee_game_server/configs"
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

//腾讯服务器地址:1.116.109.113
var hallRpcPortFlag = flag.String("hallRpcPort", ":8000", "set the port of rpc in order to communicate with hall")
var clientTcpAddrFlag = flag.String("clientAddr", "49.234.245.172:32004", "set the port of tcp in order to communicate with client")
var etcdAddrFlag = flag.String("etcdAddr", "42.192.200.194:2379", "set the address of etcd")
var testFlag = flag.Bool("t", false, "if this is a local test")

func ParseFlags() {
	flag.Parse()
	if *testFlag {
		//如果当前是本机测试
		*hallRpcPortFlag = ":8000"
		*clientTcpAddrFlag = "localhost:8001"
		*etcdAddrFlag = "42.192.200.194:2379"
	}
}

func main() {
	ParseFlags()
	server.GS.HallRpcPort = *hallRpcPortFlag
	server.GS.ClientAddr = *clientTcpAddrFlag
	configs.EtcdAddr = *etcdAddrFlag
	server.EtcdCli = server.NewEtcdCli()
	server.GS.Run()
	logger.Infof("开启game server:大厅服务器rpc端口[%s],客户端kcp地址[%s]\n", *hallRpcPortFlag, *clientTcpAddrFlag)
	time.Sleep(100 * time.Minute)
}
