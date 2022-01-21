package game_room

/**
*@Author chenjiajia
*@Date 2022/1/17
*@Version 1.0
*@Description:game_room的接口
 */

type RoomInfo struct {
}
type GameInfo struct {
}

//PlayerInfo 暂时只有playerId,用于normal_game一开始玩家进入游戏时的身份校验
type PlayerInfo struct {
	PlayerId int32
}

type RoomInitInfo struct {
	Id          int32            //room分配的id
	Port        string           //room分配的端口
	Over        chan interface{} //是否结束的标志,如果结束则由game_room关闭管道通知game_server
	JoinPlayers []*PlayerInfo    //待加入的玩家的信息
}

type GameRoom interface {
	Init(info *RoomInitInfo) //初始化一个GameRoom,over传入一个通道,游戏结束由GameRoom方关闭
	Start()                  //开始游戏
	GetGameInfo() *GameInfo  //获取游戏快照信息
	GetRoomInfo() *RoomInfo  //获取游戏房间信息
}
