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
	/*	switch rm.Msg.TopMessageType {
		case proto.TopMessageType_BroadcastType:
			broad := rm.Msg.Broadcast
			switch  broad.BroadcastCode{
			case proto.BroadcastCode_HeroBulletLaunchBroadcastCode:
				logger.Testf("[Send]给%d名玩家发送广播%v", len(rm.ConnSlice),broad)
			case proto.BroadcastCode_HeroPositionReportBroadcastCode:
				logger.Testf("[Send]给%d名玩家发送广播%v", len(rm.ConnSlice),broad)
			}
		}*/
}

func (t *TestKcpNet) Shutdown() {
	logger.Testf("this should not be called")
}
