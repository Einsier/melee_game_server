package game_room

import (
	pb "melee_game_server/api/client/proto"
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/entity"
	"melee_game_server/framework/game_net/api"
	"melee_game_server/internal/normal_game/codec"
	gt "melee_game_server/internal/normal_game/game_type"
	"melee_game_server/plugins/logger"
	"sync/atomic"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/19
*@Version 1.0
*@Description:
 */

//PlayerEnterGameRequestCallback 为请求的玩家申请一个heroId发给玩家,测试已完成(可见TestGameRequest)
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
		pm.hToP[hId] = pId
		pm.pToH[pId] = hId
		pm.RegisterLock.Unlock() //注册完再解锁,这样下次同一个玩家重复注册无论什么顺序都会检测到
		logger.Infof("room %d 完成玩家id:%d的玩家的注册,其heroId为:%d\n", room.Id, pId, hId)
		atomic.AddInt32(&room.PlayerNum, 1) //GameRoom记录的玩家人数+1
		resp.HeroId = hId
		net.SendBySingleConn(msg.Conn, codec.Encode(&resp))

		//最后一个加入对局的请求的处理函数负责通知game_room所有玩家都准备完毕
		if hId == configs.MaxNormalGamePlayerNum {
			time.Sleep(50 * time.Millisecond)
			close(room.Prepare)
		}
		return
	}

	logger.Errorf("[PlayerEnterGameRequestCallback]游戏开始后有有玩家请求加入!playerId:%d\n", msg.Msg.Request.PlayerEnterGameRequest.PlayerId)
	room.GetNetServer().SendBySingleConn(msg.Conn, codec.Encode(&resp))
}

//PlayerQuitGameRequestCallback 玩家发来退出游戏请求时的相应,只是把当前游戏人数-1,删除对应的hero和player的连接记录,不做hero有关的调整,不会删除hero,测试已完成(见TestGameRequest)
func PlayerQuitGameRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	resp := pb.PlayerQuitGameResponse{Success: true} //如果是错误请求,同样返回true
	req := msg.Msg.Request.PlayerQuitGameRequest
	if room.Status == configs.NormalGameStartStatus {
		ok := room.DeletePlayer(req.PlayerId)
		resp.Success = ok
		room.GetNetServer().SendBySingleConn(msg.Conn, codec.Encode(&resp))
		return
	}
	logger.Errorf("[PlayerQuitGameRequestCallback]在%d阶段收到了playerId:%d,heroId:%d,的QuitGame请求!\n", room.Status, req.PlayerId, req.HeroId)
}

func HeroPositionReportRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	req := msg.Msg.Request.HeroPositionReportRequest
	hid := req.HeroId
	if room.Status == configs.NormalGameStartStatus {
		X := req.Position.X
		Y := req.Position.Y

		room.heroManager.UpdateHeroPositionInfo(hid, entity.NewVector2(float64(X), float64(Y)), req.Time, req.HeroMovementType)
		broadHeroes := room.heroManager.GetHeroesNearby(hid)
		broad := codec.Encode(&pb.HeroPositionReportBroadcast{
			HeroId:       hid,
			HeroPosition: &pb.Vector2{X: X, Y: Y},
		})
		room.GetNetServer().SendByHeroId(broadHeroes, broad)
		return
	}
	logger.Errorf("[HeroPositionReportRequestCallback]在%d阶段收到了heroId:%d,的HeroPositionReportRequest!\n", room.Status, hid)
}

func HeroMovementChangeRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	req := msg.Msg.Request.HeroMovementChangeRequest
	hid := req.HeroId
	if room.Status == configs.NormalGameStartStatus {
		position := entity.NewVector2(float64(req.Position.X), float64(req.Position.Y))
		room.heroManager.UpdateHeroPositionInfo(hid, position, req.Time, req.HeroMovementType)
		broadHeroes := room.heroManager.GetHeroesNearby(hid)
		broad := codec.Encode(&pb.HeroMovementChangeBroadcast{
			HeroId:           hid,
			HeroMovementType: req.HeroMovementType,
			Time:             req.Time,
			HeroPosition:     req.Position,
		})
		room.GetNetServer().SendByHeroId(broadHeroes, broad)
		return
	}
	logger.Errorf("[HeroMovementChangeRequestCallback]在%d阶段收到了heroId:%d,的HeroPositionReportRequest!\n", room.Status, hid)
}

func HeroBulletLaunchRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	req := msg.Msg.Request.HeroBulletLaunchRequest
	hid := req.HeroId
	if room.Status == configs.NormalGameStartStatus {
		position := entity.Vector2{X: float64(req.Position.X), Y: float64(req.Position.Y)}
		direction := entity.Vector2{X: float64(req.Direction.X), Y: float64(req.Direction.Y)}
		bullet := gt.NewBullet(hid, req.BulletIdByHero, req.LaunchTime, position, direction)
		room.GetBulletsManager().AddBullets(bullet)
		broadHeroes := room.heroManager.GetHeroesNearby(hid)
		broad := codec.Encode(&pb.HeroBulletLaunchBroadcast{
			BulletId:  gt.CountBulletId(hid, req.BulletIdByHero),
			Position:  req.Position,
			Direction: req.Direction,
			Time:      req.LaunchTime,
		})
		room.GetNetServer().SendByHeroId(broadHeroes, broad)
		return
	}
	logger.Errorf("[HeroBulletLaunchRequestCallback]在%d阶段收到了heroId:%d,的HeroBulletLaunchRequest!\n", room.Status, hid)
}

func HeroSwordAttackRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	//todo:判断周围有没有英雄
}

func HeroGetPropRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	req := msg.Msg.Request.HeroGetPropRequest
	hid := req.HeroId
	pid := req.PropId

	if room.Status == configs.NormalGameStartStatus {
		//判断英雄有没有狗带,如果已经狗带则不能吃道具
		if room.GetHeroesManager().GetHeroStatus(hid) != configs.HeroAlive {
			topMsg := codec.Encode(&pb.HeroGetPropResponse{Success: false})
			room.GetNetServer().SendBySingleConn(msg.Conn, topMsg)
			logger.Errorf("[HeroGetPropRequestCallback]已经狗带的英雄heroId:%d要吃道具id:%d!\n", hid, pid)
			return
		}

		//初版后端英雄吃道具暂时不加校验,后期再加
		pType, ok := room.propsManager.EatProp(pid)
		if !ok {
			topMsg := codec.Encode(&pb.HeroGetPropResponse{Success: false})
			room.GetNetServer().SendBySingleConn(msg.Conn, topMsg)
			logger.Infof("[HeroGetPropRequestCallback]heroId:%d要吃道具id:%d失败!\n", hid, pid)
			return
		}

		//如果已经确定吃到,判断类型
		switch pType {
		case pb.PropType_StrawberryType:
			//如果是草莓,加血,判断有没有
			isChange, isDead, newHealth := room.GetHeroesManager().GetHero(hid).ChangeHeath(1)
			//万一如果isDead,说明刚刚刚被狗带了,吃个草莓也没用...狗带的消息不用自己发,所以return
			if isDead {
				return
			}
			if isChange {
				room.netServer.SendToAllPlayerConn(codec.Encode(
					&pb.HeroChangeHealthBroadcast{HeroId: hid, HeroHealth: newHealth}))
			}
		case pb.PropType_BulletPocketType:
			//暂定不做处理,后期加上防止前端作弊逻辑
		default:
			logger.Errorf("[HeroGetPropRequestCallback]收到了不在编号上的道具的请求,propType:%d", pType)
		}
		return
	}
	logger.Errorf("[HeroGetPropRequestCallback]在%d阶段收到了heroId:%d,的HeroGetPropRequest,要吃道具id:%d!\n", room.Status, hid, pid)
}

func HeroBulletColliderHeroRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	req := msg.Msg.Request.HeroBulletColliderHeroRequest
	hid := req.ColliderHeroId
	bid := req.BulletId
	if room.Status == configs.NormalGameStartStatus {
		if room.GetHeroesManager().GetHeroStatus(hid) != configs.HeroAlive {
			//如果射击到的英雄已经狗带了,那么不做任何的反应,也不会删除子弹,即子弹应该穿过狗带的英雄
			return
		}

		heroPosition := room.heroManager.GetHeroPosition(hid)
		ok := room.GetBulletsManager().CheckBulletHitHero(bid, heroPosition)
		if !ok {
			//如果当前射击无效,则应该直接返回,不删除子弹
			return
		}
		room.GetBulletsManager().DeleteBullets(bid)
		room.netServer.SendToAllPlayerConn(codec.Encode(&pb.HeroBulletDestroyBroadcast{BulletId: bid}))
		//被击中的玩家扣血
		isChange, isDead, newHealth := room.GetHeroesManager().GetHero(hid).ChangeHeath(-1)
		if isChange {
			room.netServer.SendToAllPlayerConn(codec.Encode(
				&pb.HeroChangeHealthBroadcast{HeroId: hid, HeroHealth: newHealth}))
		}
		if isDead {
			room.netServer.SendToAllPlayerConn(codec.Encode(
				&pb.HeroDeadBroadcast{HeroId: hid}))
		}
	}
	logger.Errorf("[HeroBulletColliderHeroRequestCallback]在%d阶段收到了heroId:%d,的HeroBulletColliderHeroRequest\n", room.Status, bid>>32)
}

func PlayerHeartBeatRequestCallback(msg *api.Mail, room *NormalGameRoom) {
	req := msg.Msg.Request.PlayerHeartBeatRequest
	pid := req.PlayerId
	ct := time.Now()
	room.GetPlayerManager().UpdateHeartBeatTime(pid, ct)
	room.netServer.SendBySingleConn(msg.Conn, codec.Encode(&pb.PlayerHeartBeatResponse{HeartbeatId: req.GetHeartBeatId()}))
}
