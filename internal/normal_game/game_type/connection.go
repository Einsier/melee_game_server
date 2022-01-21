package game_type

import "net"

/**
*@Author Sly
*@Date 2022/1/17
*@Version 1.0
*@Description:
 */

type Connection struct {
	PlayerId int32
	HeroId   int32
	Conn     *net.Conn
}
