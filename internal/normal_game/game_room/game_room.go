package game_room

import (
	"melee_game_server/framework/game_room"
	gc "melee_game_server/internal/normal_game/game_controller"
	gn "melee_game_server/internal/normal_game/game_net"
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

const (
	NormalGameIdleStatus       = iota //被创建,但没有被初始化
	NormalGameInitStatus              //已经被初始化
	NormalGameWaitPlayerStatus        //等待Player的到来
	NormalGameStartStatus             //全部Player已经到来,游戏开始
	NormalGameGameOverStatus          //全部Player狗带,游戏结束
)

type NormalGameRoom struct {
	Id                  int32
	port                string
	Prepare             chan interface{} //所有玩家都已连入游戏
	over                chan interface{} //用于向game_server汇报的
	PlayerNum           int32            //当前的存活的玩家数目
	Status              int32            //当前的状态
	heroManager         *HeroesManager
	propsManager        *PropsManager
	bulletsManager      *BulletsManager
	netServer           *gn.NormalGameNetServer
	playerManager       *PlayerManager
	requestController   *gc.RequestController
	timeEventController *gc.TimeEventController
}

func (room *NormalGameRoom) GetHeroesManager() *HeroesManager {
	return room.heroManager
}
func (room *NormalGameRoom) GetBulletsManager() *BulletsManager {
	return room.bulletsManager
}
func (room *NormalGameRoom) GetNetServer() *gn.NormalGameNetServer {
	return room.netServer
}
func (room *NormalGameRoom) GetPlayerManager() *PlayerManager {
	return room.playerManager
}
func (room *NormalGameRoom) GetPropsManager() *PropsManager {
	return room.propsManager
}
func (room *NormalGameRoom) GetTimeEventController() *gc.TimeEventController {
	return room.timeEventController
}

//Init 初始化GameRoom的参数
func (room *NormalGameRoom) Init(info *game_room.RoomInitInfo) {

}

//Start 开始游戏,先开启kcp服务器,接收玩家的PlayerEnterGameRequest请求,根据init传过来的JoinPlayers进行身份校验,
//把PlayerManager和NormalGameNetServer的Player部分进行初始化,初始化playerId,heroId(可以参考学长的uuid进行分配),
//和网络连接的映射关系
func (room *NormalGameRoom) Start() {
	//todo 开启kcp....
	<-room.Prepare
	//PlayerManager
}

//GetGameInfo 获取房间的游戏快照信息
func (room *NormalGameRoom) GetGameInfo() *game_room.GameInfo {
	return nil
}

//GetRoomInfo 获取游戏房间信息
func (room *NormalGameRoom) GetRoomInfo() *game_room.RoomInfo {
	return nil
}
