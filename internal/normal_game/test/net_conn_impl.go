package test

import (
	"melee_game_server/plugins/logger"
	"net"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/25
*@Version 1.0
*@Description:用于测试用的实现net.Conn接口的结构
 */

type MyNetConn struct {
	hid int32
	pid int32
	msg string
}

func (m MyNetConn) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (m MyNetConn) Write(b []byte) (n int, err error) {
	logger.Testf("hid:%d,pid:%d,msg:%s", m.hid, m.pid, m.msg)
	return 0, nil
}

func (m MyNetConn) Close() error {
	panic("implement me")
}

func (m MyNetConn) LocalAddr() net.Addr {
	panic("implement me")
}

func (m MyNetConn) RemoteAddr() net.Addr {
	panic("implement me")
}

func (m MyNetConn) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (m MyNetConn) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (m MyNetConn) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}
