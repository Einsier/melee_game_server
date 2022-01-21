package test

import (
	"melee_game_server/api/proto"
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/internal/normal_game/codec"
	gr "melee_game_server/internal/normal_game/game_room"
	"sync"
	"testing"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/20
*@Version 1.0
*@Description:
 */

func TestArrangeHeroId(t *testing.T) {
	pm := gr.NewHeroesManager()
	m := make(map[int32]int32)
	lock := sync.Mutex{}
	for i := 0; i < 1000; i++ {
		go func() {
			id := pm.ArrangeHeroId()
			lock.Lock()
			defer lock.Unlock()
			if _, ok := m[id]; ok {
				m[id]++
			} else {
				m[id] = 1
			}
		}()
	}
	time.Sleep(1 * time.Second)
	for i := int32(1); i <= configs.MaxNormalGamePlayerNum; i++ {
		if m[i] != 1 {
			t.Fatalf("id:%d not arrange\n", i)
		}
	}
	if m[int32(-1)] != 1000-configs.MaxNormalGamePlayerNum {
		t.Fatalf("not arrange enough -1\n")
	}
}

func TestEncode(t *testing.T) {
	resp1 := &proto.PlayerEnterGameResponse{HeroId: -1}
	msg1 := codec.Encode(resp1)
	if msg1.TopMessageType != proto.TopMessageType_ResponseType || msg1.Response.PlayerEnterGameResponse.HeroId != -1 {
		t.Fatalf("PlayerEnterGameResponse encode wrong")
	}

	resp2 := &proto.PlayerQuitGameResponse{Success: true}
	msg2 := codec.Encode(resp2)
	if msg2.TopMessageType != proto.TopMessageType_ResponseType || msg2.Response.PlayerQuitGameResponse.Success != true {
		t.Fatalf("PlayerQuitGameResponse encode wrong")
	}

	resp3 := &proto.HeroGetPropResponse{Success: true}
	msg3 := codec.Encode(resp3)
	if msg3.TopMessageType != proto.TopMessageType_ResponseType || msg3.Response.HeroGetPropResponse.Success != true {
		t.Fatalf("HeroGetPropResponse encode wrong")
	}

	resp4 := &proto.PlayerHeartBeatResponse{ServerSendTime: 1000}
	msg4 := codec.Encode(resp4)
	if msg4.TopMessageType != proto.TopMessageType_ResponseType || msg4.Response.PlayerHeartBeatResponse.ServerSendTime != 1000 {
		t.Fatalf("PlayerHeartBeatResponse encode wrong")
	}

	//broadcast
	b1 := &proto.HeroChangeHealthBroadcast{HeroId: 100}
	msg5 := codec.Encode(b1)
	if msg5.TopMessageType != proto.TopMessageType_BroadcastType || msg5.Broadcast.HeroChangeHealthBroadcast.HeroId != 100 {
		t.Fatalf("HeroChangeHealthBroadcast encode wrong")
	}
	b2 := &proto.HeroMovementChangeBroadcast{HeroId: 1000}
	msg6 := codec.Encode(b2)
	if msg6.TopMessageType != proto.TopMessageType_BroadcastType || msg6.Broadcast.HeroMovementChangeBroadcast.HeroId != 1000 {
		t.Fatalf("HeroMovementChangeBroadcast encode wrong")
	}
	b3 := &proto.HeroPositionReportBroadcast{HeroId: 10000}
	msg7 := codec.Encode(b3)
	if msg7.TopMessageType != proto.TopMessageType_BroadcastType || msg7.Broadcast.HeroPositionReportBroadcast.HeroId != 10000 {
		t.Fatalf("HeroPositionReportBroadcast encode wrong")
	}
	b4 := &proto.HeroPropDeleteBroadcast{PropId: 5}
	msg8 := codec.Encode(b4)
	if msg8.TopMessageType != proto.TopMessageType_BroadcastType || msg8.Broadcast.HeroPropDeleteBroadcast.PropId != 5 {
		t.Fatalf("HeroPropDeleteBroadcast encode wrong")
	}
	b5 := &proto.HeroPropAddBroadcast{PropId: 1010}
	msg9 := codec.Encode(b5)
	if msg9.TopMessageType != proto.TopMessageType_BroadcastType || msg9.Broadcast.HeroPropAddBroadcast.PropId != 1010 {
		t.Fatalf("HeroPropAddBroadcast encode wrong")
	}
	b6 := &proto.HeroAttackBroadcast{HeroId: 2000}
	msg10 := codec.Encode(b6)
	if msg10.TopMessageType != proto.TopMessageType_BroadcastType || msg10.Broadcast.HeroAttackBroadcast.HeroId != 2000 {
		t.Fatalf("HeroAttackBroadcast encode wrong")
	}
	b7 := &proto.GameStartBroadcast{HeroId: 200}
	msg11 := codec.Encode(b7)
	if msg11.TopMessageType != proto.TopMessageType_BroadcastType || msg11.Broadcast.GameStartBroadcast.HeroId != 200 {
		t.Fatalf("GameStartBroadcast encode wrong")
	}
}
