package test

import (
	"melee_game_server/api/proto"
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/game_net/api"
	framework "melee_game_server/framework/game_room"
	"melee_game_server/internal/normal_game/codec"
	gr "melee_game_server/internal/normal_game/game_room"
	"melee_game_server/plugins/logger"
	"testing"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/21
*@Version 1.0
*@Description:
 */

const TestRoomId = 1

var TestJoinPlayers = make([]*framework.PlayerInfo, 0)

var TestEnterGameRequestRightMap = make(map[int32]*api.Mail)
var TestEnterGameRequestWrongMap = make(map[int32]*api.Mail)

//testInit 创建configs.MaxNormalGamePlayerNum的正确的登录请求填到TestEnterGameRequestRightMap里,
//创建3*configs.MaxNormalGamePlayerNum的测试用例填到TestEnterGameRequestWrongMap里
func testInit() {
	offset := int32(10000)
	for i := int32(0); i < configs.MaxNormalGamePlayerNum; i++ {
		pId := offset + i
		p := framework.PlayerInfo{PlayerId: pId}
		TestJoinPlayers = append(TestJoinPlayers, &p)
		rightReq := codec.EncodeRequest(&proto.PlayerEnterGameRequest{
			PlayerId:   pId,
			GameRoomId: TestRoomId,
		})
		wrongReq1 := codec.EncodeRequest(&proto.PlayerEnterGameRequest{
			PlayerId:   pId + 10000,
			GameRoomId: TestRoomId,
		})
		wrongReq2 := codec.EncodeRequest(&proto.PlayerEnterGameRequest{
			PlayerId:   pId + 20000,
			GameRoomId: TestRoomId,
		})
		wrongReq3 := codec.EncodeRequest(&proto.PlayerEnterGameRequest{
			PlayerId:   pId + 30000,
			GameRoomId: TestRoomId,
		})
		rightMail := api.Mail{
			Conn: nil,
			Msg:  rightReq,
		}
		wrongMail1 := api.Mail{
			Conn: nil,
			Msg:  wrongReq1,
		}
		wrongMail2 := api.Mail{
			Conn: nil,
			Msg:  wrongReq2,
		}
		wrongMail3 := api.Mail{
			Conn: nil,
			Msg:  wrongReq3,
		}
		TestEnterGameRequestRightMap[pId] = &rightMail
		TestEnterGameRequestWrongMap[pId+10000] = &wrongMail1
		TestEnterGameRequestWrongMap[pId+20000] = &wrongMail2
		TestEnterGameRequestWrongMap[pId+30000] = &wrongMail3
	}
}

func TestEnterGameRequest(t *testing.T) {
	logger.SetLogLevel(logger.LogDisabledLevel)
	testInit()
	room := new(gr.NormalGameRoom)
	var TestRoomInitInfo = framework.RoomInitInfo{
		Id:          1,
		Port:        "8000",
		Over:        make(chan interface{}),
		JoinPlayers: TestJoinPlayers,
	}
	room.Init(&TestRoomInitInfo)
	room.Status = configs.NormalGameWaitPlayerStatus
	//处理正确消息
	go func() {
		for _, mail := range TestEnterGameRequestRightMap {
			go gr.PlayerEnterGameRequestCallback(mail, room)
		}
	}()
	go func() {
		for _, mail := range TestEnterGameRequestWrongMap {
			go gr.PlayerEnterGameRequestCallback(mail, room)
		}
	}()
	time.Sleep(1 * time.Second)
	dup := make(map[int32]struct{})
	res := make(map[int32]struct{})
	for i := int32(1); i <= configs.MaxNormalGamePlayerNum; i++ {
		res[i] = struct{}{}
	}
	for _, v := range TestEnterGameRequestRightMap {
		pId := v.Msg.Request.PlayerEnterGameRequest.PlayerId
		hId := room.GetPlayerManager().GetPlayer(pId).HeroId
		if _, ok := dup[pId]; ok {
			t.Fatalf("pId:%d registered two times!", pId)
		}
		dup[pId] = struct{}{}
		delete(res, hId)
	}
	if len(res) != 0 {
		t.Fatalf("not all heroId distributed!,left:%v", res)
	}
}

/*func TestPlayerManager(t *testing.T) {
	p0 := gt.NewPlayer(0)
	p1 := gt.NewPlayer(1)
	p2 := gt.NewPlayer(2)
	p3 := gt.NewPlayer(3)
	p4 := gt.NewPlayer(4)
	pm := gr.NewPlayerManager()
	pm.AddPlayer(p0)
	pm.AddPlayer(p1)
	pm.AddPlayer(p2)
	pm.AddPlayer(p3)
	pm.AddPlayer(p4)
	for i := 0; i < 5; i++ {
		go func(id int32) {
			for k := 0; k < 1000; k++ {
				h := pm.GetPlayer(id)
				h.HeroId++
			}
		}(int32(i))
	}
	time.Sleep(1*time.Second)
	if p0.HeroId != 1000 || p1.HeroId != 1000||p2.HeroId != 1000||p3.HeroId != 1000||p4.HeroId != 1000{
		t.Fatalf("wrong")
	}
}*/
