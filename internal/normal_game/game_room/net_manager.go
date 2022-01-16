package game_room

import "melee_game_server/internal/normal_game/game_type"

/**
*@Author Sly
*@Date 2022/1/17
*@Version 1.0
*@Description:用于管理playerId,heroId,以及对应的connection
 */

type NetManager struct {
	connections map[int32]game_type.Connection //key:heroId
}
