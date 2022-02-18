package server

import (
	"log"
	"melee_game_server/configs"
	"melee_game_server/framework/game_net/api"
	framework "melee_game_server/framework/game_room"
	ngr "melee_game_server/internal/normal_game/game_room"
	"melee_game_server/plugins/kcp"
	"net"
	"net/http"
	"net/rpc"
)

/**
*@Author Sly
*@Date 2022/2/15
*@Version 1.0
*@Description:
 */

var GS *GameServer = newGameServer()

var GameTypeIdMap = map[int32]framework.GameRoom{
	configs.NormalGameRoomId: &ngr.NormalGameRoom{},
}

type GameServer struct {
	grm  *GameRoomManger
	Port string
	Net  api.NetPlugin
}

//newGameServer 创建一个GameServer,单例模式,所以不对外暴露
func newGameServer() *GameServer {
	gs := new(GameServer)
	gs.grm = new(GameRoomManger)
	gs.grm.gameRooms = make(map[int32]framework.GameRoom)
	gs.Net = kcp.KCP
	return gs
}

//Run 开启对内部集群的rpc服务,开启kcp网络等
func (gs *GameServer) Run() {
	//开启rpc服务
	//todo 根据部署后期调整
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("无法监听:", e)
	}
	gs.serveRpc(l)

	//开启监听客户端的kcp服务
	kcp.StartKCP("0.0.0.0:"+gs.Port, 1024, 1024)

	go gs.DispatchMail()
}

//doServer 将自己的方法注册成rpc服务
func (gs *GameServer) serveRpc(rpcListener net.Listener) {
	err := rpc.Register(gs)
	if err != nil {
		log.Fatal("无法注册rpc:", err)
	}
	rpc.HandleHTTP()
	go func() {
		err := http.Serve(rpcListener, nil)
		if err != nil {
			log.Fatal("rpc无法正确监听:", err)
		}
	}()
}
