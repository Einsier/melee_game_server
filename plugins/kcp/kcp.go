package kcp

import "melee_game_server/plugins/kcp/mailbox"

/**
*@Author Sly
*@Date 2022/2/17
*@Version 1.0
*@Description:
 */

var KCP *mailbox.Mailbox = initKCP()

func StartKCP(addr string, recvSize, sendSize uint32) {
	KCP.Init(addr, recvSize, sendSize)
	if err := KCP.Start(); err != nil {
		panic(err)
	}
}

//initKCP 单例模式,初始化kcp
func initKCP() *mailbox.Mailbox {
	return new(mailbox.Mailbox)
}
