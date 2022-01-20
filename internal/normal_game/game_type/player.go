package game_type

/**
*@Author chenjiajia
*@Date 2022/1/19
*@Version 1.0
*@Description: 存放对战中玩家有关内容
 */

type Player struct {
	Id       int32  //本局游戏中的玩家id
	Nickname string //玩家昵称
	Level    int    //玩家等级
	Score    int    //玩家得分
}

func NewPlayer(id int32) *Player {
	return &Player{
		Id: id,
	}
}
