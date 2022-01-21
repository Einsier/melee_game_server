package game_room

import (
	pb "melee_game_server/api/proto"
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/game_net/api"
	"melee_game_server/internal/normal_game/codec"
	"melee_game_server/logger"
	"sync/atomic"
)

/**
*@Author Sly
*@Date 2022/1/19
*@Version 1.0
*@Description:
 */

//PlayerEnterGameRequestCallback 为请求的玩家申请一个heroId发给玩家,测试已完成(可见TestEnterGameRequest)
func PlayerEnterGameRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	resp := pb.PlayerEnterGameResponse{HeroId: int32(-1)} //如果是错误请求,返回heroId为-1

	if room.Status == int32(configs.NormalGameWaitPlayerStatus) {
		req := msg.Msg.Request.PlayerEnterGameRequest
		rId := req.GameRoomId
		pId := req.PlayerId
		pm := room.GetPlayerManager()
		net := room.GetNetServer()
		hm := room.GetHeroesManager()

		if rId != room.Id {
			logger.Errorf("[PlayerEnterGameRequestCallback]出现了和房间id不匹配的请求加入房间!playerId:%d\n", pId)
			net.SendBySingleConn(msg.Conn, codec.Encode(&resp))
			return
		}

		//查看是不是合法用户,如果不合法,发包告诉前端注册失败
		ok := pm.IsPlayerInRoom(pId)
		if !ok {
			logger.Errorf("[PlayerEnterGameRequestCallback]出现了不该出现的玩家请求加入房间!playerId:%d\n", pId)
			net.SendBySingleConn(msg.Conn, codec.Encode(&resp))
			return
		}

		//查看是不是已经注册了,如果已经注册了,重复注册,返回已经注册过的heroId,这里应该加锁...防止同一个用户反复注册
		pm.RegisterLock.Lock()
		ok = pm.IsPlayerRegistered(pId)
		if ok {
			pm.RegisterLock.Unlock()
			hId := pm.GetPlayer(pId).HeroId
			logger.Errorf("[PlayerEnterGameRequestCallback]前端重复发包注册,playerId:%d\n", pId)
			resp.HeroId = hId
			net.SendBySingleConn(msg.Conn, codec.Encode(&resp))
			return
		}
		hId := hm.ArrangeHeroId()

		if hId == int32(-1) {
			pm.RegisterLock.Unlock()
			logger.Errorf("[PlayerEnterGameRequestCallback]Player满之后又有玩家请求加入!playerId:%d\n", msg.Msg.Request.PlayerEnterGameRequest.PlayerId)
			//如果当前房间里的人已经满了,则返回heroId为-1
			net.SendBySingleConn(msg.Conn, codec.Encode(&resp))
			return
		}

		//如果从HeroManager拿到的id不是-1,说明分到了一个heroId,此时可以将player进行注册
		net.Register(pId, hId, msg.Conn)
		pm.GetPlayer(pId).BindConn(msg.Conn) //为player绑定conn
		pm.GetPlayer(pId).BindHeroId(hId)    //为player绑定heroId
		pm.GetPlayer(pId).SetStatus(configs.PlayerEnterGameStatus)
		pm.RegisterLock.Unlock() //注册完再解锁,这样下次同一个玩家重复注册无论什么顺序都会检测到
		logger.Infof("[PlayerEnterGameRequestCallback]完成玩家id:%d的玩家的注册,其heroId为:%d\n", pId, hId)
		atomic.AddInt32(&room.PlayerNum, 1) //GameRoom记录的玩家人数+1
		resp.HeroId = hId
		net.SendBySingleConn(msg.Conn, codec.Encode(&resp))

		//最后一个加入对局的请求的处理函数负责通知game_room所有玩家都准备完毕
		if hId == configs.MaxNormalGamePlayerNum {
			close(room.Prepare)
		}
	}

	logger.Errorf("[PlayerEnterGameRequestCallback]游戏开始后有有玩家请求加入!playerId:%d\n", msg.Msg.Request.PlayerEnterGameRequest.PlayerId)
	room.GetNetServer().SendBySingleConn(msg.Conn, codec.Encode(&resp))
}

//PlayerQuitGameRequestCallback 玩家发来退出游戏请求时的相应,只是把当前游戏人数-1,不做hero有关的调整,不会删除hero todo:待测试
func PlayerQuitGameRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	resp := pb.PlayerQuitGameResponse{Success: true} //如果是错误请求,返回heroId为-1
	if room.Status == configs.NormalGameStartStatus {
		req := msg.Msg.Request.PlayerQuitGameRequest
		pid := req.PlayerId
		pm := room.GetPlayerManager()

		pm.LeaveLock.Lock()
		//获取用户信息状态,如果状态为已离开,不做处理,如果没有离开,则状态设为已离开并且判断是否所有玩家都离开
		status := pm.GetPlayer(pid).GetStatus()
		if status == configs.PlayerLeaveGameStatus {
			pm.LeaveLock.Unlock()
			logger.Errorf("重复收到玩家离开消息,playerId:%d", pid)
			room.GetNetServer().SendBySingleConn(msg.Conn, codec.Encode(&resp))
			return
		}

		pm.GetPlayer(pid).SetStatus(configs.PlayerLeaveGameStatus)
		pm.LeaveLock.Unlock()
		room.GetNetServer().SendBySingleConn(msg.Conn, codec.Encode(&resp))
		n := atomic.AddInt32(&room.PlayerNum, -1)
		if n == 0 {
			//如果所有玩家都已经退出,则通知game_room
			close(room.leave)
		}
	}
}

func HeroPositionReportRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	//if room.Status ==
}

func HeroMovementChangeRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	//req := msg.Msg.Request.HeroMovementChangeRequest

}

func HeroAttackRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	//req := msg.Msg.Request.HeroAttackRequest
}

func HeroGetPropRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	//req := msg.Msg.Request.HeroGetPropRequest
}

func HeroBulletColliderHeroRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	//req := msg.Msg.Request.HeroBulletColliderHeroRequest
}

func PlayerHeartBeatRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	//req := msg.Msg.Request.PlayerHeartBeatRequest
}
