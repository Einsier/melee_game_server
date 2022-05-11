package adapter

import (
	"fmt"
	pb "melee_game_server/api/client/proto"
	"melee_game_server/configs"
	codec "melee_game_server/plugins/kcp/codec"
	"net"
)

const (
	bufferSize int32 = 1024
)

func Send(conn net.Conn, msg *pb.TopMessage) {
	cMsg := codec.Code(msg)
	//count, err := conn.Write(cMsg)
	_, err := conn.Write(cMsg)
	if err != nil {
		//log.Println("Usage:Send()调用失败，错误信息:", err)
	} else {
		if configs.ShowTcpMsg {
			fmt.Printf("%v\n", msg)
		}
		//fmt.Println("Msg:Send()调用成功，共发送", count, "字节")
	}
}

func Receive(conn net.Conn) *pb.TopMessage {
	buffer := make([]byte, bufferSize)
	count, err := conn.Read(buffer)
	if err != nil {
		//logger.Infof("tcp:检测到%s断开与gs的连接", conn.RemoteAddr())
		return nil
	}
	msg := codec.Decode(buffer[0:count])
	return msg
}
