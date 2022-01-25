package simu_client

import (
	"fmt"
	"math/rand"
	"melee_game_server/api/proto"
	"melee_game_server/internal/normal_game/codec"
	"melee_game_server/plugins/kcp_net/adapter"
	"net"
	"time"

	"github.com/xtaci/kcp-go"
)

/**
*@Author chenjiajia
*@Date 2022/1/24
*@Version 1.0
*@Description: 模拟客户端
 */

type SimuClient struct {
	conn             net.Conn                //与服务器的连接
	playerId         int32                   //玩家id
	roomId           int32                   //游戏房间id
	heroId           int32                   //英雄id
	heroMovementType *proto.HeroMovementType //英雄移动状态
	position         *proto.Vector2          //英雄位置
	bulletIdByHero   int32                   //该英雄发送的第几颗子弹
}

func NewSimuClient(playerId, roomId int32, port string) *SimuClient {
	conn, err := kcp.DialWithOptions("localhost:"+port, nil, 0, 0)
	conn.SetNoDelay(1, 10, 2, 1)
	if err != nil {
		fmt.Printf("连接服务器失败，错误信息为 %v", err)
	}
	return &SimuClient{
		conn:           conn,
		playerId:       playerId,
		roomId:         roomId,
		position:       &proto.Vector2{X: 0.0, Y: 0.0},
		bulletIdByHero: 0,
	}
}

func (client *SimuClient) Start() {
	client.PlayerEnterGame()
	go client.ReceiveHandle()
	// 定时执行
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			client.HeroPositionReport()
			client.PlayerHeartBeat()
		}
	}()
}

// PlayerEnterGame 玩家进入游戏
func (client *SimuClient) PlayerEnterGame() {
	// 封装请求
	req := codec.EncodeRequest(&proto.PlayerEnterGameRequest{
		PlayerId:   client.playerId,
		GameRoomId: client.roomId,
	})
	adapter.Send(&client.conn, req)
}

// PlayerQuitGame 玩家退出游戏
func (client *SimuClient) PlayerQuitGame() {
	req := codec.EncodeRequest(&proto.PlayerQuitGameRequest{
		HeroId:   client.heroId,
		PlayerId: client.playerId,
	})
	adapter.Send(&client.conn, req)
}

// HeroPositionReport 英雄位置同步
func (client *SimuClient) HeroPositionReport() {
	req := codec.EncodeRequest(&proto.HeroPositionReportRequest{
		HeroId:           client.heroId,
		HeroMovementType: *client.heroMovementType,
		Position:         client.position,
		Time:             time.Now().UnixNano(),
	})
	adapter.Send(&client.conn, req)
}

//changeHeroPositionByRandom 随机改变英雄位置
func (client *SimuClient) changeHeroPositionByRandom() {
	client.position = &proto.Vector2{X: randomFloat(0, 1), Y: randomFloat(0, 1)}
}

//randomFloat 获取随机浮点数
func randomFloat(min, max float32) float32 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float32()*(max-min)
}

//HeroMovementChange 英雄移动状态改变
func (client *SimuClient) HeroMovementChange(heroMovementType *proto.HeroMovementType) {
	client.heroMovementType = heroMovementType
	client.changeHeroPositionByRandom()
	req := codec.EncodeRequest(&proto.HeroMovementChangeRequest{
		HeroId:           client.heroId,
		HeroMovementType: *client.heroMovementType,
		Position:         client.position,
		Time:             time.Now().UnixNano(),
	})
	adapter.Send(&client.conn, req)
}

//HeroSwordAttack 英雄近战攻击
func (client *SimuClient) HeroSwordAttack() {
	req := codec.EncodeRequest(&proto.HeroSwordAttackRequest{
		//todo
	})
	adapter.Send(&client.conn, req)
}

