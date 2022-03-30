package game_room

import (
	"melee_game_server/framework/game_net/api"
)

type RoomInfo struct {
}
type GameInfo struct {
}

//PlayerInfo 用于normal_game一开始玩家进入游戏时的身份校验和玩家名称的显示
type PlayerInfo struct {
	PlayerId int32
	NickName string
}

type RoomInitInfo struct {
	Id          int32         //room分配的id
	Over        chan struct{} //是否结束的标志,如果结束则由game_room关闭管道通知game_server,后期可以改成对局结算信息
	JoinPlayers []*PlayerInfo //待加入的玩家的信息
}

type RoomConnectionInfo struct {
	Id         int32
	ClientAddr string
}

type GameRoom interface {
	Init(info *RoomInitInfo) //初始化一个GameRoom,over传入一个通道,游戏结束由GameRoom方关闭
	Start()                  //开始游戏
	GetGameInfo() *GameInfo  //获取游戏快照信息
	GetRoomInfo() *RoomInfo  //获取游戏房间信息
	PutMsg(mail *api.Mail)
	CloseMsgChan()
	//ForceStopGame 强制停止游戏,包括给所有正在游戏中的玩家发送游戏结束广播,断开所有的kcp/tcp连接,向外围的game_server发送
	//游戏结束信号(通过关闭信道等)
	ForceStopGame() (ok bool)
	//GetGameAccount()*hall.GameAccountInfo
	GetGameAccount() interface{}
}
