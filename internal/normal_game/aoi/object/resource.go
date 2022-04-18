package object

import "melee_game_server/framework/entity"

/**
*@Author Sly
*@Date 2022/4/12
*@Version 1.0
*@Description:
 */

type Resource struct {
	Position		entity.Vector2		//左下角的位置
	Height			float64
	Weight			float64
}