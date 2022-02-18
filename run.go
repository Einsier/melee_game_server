package main

import (
	"flag"
	"fmt"
	"melee_game_server/server"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/24
*@Version 1.0
*@Description:mvp1版本的game_room的暂时的启动方式
 */

var portFlag = flag.String("port", "1234", "set the port of this game server")

func main() {
	flag.Parse()
	server.GS.Port = *portFlag
	server.GS.Run()
	fmt.Printf("start new game server at:0.0.0.0:%s\n", *portFlag)
	time.Sleep(100 * time.Minute)
}
