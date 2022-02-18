package adapter

import (
	"fmt"
	"log"
	pb "melee_game_server/api/client/proto"
	codec "melee_game_server/plugins/kcp/codec"
	"net"
)

const (
	bufferSize int32 = 1024
)

func Send(conn net.Conn, msg *pb.TopMessage) {
	cMsg := codec.Code(msg)
	count, err := conn.Write(cMsg)
	if err != nil {
		log.Println("Usage:Send()调用失败，错误信息:", err)
	} else {
		fmt.Println("Msg:Send()调用成功，共发送", count, "字节")
	}
}

func Receive(conn net.Conn) *pb.TopMessage {
	buffer := make([]byte, bufferSize)
	count, err := conn.Read(buffer)
	if err != nil {
		log.Println(err)
		return nil
	}
	msg := codec.Decode(buffer[0:count])
	return msg
}
