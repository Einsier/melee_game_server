package game_net

import (
	"melee_game_server/api/proto"
	gn "melee_game_server/framework/game_net/api"
	"melee_game_server/plugins/kcp_net/mailbox"
	"net"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/18
*@Version 1.0
*@Description:
 */

//NormalGameNetServer 普通游戏(玩家全部进入才开始)的GameNetServer,全部玩家都完成与game_server的注册之后才开始游戏,
//所以初步设计在SendByHeroId和SendByPlayerId中访问两个map不加锁,因为中间不会产生对map的增删操作
type NormalGameNetServer struct {
	np         gn.NetPlugin
	heroConn   map[int32]*net.Conn
	playerConn map[int32]*net.Conn
	port       string
	lock       sync.Mutex
}

func NewNormalGameNetServer() *NormalGameNetServer {
	kcpNet := mailbox.Mailbox{}
	return &NormalGameNetServer{
		np:         &kcpNet,
		heroConn:   make(map[int32]*net.Conn),
		playerConn: make(map[int32]*net.Conn),
		lock:       sync.Mutex{},
	}
}

func (ngs *NormalGameNetServer) Init(port string, recvSize, sendSize uint32) {
	ngs.np.Init(port, recvSize, sendSize)
}

func (ngs *NormalGameNetServer) Start() {
	if err := ngs.np.Start(); err != nil {
		panic(err)
	}
}

func (ngs *NormalGameNetServer) Register(playerId, heroId int32, conn *net.Conn) {
	ngs.lock.Lock()
	defer ngs.lock.Unlock()
	ngs.playerConn[playerId] = conn
	ngs.heroConn[heroId] = conn
}

func (ngs *NormalGameNetServer) SendBySingleConn(conn *net.Conn, msg *proto.TopMessage) {
	s := []*net.Conn{conn}
	ngs.np.Send(gn.NewReplyMail(s, msg))
}

func (ngs *NormalGameNetServer) SendByConn(conn []*net.Conn, msg *proto.TopMessage) {
	ngs.np.Send(gn.NewReplyMail(conn, msg))
}

func (ngs *NormalGameNetServer) SendByHeroId(hIdSlice []int32, msg *proto.TopMessage) {
	sendConn := make([]*net.Conn, 0)
	for _, hId := range hIdSlice {
		conn := ngs.heroConn[hId]
		if conn != nil {
			sendConn = append(sendConn, conn)
		}
	}
	ngs.np.Send(gn.NewReplyMail(sendConn, msg))
}

func (ngs *NormalGameNetServer) SendByPlayerId(pIdSlice []int32, msg *proto.TopMessage) {
	sendConn := make([]*net.Conn, 0)
	for _, hId := range pIdSlice {
		conn := ngs.heroConn[hId]
		if conn != nil {
			sendConn = append(sendConn, conn)
		}
	}
	ngs.np.Send(gn.NewReplyMail(sendConn, msg))
}

//Receive 如果没有消息则阻塞
func (ngs *NormalGameNetServer) Receive() *proto.TopMessage {
	return ngs.np.Receive().Msg
}
