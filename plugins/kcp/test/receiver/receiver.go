package main

import (
	"fmt"
	"log"
	pb "melee_game_server/api/client/proto"
	adapter "melee_game_server/plugins/kcp/adapter"
	"net"

	kcp "github.com/xtaci/kcp-go"
)

func main() {
	conn, err := kcp.DialWithOptions("localhost:7777", nil, 0, 0)
	conn.SetNoDelay(1, 10, 2, 1)
	if err != nil {
		log.Fatal("Usage:client连接服务器失败，错误信息:", err)
		return
	}
	var msg pb.TopMessage
	msg.SeqId = 1
	send(conn, &msg)
	receive(conn)
}

func receive(conn net.Conn) {
	for {
		msg := adapter.Receive(conn)
		fmt.Print("接收到了:")
		fmt.Println(msg)
	}
}

func send(conn net.Conn, msg *pb.TopMessage) {
	adapter.Send(conn, msg)
	fmt.Print("发送了:")
	fmt.Println(msg)
}
