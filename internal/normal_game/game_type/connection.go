package game_type

import "net"

/**
*@Author Sly
*@Date 2022/1/17
*@Version 1.0
*@Description:
 */

type Connection struct {
	playerId int32
	heroId   int32
	conn     *net.Conn
}
