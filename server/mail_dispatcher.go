package server

import (
	"fmt"
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
			logger.Errorf("receive error msg:%v", mail.Msg)
			continue
		}
		fmt.Printf("receive:%v\n", mail.Msg.Request)
		room, ok := gs.grm.GetRoom(mail.Msg.Request.RoomId)
		if !ok {
			logger.Errorf("receive room not exist msg:%v", mail.Msg)
			continue
		}
		room.PutMsg(mail)
	}
}
