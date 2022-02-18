package api

import (
	pb "melee_game_server/api/client/proto"
	"net"
)

/**
*@Author Sly
*@Date 2022/1/17
*@Version 1.0
*@Description:网络插件需要实现以下接口,才可以被注册到game_net_server中使用
 */

type Mail struct {
	Conn net.Conn
	Msg  *pb.TopMessage
}

type ReplyMail struct {
	ConnSlice []net.Conn
	Msg       *pb.TopMessage
}

func NewReplyMail(connSlice []net.Conn, msg *pb.TopMessage) *ReplyMail {
	return &ReplyMail{
		ConnSlice: connSlice,
		Msg:       msg,
	}
}

type NetPlugin interface {
	Init(port string, recvSize, sendSize uint32)
	Start() error   //开始工作
	Receive() *Mail //接收信息
	//Listen(conn net.Conn)*Mail
	Send(replyPtr *ReplyMail) //发送
	Shutdown()                //停止工作
}
