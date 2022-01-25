package test_kcp_net

import (
	"melee_game_server/framework/game_net/api"
	"melee_game_server/plugins/logger"
)

/**
*@Author Sly
*@Date 2022/1/24
*@Version 1.0
*@Description:
 */

type TestKcpNet struct {
}

func NewTestKcpNet() TestKcpNet {
	return TestKcpNet{}
}

func (t *TestKcpNet) Init(port string, recvSize, sendSize uint32) {
	logger.Test("test kcp net init success.")
}

func (t *TestKcpNet) Start() error {
	logger.Test("test kcp net start success.")
	return nil
}

func (t *TestKcpNet) Receive() *api.Mail {
	logger.Test("this should not be called")
	return nil
}

func (t *TestKcpNet) Send(rm *api.ReplyMail) {
	c := 0
	for _, conn := range rm.ConnSlice {
		if conn != nil {
			c++
		}
	}
	if c != 0 {
		logger.TestErrf("[TestKcpNet]给%d个nil Conn发送消息:%v\n", c, rm.Msg)
	}
	logger.Testf("[TestKcpNet]给%d个玩家发送如下消息:%v\n", len(rm.ConnSlice), rm.Msg)
}

func (t *TestKcpNet) Shutdown() {
	logger.Testf("this should not be called")
}
