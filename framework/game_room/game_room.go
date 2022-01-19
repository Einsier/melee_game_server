package game_room

/**
*@Author Sly
*@Date 2022/1/17
*@Version 1.0
*@Description:game_room的接口
 */

type RoomInfo struct {
}
type GameInfo struct {
}

type GameRoom interface {
	Init(id int, port string, over chan interface{}) //初始化一个GameRoom,over传入一个通道,游戏结束由GameRoom方关闭
	Start()                                          //开始游戏
	GetGameInfo() *GameInfo                          //获取游戏快照信息
	GetRoomInfo() *RoomInfo                          //获取游戏房间信息
}
