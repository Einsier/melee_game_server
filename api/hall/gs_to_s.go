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
