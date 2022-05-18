package game_room

import (
	"melee_game_server/api/client/proto"
	"melee_game_server/api/hall"
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/game_net/api"
	"melee_game_server/framework/game_room"
	"melee_game_server/internal/normal_game/aoi"
	"melee_game_server/internal/normal_game/codec"
	gn "melee_game_server/internal/normal_game/game_net"
	gt "melee_game_server/internal/normal_game/game_type"
	"melee_game_server/internal/normal_game/metrics"
	"melee_game_server/plugins/logger"
	"strconv"
	"sync/atomic"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:用于管理一次游戏的对局,包括玩家管理,对局中的prop(道具)管理以及子弹管理
 */

type NormalGameRoom struct {
	Id int32
	//port                string
	Prepare         chan interface{} //所有玩家都已连入游戏
	leave           chan interface{} //玩家都已经离开
	over            chan struct{}    //用于向game_server汇报的
	TestRequestChan chan *api.Mail   //用于测试的chan,把Mail传进去交给RequestController处理
	PlayerNum       int32            //当前的存活的玩家数目
	Status          int32            //当前的状态

	StartTime time.Time
	EndTime   time.Time

	heroManager         *HeroesManager
	propsManager        *PropsManager
	bulletsManager      *BulletsManager
	netServer           *gn.NormalGameNetServer
	playersManager      *PlayersManager
	requestController   *RequestController
	timeEventController *TimeEventController
	honorManager        *HonorManager
	aoi                 *aoi.AOI
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
	room.over = info.Over
	room.leave = make(chan interface{})
	room.Prepare = make(chan interface{})
	room.PlayerNum = 0
	room.Status = configs.NormalGameInitStatus
	//todo 删除掉这行
	room.TestRequestChan = make(chan *api.Mail)
	room.heroManager = NewHeroesManager()
	room.propsManager = NewPropsManager()
	room.bulletsManager = NewBulletsManager()
	room.netServer = gn.NewNormalGameNetServer(info.Id)
	room.playersManager = NewPlayersManager()
	room.honorManager = NewHonorManager()
	for _, pi := range info.JoinPlayers {
		room.playersManager.AddPlayer(gt.NewPlayer(pi))
		room.honorManager.AddPlayerHonor(pi.PlayerId)
	}
	room.requestController = NewRequestController()
	room.timeEventController = NewTimeEventController(room)
	logger.Infof("房间%d初始化完毕", room.Id)
}

//Start 开始游戏,先开启kcp服务器,接收玩家的PlayerEnterGameRequest请求,根据init传过来的JoinPlayers进行身份校验,
//把PlayerManager和NormalGameNetServer的Player部分进行初始化,初始化playerId,heroId,记录网络和玩家的映射消息,
//注意应该使用go Start()进行调用,然后外面监听room.over
func (room *NormalGameRoom) Start() {
	room.Status = configs.NormalGameWaitPlayerStatus
	//room.netServer.Start()x
	//go room.requestController.Work1(room)
	go room.requestController.Work2(room)
	<-room.Prepare
	logger.Infof("room:%d 所有玩家准备就绪,开始游戏", room.Id)
	//代码执行到这里,所有的玩家都已经准备好
	//todo 将测试地图改成真实的地图
	/*	room.aoi = aoi.NewAOI(aoi.NewRandomHeroesInitInfo(configs.MaxNormalGamePlayerNum, aoi.TestHeroSpeed, aoi.TestMapQT),
		aoi.TestGameMapWidth, aoi.TestGameMapHeight, aoi.TestGridWidth, aoi.TestGridHeight, 100*time.Millisecond, room.GetNetServer(), aoi.TestMapQT)*/
	room.aoi = aoi.NewAOI(aoi.NewRandomHeroesInitInfo(configs.MaxNormalGamePlayerNum, configs.HeroMoveSpeed, aoi.NormalGameMapQt),
		configs.MapWidth, configs.MapHeight, configs.GridWidth, configs.GridHeight, configs.FrameSyncSlice, room.GetNetServer(), aoi.NormalGameMapQt)

	time.Sleep(20 * time.Millisecond) //等待最后一个分配heroId的包到达
	//room.netServer.SendToAllPlayerConn(room.GetNormalGameStartBroadcastInfo()) //发消息通知所有的玩家游戏开始
	//todo 做开始的包的删减
	room.SendToAllPlayerInRoom(room.GetNormalGameStartBroadcastInfo())
	room.Status = configs.NormalGameStartStatus
	time.Sleep(300 * time.Millisecond)
	room.StartTime = time.Now()
	room.aoi.Work(room.StartTime)
	//注册定时事件
	room.GetTimeEventController().AddEvent(CheckHeartBeatTimeEvent)
	room.GetTimeEventController().AddEvent(CleanOverTimeBulletTimeEvent)
	//room.GetTimeEventController().AddEvent(RefreshPropsTimeEvent)

	//等待游戏结束
	<-room.leave
	time.Sleep(20 * time.Second)
	room.aoi.Stop()
	room.GetTimeEventController().Destroy()
	room.Status = configs.NormalGameGameDestroyingStatus
	room.EndTime = time.Now()
	logger.Infof("room %d 所有玩家已经离开,准备进行清理工作", room.Id)

	close(room.over) //通知本game_room已经结束
}

//GetNormalGameStartBroadcastInfo 为了方便前端显示名字以及其他内容,将每个英雄对应的玩家的名字加入集合发给前端
func (room *NormalGameRoom) GetNormalGameStartBroadcastInfo() *proto.TopMessage {
	pm := room.GetPlayerManager()
	m := make(map[int32]string) //key为玩家的heroId,value为对应玩家的nickname
	if configs.MaxNormalGamePlayerNum < 200 {
		for i := int32(1); i <= configs.MaxNormalGamePlayerNum; i++ {
			nickname := pm.GetNicknameByHeroId(i)
			if nickname == "" {
				logger.Errorf("获取hero id:%d nickname出错,游戏不能正常显示nickname!\n", i)
				nickname = "error msg"
			}
			m[i] = nickname
		}
	}
	return codec.Encode(&proto.GameStartBroadcast{NickNameMap: m, HeroNum: configs.MaxNormalGamePlayerNum})
}

//GetGameInfo 获取房间的游戏快照信息
func (room *NormalGameRoom) GetGameInfo() *game_room.GameInfo {
	return nil
}

//GetRoomInfo 获取游戏房间信息
func (room *NormalGameRoom) GetRoomInfo() *game_room.RoomInfo {
	return nil
}

func (room *NormalGameRoom) PutMsg(mail *api.Mail) {
	room.netServer.ReqChan <- mail
}

func (room *NormalGameRoom) CloseMsgChan() {
	close(room.netServer.ReqChan)
}

//ForceStopGame todo
func (room *NormalGameRoom) ForceStopGame() (ok bool) {
	return true
}

//GetGameAccount 获取游戏结算信息(为了接口的不能循环依赖,这里把返回值改成了interface类型)
func (room *NormalGameRoom) GetGameAccount() interface{} {
	gameAccount := new(hall.GameAccountInfo)
	gameAccount.StartTime = room.StartTime.UnixNano()
	gameAccount.EndTime = room.EndTime.UnixNano()
	gameAccount.PlayerAccountMap = room.honorManager.GetAllPlayerHonor()
	return gameAccount
}

//DeletePlayer 删除一个游戏内的玩家,包括断开其连接,在各种manager中进行删除等,成功返回true,不成功返回false
//删除英雄/玩家的唯一方式,有三种方式会导致玩家被删除:	1.玩家掉血死亡	2.玩家心跳过期	3.玩家自己退出
func (room *NormalGameRoom) DeletePlayer(pid int32) bool {
	pm := room.GetPlayerManager()
	//获取用户信息状态,如果状态为已离开,不做处理,如果没有离开,则状态设为已离开并且判断是否所有玩家都离开
	p := pm.GetPlayer(pid)
	pm.LeaveLock.Lock()
	if p == nil {
		pm.LeaveLock.Unlock()
		//logger.Errorf("出现不属于本房间的玩家发送退出请求!playerId:%d", pid)
		return false
	}
	status := p.GetStatus()
	if status == configs.PlayerLeaveGameStatus {
		pm.LeaveLock.Unlock()
		//logger.Errorf("重复收到玩家离开消息,playerId:%d", pid)
		return false
	}

	pm.GetPlayer(pid).SetStatus(configs.PlayerLeaveGameStatus)
	pm.LeaveLock.Unlock()
	before := atomic.LoadInt32(&room.PlayerNum)
	if before <= 1 {
		return false
	}
	n := atomic.AddInt32(&room.PlayerNum, -1)
	metrics.GaugeVecGameRoomPlayerCount.WithLabelValues(strconv.Itoa(int(room.Id))).Dec()
	hid := pm.GetHeroId(pid)
	room.aoi.RemoveHero(hid) //从aoi中删除该英雄
	//room.GetNetServer().DeleteConn(hid, pid)
	pm.DeletePlayer(pm.GetPlayer(pid))
	logger.Infof("room%d playerId为:%d,heroId为:%d的玩家已退出游戏,当前剩余%d人\n", room.Id, pid, hid, n)

	//更新玩家荣誉信息
	room.honorManager.GetPlayerHonor(pid).SetAliveTime(room.StartTime.UnixNano() - time.Now().UnixNano())
	if n <= 1 {
		//如果当前场上只有一个玩家了,那么这个玩家是最终的胜利者,广播游戏结束报文,并且通知上层游戏结束
		broad := &proto.GameOverBroadcast{}
		room.SendToAllPlayerInRoom(codec.Encode(broad))
		alivePlayer := room.GetAllPlayerIdInRoom()
		if len(alivePlayer) == 1 {
			logger.Infof("胜利者是:player%d", alivePlayer[0])
		}
		for _, alivePlayerId := range alivePlayer {
			//将剩余的那个玩家从aoi模块中删除,不同步消息了
			if aliveHeroId := room.GetPlayerManager().GetHeroId(alivePlayerId); aliveHeroId != 0 {
				room.aoi.RemoveHero(aliveHeroId)
			}
			//pm中删除
			pm.DeletePlayer(pm.GetPlayer(alivePlayerId))
			metrics.GaugeVecGameRoomPlayerCount.WithLabelValues(strconv.Itoa(int(room.Id))).Set(0)
		}
		close(room.leave)
	}
	return true
}

//GetAllPlayerIdInRoom 获取当前还在房间中的玩家的id
func (room *NormalGameRoom) GetAllPlayerIdInRoom() []int32 {
	pm := room.GetPlayerManager()
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	ret := make([]int32, len(pm.players))
	i := 0
	for pid := range pm.players {
		ret[i] = pid
		i++
	}
	return ret
}

func (room *NormalGameRoom) SendToAllPlayerInRoom(msg *proto.TopMessage) {
	alivePlayer := room.GetAllPlayerIdInRoom()
	room.GetNetServer().SendByPlayerId(alivePlayer, msg)
}
