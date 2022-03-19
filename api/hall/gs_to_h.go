package hall

import framework "melee_game_server/framework/game_room"

/**
*@Author Sly
*@Date 2022/2/16
*@Version 1.0
*@Description:
 */

type CreateNormalGameResponse struct {
	Ok             bool
	ConnectionInfo *framework.RoomConnectionInfo
}

type StartNormalGameResponse struct {
	Ok bool
}

// RoomStatus 结束游戏/返回游戏当前状态有关
type RoomStatus int

const (
	RoomInitStatus       RoomStatus = 1
	RoomStartStatus      RoomStatus = 2
	RoomDestroyingStatus RoomStatus = 3
)

type DestroyGameRoomResponse struct {
	Status RoomStatus
	Ok     bool
}

//GameAccountInfo 对局结算信息
type GameAccountInfo struct {
	StartTime        int64
	EndTime          int64
	PlayerAccountMap map[int32]*PlayerAccountInfo
}

//PlayerAccountInfo 玩家结算信息
type PlayerAccountInfo struct {
	Id        int32 //玩家id
	KillNum   int32 //击杀数
	AliveTime int64 //生存时间
}