//HeroBulletLaunch 英雄发射子弹
func (client *SimuClient) HeroBulletLaunch(direction *proto.Vector2) {
	client.bulletIdByHero++
	client.changeHeroPositionByRandom()
	req := codec.EncodeRequest(&proto.HeroBulletLaunchRequest{
		HeroId:         client.heroId,
		Position:       client.position,
		Direction:      direction,
		BulletIdByHero: client.bulletIdByHero,
	})
	adapter.Send(&client.conn, req)
}

//HeroGetProp 英雄获取道具
func (client *SimuClient) HeroGetProp() {
	req := codec.EncodeRequest(&proto.HeroGetPropRequest{
		HeroId: client.heroId,
		//todo
	})
	adapter.Send(&client.conn, req)
}

//HeroBulletColliderHero 英雄被子弹打中
func (client *SimuClient) HeroBulletColliderHero() {
	req := codec.EncodeRequest(&proto.HeroBulletColliderHeroRequest{
		//todo
	})
	adapter.Send(&client.conn, req)
}

//PlayerHeartBeat 玩家心跳检测
func (client *SimuClient) PlayerHeartBeat() {
	req := codec.EncodeRequest(&proto.PlayerHeartBeatRequest{
		PlayerId:       client.playerId,
		ClientSendTime: time.Now().UnixNano(),
	})
	adapter.Send(&client.conn, req)
}

// ReceiveHandle 处理服务器发来的消息
func (client *SimuClient) ReceiveHandle() {
	for {
		msg := adapter.Receive(&client.conn)
		fmt.Printf("接收到消息 %v", msg)
		// response
		if msg.TopMessageType == proto.TopMessageType_ResponseType {
			switch msg.Response.ResponseCode {
			case proto.ResponseCode_PlayerEnterGameResponseCode: //玩家进入游戏
				heroId := msg.Response.PlayerEnterGameResponse.HeroId
				if heroId == -1 {
					fmt.Printf("玩家%d进入游戏请求失败\n", client.playerId)
					break
				}
				fmt.Printf("玩家%d分配的heroId为%d\n", client.playerId, heroId)
				client.heroId = msg.Response.PlayerEnterGameResponse.HeroId //获取heroId
			case proto.ResponseCode_PlayerQuitGameResponseCode: //玩家退出游戏
				isSuccess := msg.Response.PlayerQuitGameResponse.Success
				if !isSuccess {
					fmt.Printf("玩家%d退出游戏请求失败\n", client.playerId)
					break
				}
				client.conn.Close()
			case proto.ResponseCode_HeroGetPropResponseCode: //英雄获取道具
				//tode
			case proto.ResponseCode_PlayerHeartBeatResponseCode: //玩家心跳检测
				st := msg.Response.PlayerHeartBeatResponse.ServerSendTime
				ct := time.Now().UnixNano()
				fmt.Printf("玩家%d心跳检测成功，时延为%dms\n", client.playerId, (ct-st)/1000)
			}
		}
		// broadcast
		if msg.TopMessageType == proto.TopMessageType_BroadcastType {
			switch msg.Broadcast.BroadcastCode {
			case proto.BroadcastCode_HeroPositionReportBroadcastCode: //英雄位置同步
				res := msg.Broadcast.HeroPositionReportBroadcast
				fmt.Printf("英雄%d的当前位置为%v", res.HeroId, res.HeroPosition)
			case proto.BroadcastCode_HeroMovementChangeBroadcastCode: //英雄移动状态改变
				res := msg.Broadcast.HeroMovementChangeBroadcast
				fmt.Printf("英雄%d的当前位置为%v，移动状态为%v\n", res.HeroId, res.HeroPosition, res.HeroMovementType)
			case proto.BroadcastCode_HeroBulletLaunchBroadcastCode: //英雄发射子弹
				res := msg.Broadcast.HeroBulletLaunchBroadcast
				fmt.Printf("子弹%d于%d在位置%v发射，方向为%v\n", res.BulletId, res.Time, res.Position, res.Direction)
			case proto.BroadcastCode_HeroChangeHealthBroadcastCode: //英雄血量变化
				//todo
			case proto.BroadcastCode_HeroDeadBroadcastCode: //英雄死亡
				//todo
			}
		}
	}
}
