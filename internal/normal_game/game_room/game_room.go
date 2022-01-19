package game_room

import (
	"melee_game_server/framework/game_room"
	"melee_game_server/internal/normal_game/game_net"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:用于管理一次游戏的对局,包括玩家管理,对局中的prop(道具)管理以及子弹管理
 */

type PlayerManager struct {
}

type Player struct {
}

type NormalGameRoom struct {
	Id             int
	port           string
	heroManager    HeroesManager
	propsManager   PropsManager
	bulletsManager BulletsManager
	netServer      game_net.NormalGameNetServer
	playerManager  PlayerManager
	over           chan interface{}
	PlayerNum      int32
}

//Init 初始化GameRoom的参数
func (room *NormalGameRoom) Init(info *game_room.RoomInitInfo) {

}

//Start 开始游戏,先开启kcp服务器,接收玩家的PlayerEnterGameRequest请求,根据init传过来的JoinPlayers进行身份校验,
//把PlayerManager和NormalGameNetServer的Player部分进行初始化,初始化playerId,heroId(可以参考学长的uuid进行分配),
//和网络连接的映射关系
func (room *NormalGameRoom) Start() {

}

//GetGameInfo 获取房间的游戏快照信息
func (room *NormalGameRoom) GetGameInfo() *game_room.GameInfo {
	return nil
}

//GetRoomInfo 获取游戏房间信息
func (room *NormalGameRoom) GetRoomInfo() *game_room.RoomInfo {
	return nil
}
