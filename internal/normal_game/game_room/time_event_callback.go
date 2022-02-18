package game_room

import (
	pb "melee_game_server/api/client/proto"
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/internal/normal_game/codec"
	gt "melee_game_server/internal/normal_game/game_type"
	"melee_game_server/plugins/logger"
	"melee_game_server/utils"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/24
*@Version 1.0
*@Description:用于注册定时事件
 */

var CleanOverTimeBulletTimeEvent = NewTimeEvent(CleanOverTimeBulletTimeEventCallbackId, time.Nanosecond*CleanOverTimeBulletTimeEventCallbackTimeSlice, CleanOverTimeBulletTimeEventCallback)
var RefreshPropsTimeEvent = NewTimeEvent(RefreshPropsTimeEventCallbackId, time.Nanosecond*RefreshPropsTimeEventCallbackTimeSlice, RefreshPropsTimeEventCallback)

const (
	CleanOverTimeBulletTimeEventCallbackId = iota
	RefreshPropsTimeEventCallbackId
)

//注意单位均为ns
const (
	CleanOverTimeBulletTimeEventCallbackTimeSlice = configs.BulletDuration
	RefreshPropsTimeEventCallbackTimeSlice        = configs.PropRefreshTime
)

//CleanOverTimeBulletTimeEventCallback 定期清理过期的子弹,客户端子弹如果没撞到人过2s由unity自动清除(射程有限),服务器子弹过2s由服务器清除
func CleanOverTimeBulletTimeEventCallback(room *NormalGameRoom) {
	bm := room.GetBulletsManager()
	bm.refreshLock.Lock()
	deleteBullet := bm.oldBullets
	bm.oldBullets = bm.newBullets
	bm.newBullets = make([]int64, 0)
	bm.refreshLock.Unlock()

	for _, id := range deleteBullet {
		//todo 测试用,待删除
		hid, bid := countHidBid(id)
		t := time.Now().Format(time.StampNano)
		logger.Testf("[%s]删除了hero[%d]的第[%d]颗子弹", t, hid, bid)
		bm.bullets.Delete(id)
	}
}

//RefreshPropsTimeEventCallback 定期刷新道具
func RefreshPropsTimeEventCallback(room *NormalGameRoom) {
	pm := room.GetPropsManager()
	for i := 0; i < configs.PropRefreshNumPerTime; i++ {
		X1 := utils.RandomFloat64(0, configs.NormalGameMapWidth)
		X2 := utils.RandomFloat64(0, configs.NormalGameMapWidth)
		Y1 := utils.RandomFloat64(0, configs.NormalGameMapHeight)
		Y2 := utils.RandomFloat64(0, configs.NormalGameMapHeight)
		id1 := pm.GetId()
		id2 := pm.GetId()
		p1 := gt.NewProp(id1, pb.PropType_StrawberryType, gt.NewVector2(X1, Y1))
		p2 := gt.NewProp(id2, pb.PropType_BulletPocketType, gt.NewVector2(X2, Y2))
		pm.AddProp(p1)
		pm.AddProp(p2)

		m1 := codec.Encode(&pb.HeroPropAddBroadcast{
			PropId:       id1,
			PropPosition: &pb.Vector2{X: float32(X1), Y: float32(Y1)},
			PropType:     pb.PropType_StrawberryType,
		})
		m2 := codec.Encode(&pb.HeroPropAddBroadcast{
			PropId:       id2,
			PropPosition: &pb.Vector2{X: float32(X2), Y: float32(Y2)},
			PropType:     pb.PropType_BulletPocketType,
		})
		room.GetNetServer().SendToAllPlayerConn(m1)
		room.GetNetServer().SendToAllPlayerConn(m2)
	}
}
