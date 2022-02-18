package main

import (
	"fmt"
	mailbox "melee_game_server/plugins/kcp/mailbox"
)

func main() {
	var box mailbox.Mailbox
	box.Init("0.0.0.0:7777", 10, 10)

	box.Start()

	for {
		msg := box.Receive()
		fmt.Println(*msg)
		msg.Msg.SeqId = msg.Msg.SeqId + 1000
		box.Broadcast(msg.Msg)
	}

}
