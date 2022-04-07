package server

import (
	"fmt"
	"melee_game_server/configs"
	"melee_game_server/plugins/logger"
)

/**
*@Author Sly
*@Date 2022/2/17
*@Version 1.0
*@Description:
 */

//DispatchMail 分发信息给各个房间的网络模块
func (gs *GameServer) DispatchMail() {
	for {
		mail := gs.Net.Receive()
		if mail.Msg == nil || mail.Msg.Request == nil {
			logger.Errorf("Dispatcher receive empty msg from:%s", mail.Conn.RemoteAddr())
			continue
		}
		if configs.ShowTcpMsg {
			fmt.Printf("receive:%v\n", mail.Msg.Request)
		}

		//修改在putMsg的时候对整体的room manager加读锁,这样不会出现没有删除房间,但是管道关闭,往空管道写数据引发panic的情况
		gs.grm.PutMsg(mail.Msg.Request.RoomId, mail)
	}
}
