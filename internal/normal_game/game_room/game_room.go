package game_room

import (
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/game_room"
	gn "melee_game_server/internal/normal_game/game_net"
	gt "melee_game_server/internal/normal_game/game_type"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:用于管理一次游戏的对局,包括玩家管理,对局中的prop(道具)管理以及子弹管理
 */

type NormalGameRoom struct {
	Id                  int32
	port                string
	Prepare             chan interface{} //所有玩家都已连入游戏
	leave               chan interface{} //玩家都已经离开
	over                chan interface{} //用于向game_server汇报的
	PlayerNum           int32            //当前的存活的玩家数目
	Status              int32            //当前的状态
	heroManager         *HeroesManager
	propsManager        *PropsManager
	bulletsManager      *BulletsManager
	netServer           *gn.NormalGameNetServer
	playersManager      *PlayersManager
	requestController   *RequestController
	timeEventController *TimeEventController
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
func (room *NormalGameRoom) GetPlayerManager() *PlayersManager {
	return room.playersManager
}
func (room *NormalGameRoom) GetPropsManager() *PropsManager {
	return room.propsManager
}
func (room *NormalGameRoom) GetTimeEventController() *TimeEventController {
	return room.timeEventController
}

//Init 初始化GameRoom的参数
func (room *NormalGameRoom) Init(info *game_room.RoomInitInfo) {
	room.Id = info.Id
	room.port = info.Port
	room.over = info.Over
	room.Prepare = make(chan interface{})
	room.PlayerNum = 0
	room.Status = configs.NormalGameInitStatus
	room.heroManager = NewHeroesManager()
	//room.propsManager
	//room.bulletsManager
	room.netServer = gn.NewNormalGameNetServer()
	room.netServer.Init(room.port, configs.KcpRecvSize, configs.KcpSendSize)
	room.playersManager = NewPlayersManager()
	for _, pi := range info.JoinPlayers {
		room.playersManager.AddPlayer(gt.NewPlayer(pi.PlayerId))
	}
	room.requestController = NewRequestController()
	room.timeEventController = NewTimeEventController()
}

//Start 开始游戏,先开启kcp服务器,接收玩家的PlayerEnterGameRequest请求,根据init传过来的JoinPlayers进行身份校验,
//把PlayerManager和NormalGameNetServer的Player部分进行初始化,初始化playerId,heroId(可以参考学长的uuid进行分配),
//和网络连接的映射关系
func (room *NormalGameRoom) Start() {
	room.Status = configs.NormalGameWaitPlayerStatus
	room.netServer.Start()
	<-room.Prepare
	room.Status = configs.NormalGameStartStatus
	<-room.leave
	room.Status = configs.NormalGameGameDestroyingStatus
	//todo 清理工作/持久化数据给数据库等
}

//GetGameInfo 获取房间的游戏快照信息
func (room *NormalGameRoom) GetGameInfo() *game_room.GameInfo {
	return nil
}

//GetRoomInfo 获取游戏房间信息
func (room *NormalGameRoom) GetRoomInfo() *game_room.RoomInfo {
	return nil
}
