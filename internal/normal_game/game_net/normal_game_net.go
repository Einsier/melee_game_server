package game_net

import (
	"melee_game_server/api/client/proto"
	gn "melee_game_server/framework/game_net/api"
	"melee_game_server/plugins/kcp"
	"net"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/18
*@Version 2.0
*@Description:每个game_room拥有一个的网络接口.用于从全局消息队列中接收消息,或者发给玩家
2.0版本:整个的game_server使用一个kcp网络
*/

//NormalGameNetServer 普通游戏(玩家全部进入才开始)的GameNetServer,全部玩家都完成与game_server的注册之后才开始游戏,
//所以初步设计在SendByHeroId和SendByPlayerId中访问两个map不加锁,因为中间不会产生对map的增删操作
type NormalGameNetServer struct {
	np         gn.NetPlugin
	RoomId     int32
	ReqChan    chan *gn.Mail      //2.0版本,有一个接收消息的信道
	heroConn   map[int32]net.Conn //key为heroId,value为联系方式
	playerConn map[int32]net.Conn //key为playerId,value为联系方式
	lock       sync.Mutex
}

//PutReqMsg 2.0版本,向网络的NormalGameNetServer添加
func (ngs *NormalGameNetServer) PutReqMsg(msg *gn.Mail) {
	ngs.ReqChan <- msg
}

func NewNormalGameNetServer(roomId int32) *NormalGameNetServer {
	return &NormalGameNetServer{
		RoomId:     roomId,
		ReqChan:    make(chan *gn.Mail, 8192),
		np:         kcp.KCP,
		heroConn:   make(map[int32]net.Conn),
		playerConn: make(map[int32]net.Conn),
		lock:       sync.Mutex{},
	}
}

func (ngs *NormalGameNetServer) Register(playerId, heroId int32, conn net.Conn) {
	ngs.lock.Lock()
	defer ngs.lock.Unlock()
	ngs.playerConn[playerId] = conn
	ngs.heroConn[heroId] = conn
}

func (ngs *NormalGameNetServer) SendBySingleConn(conn net.Conn, msg *proto.TopMessage) {
	s := []net.Conn{conn}
	ngs.np.Send(gn.NewReplyMail(s, msg))
}

func (ngs *NormalGameNetServer) SendByConn(conn []net.Conn, msg *proto.TopMessage) {
	ngs.np.Send(gn.NewReplyMail(conn, msg))
}

func (ngs *NormalGameNetServer) SendByHeroId(hIdSlice []int32, msg *proto.TopMessage) {
	sendConn := make([]net.Conn, 0)
	for _, hId := range hIdSlice {
		conn := ngs.heroConn[hId]
		if conn != nil {
			sendConn = append(sendConn, conn)
		}
	}
	ngs.np.Send(gn.NewReplyMail(sendConn, msg))
}

//因为normal_game_net无法感知当前在线的玩家数目(因为为了效率没有加锁,所以不能中途增删改玩家 map,所以只能上层来检测当前的在线玩家.)
////SendToAllPlayerConn 向所有的Player的Conn发送消息
//func (ngs *NormalGameNetServer) SendToAllPlayerConn(msg *proto.TopMessage) {
//	sendConn := make([]net.Conn, 0)
//	for _, conn := range ngs.playerConn {
//		if conn != nil {
//			sendConn = append(sendConn, conn)
//		}
//	}
//	ngs.np.Send(gn.NewReplyMail(sendConn, msg))
//}

func (ngs *NormalGameNetServer) SendByPlayerId(pIdSlice []int32, msg *proto.TopMessage) {
	sendConn := make([]net.Conn, 0)
	for _, pId := range pIdSlice {
		conn := ngs.playerConn[pId]
		if conn != nil {
			sendConn = append(sendConn, conn)
		}
	}
	ngs.np.Send(gn.NewReplyMail(sendConn, msg))
}

//Receive 如果没有消息则阻塞
func (ngs *NormalGameNetServer) Receive() (*gn.Mail, bool) {
	mail, ok := <-ngs.ReqChan
	return mail, ok
}

//DeleteConn player退出游戏的时候删除player和hero的联系方式，不可以在游戏中间调用
func (ngs *NormalGameNetServer) DeleteConn(hid, pid int32) {
	ngs.lock.Lock()
	defer ngs.lock.Unlock()
	delete(ngs.playerConn, pid)
	delete(ngs.heroConn, hid)
}
