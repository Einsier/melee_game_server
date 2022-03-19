package hall

import framework "melee_game_server/framework/game_room"

/**
*@Author Sly
*@Date 2022/2/16
*@Version 1.0
*@Description: 大厅服务器向game_server发送的内容
 */

type CreateNormalGameRequest struct {
	PlayerInfo []*framework.PlayerInfo
	GameId     string //作为etcd的路径,游戏结束之后由server通知hall,方便hall落库
}

type StartNormalGameRequest struct {
	RoomId int32
}

type DestroyGameRoomRequest struct {
	RoomId int32
}
