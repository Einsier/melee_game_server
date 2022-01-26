package api

import (
	"melee_game_server/api/proto"
	"net"
)

/**
*@Author Sly
*@Date 2022/1/18
*@Version 1.0
*@Description:直接作为game_room的一部分使用
 */

type GameNetServer interface {
	Init(port string, recvSize, sendSize uint32)    //初始化
	Start()                                         //开始工作
	Register(playerId, heroId int32, conn net.Conn) //注册playerId,heroId,对应连接的对应关系
	SendByHeroId(hId []int32, msg *proto.TopMessage)
	SendByPlayerId(pId []int32, msg *proto.TopMessage)
	Receive() *Mail
	SendByConn(conn []net.Conn, msg *proto.TopMessage)
	SendBySingleConn(conn net.Conn, msg *proto.TopMessage)
}
