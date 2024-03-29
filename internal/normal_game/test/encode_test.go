package test

import (
	"melee_game_server/api/client/proto"
	"melee_game_server/internal/normal_game/codec"
	"testing"
)

/**
*@Author Sly
*@Date 2022/1/25
*@Version 1.0
*@Description:测试codec/encode
 */
func TestEncode(t *testing.T) {
	resp1 := &proto.PlayerEnterGameResponse{HeroId: -1}
	msg1 := codec.Encode(resp1)
	if msg1.TopMessageType != proto.TopMessageType_ResponseType || msg1.Response.ResponseCode != proto.ResponseCode_PlayerEnterGameResponseCode ||
		msg1.Response.PlayerEnterGameResponse.HeroId != -1 {
		t.Fatalf("PlayerEnterGameResponse encode wrong")
	}

	resp2 := &proto.PlayerQuitGameResponse{Success: true}
	msg2 := codec.Encode(resp2)
	if msg2.TopMessageType != proto.TopMessageType_ResponseType || msg2.Response.ResponseCode != proto.ResponseCode_PlayerQuitGameResponseCode ||
		msg2.Response.PlayerQuitGameResponse.Success != true {
		t.Fatalf("PlayerQuitGameResponse encode wrong")
	}

	resp3 := &proto.HeroGetPropResponse{Success: true}
	msg3 := codec.Encode(resp3)
	if msg3.TopMessageType != proto.TopMessageType_ResponseType || msg3.Response.ResponseCode != proto.ResponseCode_HeroGetPropResponseCode ||
		msg3.Response.HeroGetPropResponse.Success != true {
		t.Fatalf("HeroGetPropResponse encode wrong")
	}

	resp4 := &proto.PlayerHeartBeatResponse{HeartbeatId: 1000}
	msg4 := codec.Encode(resp4)
	if msg4.TopMessageType != proto.TopMessageType_ResponseType || msg4.Response.ResponseCode != proto.ResponseCode_PlayerHeartBeatResponseCode ||
		msg4.Response.PlayerHeartBeatResponse.HeartbeatId != 1000 {
		t.Fatalf("PlayerHeartBeatResponse encode wrong")
	}

	//broadcast
	b1 := &proto.HeroChangeHealthBroadcast{HeroId: 100}
	msg5 := codec.Encode(b1)
	if msg5.TopMessageType != proto.TopMessageType_BroadcastType || msg5.Broadcast.BroadcastCode != proto.BroadcastCode_HeroChangeHealthBroadcastCode ||
		msg5.Broadcast.HeroChangeHealthBroadcast.HeroId != 100 {
		t.Fatalf("HeroChangeHealthBroadcast encode wrong")
	}
	b2 := &proto.HeroMovementChangeBroadcast{HeroId: 1000}
	msg6 := codec.Encode(b2)
	if msg6.TopMessageType != proto.TopMessageType_BroadcastType || msg6.Broadcast.BroadcastCode != proto.BroadcastCode_HeroMovementChangeBroadcastCode ||
		msg6.Broadcast.HeroMovementChangeBroadcast.HeroId != 1000 {
		t.Fatalf("HeroMovementChangeBroadcast encode wrong")
	}
	b3 := &proto.HeroPositionReportBroadcast{HeroId: 10000}
	msg7 := codec.Encode(b3)
	if msg7.TopMessageType != proto.TopMessageType_BroadcastType || msg7.Broadcast.BroadcastCode != proto.BroadcastCode_HeroPositionReportBroadcastCode ||
		msg7.Broadcast.HeroPositionReportBroadcast.HeroId != 10000 {
		t.Fatalf("HeroPositionReportBroadcast encode wrong")
	}
	b4 := &proto.HeroPropDeleteBroadcast{PropId: 5}
	msg8 := codec.Encode(b4)
	if msg8.TopMessageType != proto.TopMessageType_BroadcastType || msg8.Broadcast.BroadcastCode != proto.BroadcastCode_HeroPropDeleteBroadcastCode ||
		msg8.Broadcast.HeroPropDeleteBroadcast.PropId != 5 {
		t.Fatalf("HeroPropDeleteBroadcast encode wrong")
	}
	b5 := &proto.HeroPropAddBroadcast{PropId: 1010}
	msg9 := codec.Encode(b5)
	if msg9.TopMessageType != proto.TopMessageType_BroadcastType || msg9.Broadcast.BroadcastCode != proto.BroadcastCode_HeroPropAddBroadcastCode ||
		msg9.Broadcast.HeroPropAddBroadcast.PropId != 1010 {
		t.Fatalf("HeroPropAddBroadcast encode wrong")
	}
	b6 := &proto.HeroSwordAttackBroadcast{HeroId: 2000}
	msg10 := codec.Encode(b6)
	if msg10.TopMessageType != proto.TopMessageType_BroadcastType || msg10.Broadcast.BroadcastCode != proto.BroadcastCode_HeroSwordAttackBroadcastCode ||
		msg10.Broadcast.HeroSwordAttackBroadcast.HeroId != 2000 {
		t.Fatalf("HeroAttackBroadcast encode wrong")
	}
	b7 := &proto.GameStartBroadcast{NickNameMap: map[int32]string{1: "abc"}}
	msg11 := codec.Encode(b7)
	if msg11.TopMessageType != proto.TopMessageType_BroadcastType || msg11.Broadcast.BroadcastCode != proto.BroadcastCode_GameStartBroadcastCode ||
		msg11.Broadcast.GameStartBroadcast.NickNameMap == nil {
		t.Fatalf("GameStartBroadcast encode wrong")
	}
	b8 := &proto.HeroBulletLaunchBroadcast{BulletId: 1200}
	msg12 := codec.Encode(b8)
	if msg11.TopMessageType != proto.TopMessageType_BroadcastType || msg12.Broadcast.BroadcastCode != proto.BroadcastCode_HeroBulletLaunchBroadcastCode ||
		msg12.Broadcast.HeroBulletLaunchBroadcast.BulletId != 1200 {
		t.Fatalf("HeroBulletLaunchBroadcast encode wrong")
	}
	b9 := &proto.HeroDeadBroadcast{HeroId: 12100}
	msg13 := codec.Encode(b9)
	if msg13.TopMessageType != proto.TopMessageType_BroadcastType || msg13.Broadcast.BroadcastCode != proto.BroadcastCode_HeroDeadBroadcastCode ||
		msg13.Broadcast.HeroDeadBroadcast.HeroId != 12100 {
		t.Fatalf("HeroDeadBroadcast encode wrong")
	}
	b10 := &proto.HeroBulletDestroyBroadcast{BulletId: 12010}
	msg14 := codec.Encode(b10)
	if msg14.TopMessageType != proto.TopMessageType_BroadcastType || msg14.Broadcast.BroadcastCode != proto.BroadcastCode_HeroBulletDestroyBroadcastCode ||
		msg14.Broadcast.HeroBulletDestroyBroadcast.BulletId != 12010 {
		t.Fatalf("HeroBulletDestroyBroadcast encode wrong")
	}
}
