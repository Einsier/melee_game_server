package game_type

import "melee_game_server/api/client/proto"

/**
*@Author Sly
*@Date 2022/1/17
*@Version 1.0
*@Description:游戏内的道具,具体可以见api/proto/consts下的PropType枚举类型
 */

type Prop struct {
	Id       int32
	PropType proto.PropType
	Position Vector2
}

func NewProp(id int32, propType proto.PropType, position Vector2) *Prop {
	return &Prop{
		Id:       id,
		PropType: propType,
		Position: position,
	}
}
