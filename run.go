package main

import (
	"flag"
	"log"
	"melee_game_server/configs"
	"melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/internal/normal_game/metrics"
	"melee_game_server/plugins/logger"
	"melee_game_server/server"
	"os"
	"runtime"
	"strconv"
	"strings"
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
var clientTcpAddrFlag = flag.String("clientAddr", "1.116.109.113:8001", "set the port of tcp in order to communicate with client")
var etcdAddrFlag = flag.String("etcdAddr", "42.192.200.194:2379", "set the address of etcd")
var testFlag = flag.Bool("t", false, "if this is a none-deploy test")
var playerNumFlag = flag.Int("playerNum", 3, "configs the number of players in each game which must be same as the hall's config")
var showTcpMsg = flag.Bool("msg", false, "whether to display Tcp messages receive from clients/send to clients")

func ParseFlags() {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)
	flag.Parse()
	configs.ShowTcpMsg = *showTcpMsg
	if *testFlag {
		//如果当前是本机测试
		*hallRpcPortFlag = ":8000"
		*etcdAddrFlag = "42.192.200.194:2379"
		configs.IsTest = true
	} else {
		//如果当前是集群部署,查看自己的环境变量然后设置监听外网端口
		hostname := os.Getenv("HOSTNAME")
		if len(hostname) == 0 {
			log.Fatalln("请设置-t参数表示当前非集群部署,或者设置HOSTNAME参数")
		}
		//分割类似"game-server-0"这样的字符串
		sp := strings.Split(hostname, "-")
		if len(sp) == 0 {
			log.Fatalf("HOSTNAME不正确,hostname:%s", hostname)
		}
		if id, err := strconv.Atoi(sp[len(sp)-1]); err != nil {
			log.Fatalf("HOSTNAME不正确,hostname:%s", hostname)
		} else {
			*clientTcpAddrFlag = "49.234.245.172:" + strconv.Itoa(33000+id)
		}
	}
}

func main() {
	ParseFlags()
	server.GS.HallRpcPort = *hallRpcPortFlag
	server.GS.ClientAddr = *clientTcpAddrFlag
	configs.EtcdAddr = *etcdAddrFlag
	normal_game_type_configs.MaxNormalGamePlayerNum = int32(*playerNumFlag)
	server.EtcdCli = server.NewEtcdCli()
	go metrics.Start()
	server.GS.Run()
	logger.Infof("开启game server:大厅服务器rpc端口[%s],客户端kcp地址[%s]\n", *hallRpcPortFlag, *clientTcpAddrFlag)
	time.Sleep(1000000 * time.Minute)
}
