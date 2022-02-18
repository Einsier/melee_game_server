package server

import (
	"melee_game_server/api/hall"
	"melee_game_server/plugins/logger"
)

/**
*@Author Sly
*@Date 2022/2/15
*@Version 1.0
*@Description:大厅服务器排完队,凑齐游戏人数之后调用的api
 */

//CreateNormalGameRoom hall排完队之后通知game_server,开启一个game_room,用于玩家加入对局
func (gs *GameServer) CreateNormalGameRoom(req *hall.CreateNormalGameRequest, resp *hall.CreateNormalGameResponse) error {
	connectInfo, err := gs.grm.AddNormalGameRoom(req.PlayerInfo)
	if err != nil {
		resp.Ok = false
		return err
	}
	resp.ConnectionInfo = connectInfo
	resp.Ok = true
	logger.Testf("分配%d房间号给hall", resp.ConnectionInfo.Id)
	return nil
}

//StartNormalGame hall通知所有的玩家就绪之后,通知NormalGameServer中的room开始游戏
func (gs *GameServer) StartNormalGame(req *hall.StartNormalGameRequest, resp *hall.StartNormalGameResponse) error {
	err := gs.grm.StartNormalGame(req.RoomId)
	if err != nil {
		resp.Ok = false
		return err
	}
	resp.Ok = true
	return nil
}
