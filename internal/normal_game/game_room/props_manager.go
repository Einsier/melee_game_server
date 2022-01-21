package game_room

import (
	pb "melee_game_server/api/proto"
	gt "melee_game_server/internal/normal_game/game_type"
)

/**
*@Author Sly
*@Date 2022/1/17
*@Version 1.0
*@Description:
 */

type PropsManager struct {
	position gt.Vector2
	propType pb.PropType
}
