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

var TestEnterGameRequestRightSlice = make([]*api.Mail, 0)
var TestEnterGameRequestWrongSlice = make([]*api.Mail, 0)
var TestQuitGameRequestRightSlice = make([]*api.Mail, 0)
var TestQuitGameRequestWrongSlice = make([]*api.Mail, 0)

func NewEnterGameRequestTopMsg(pid int32, rid int32) *proto.TopMessage {
	return codec.EncodeRequest(&proto.PlayerEnterGameRequest{
		PlayerId:   pid,
		GameRoomId: rid,
	})
}
func NewQuitGameRequestTopMsg(pid int32) *proto.TopMessage {
	return codec.EncodeRequest(&proto.PlayerQuitGameRequest{
		PlayerId: pid,
		HeroId:   -1,
	})
}

//testInit 创建configs.MaxNormalGamePlayerNum的正确的登录请求填到TestEnterGameRequestRightMap里,
//创建3*configs.MaxNormalGamePlayerNum的测试用例填到TestEnterGameRequestWrongMap里
func testInit() {
	offset := int32(10000)
	for i := int32(0); i < configs.MaxNormalGamePlayerNum; i++ {
		pId := offset + i
		//填充PlayerInfo
		p := framework.PlayerInfo{PlayerId: pId}
		TestJoinPlayers = append(TestJoinPlayers, &p)

		//填充EnterGameRequest
		rightEnterMsg := NewEnterGameRequestTopMsg(pId, TestRoomId)        //正确的进入游戏的请求
		rightEnterMsgDup := NewEnterGameRequestTopMsg(pId, TestRoomId)     //正确但是重复的进入游戏的请求
		wrongEnterMsg1 := NewEnterGameRequestTopMsg(pId+10000, TestRoomId) //错误的进入游戏的请求(Pid不在GameServer传入的PlayerInfo中)
		wrongEnterMsg2 := NewEnterGameRequestTopMsg(pId+20000, TestRoomId) //错误的进入游戏的请求(Pid不在GameServer传入的PlayerInfo中)
		wrongEnterMsg3 := NewEnterGameRequestTopMsg(pId, TestRoomId+1)     //错误的进入游戏的请求(RoomId错误)
		rightEnterMail := api.Mail{
			Conn: nil,
			Msg:  rightEnterMsg,
		}
		rightEnterMailDup := api.Mail{
			Conn: nil,
			Msg:  rightEnterMsgDup,
		}
		wrongEnterMail1 := api.Mail{
			Conn: nil,
			Msg:  wrongEnterMsg1,
		}
		wrongEnterMail2 := api.Mail{
			Conn: nil,
			Msg:  wrongEnterMsg2,
		}
		wrongEnterMail3 := api.Mail{
			Conn: nil,
			Msg:  wrongEnterMsg3,
		}
		TestEnterGameRequestRightSlice = append(TestEnterGameRequestRightSlice, &rightEnterMail)
		TestEnterGameRequestWrongSlice = append(TestEnterGameRequestWrongSlice, &wrongEnterMail1)
		TestEnterGameRequestWrongSlice = append(TestEnterGameRequestWrongSlice, &wrongEnterMail2)
		TestEnterGameRequestWrongSlice = append(TestEnterGameRequestWrongSlice, &wrongEnterMail3)
		TestEnterGameRequestWrongSlice = append(TestEnterGameRequestWrongSlice, &rightEnterMailDup)

		//填充QuitGameRequest
		rightQuitMsg := NewQuitGameRequestTopMsg(pId)          //正确的退出游戏请求
		rightQuitMsgDup := NewQuitGameRequestTopMsg(pId)       //正确但是重复的退出游戏请求
		wrongQuitMsg1 := NewQuitGameRequestTopMsg(pId + 10000) //错误的退出游戏请求(Pid不在本局游戏的Pid中)
		wrongQuitMsg2 := NewQuitGameRequestTopMsg(pId + 20000) //错误的退出游戏请求(Pid不在本局游戏的Pid中)
		wrongQuitMsg3 := NewQuitGameRequestTopMsg(pId + 30000) //错误的退出游戏请求(Pid不在本局游戏的Pid中)
		rightQuitMail := api.Mail{
			Conn: nil,
			Msg:  rightQuitMsg,
		}
		rightQuitMailDup := api.Mail{
			Conn: nil,
			Msg:  rightQuitMsgDup,
		}
		WrongQuitMail1 := api.Mail{
			Conn: nil,
			Msg:  wrongQuitMsg1,
		}
		WrongQuitMail2 := api.Mail{
			Conn: nil,
			Msg:  wrongQuitMsg2,
		}
		WrongQuitMail3 := api.Mail{
			Conn: nil,
			Msg:  wrongQuitMsg3,
		}
		TestQuitGameRequestRightSlice = append(TestQuitGameRequestRightSlice, &rightQuitMail)
		TestQuitGameRequestWrongSlice = append(TestQuitGameRequestWrongSlice, &rightQuitMailDup)
		TestQuitGameRequestWrongSlice = append(TestQuitGameRequestWrongSlice, &WrongQuitMail1)
		TestQuitGameRequestWrongSlice = append(TestQuitGameRequestWrongSlice, &WrongQuitMail2)
		TestQuitGameRequestWrongSlice = append(TestQuitGameRequestWrongSlice, &WrongQuitMail3)
	}
}

//TestEnterGameRequest 测试登录的时候有很多个假的请求和很多个真的请求的情况
func TestGameRequest(t *testing.T) {
	logger.SetLogLevel(logger.LogInfoLevel) //打印正常游戏过程的日志
	logger.SetLogLevelToTestOnly()          //打印测试日志
	testInit()
	room := new(gr.NormalGameRoom)
	var TestRoomInitInfo = framework.RoomInitInfo{
		Id:          1,
		Port:        "8000",
		Over:        make(chan interface{}),
		JoinPlayers: TestJoinPlayers,
	}
	room.Init(&TestRoomInitInfo)
	go room.Start()
	//把正确消息投送到消息处理管道中
	go func() {
		for _, mail := range TestEnterGameRequestRightSlice {
			time.Sleep(100 * time.Millisecond)
			room.TestRequestChan <- mail
		}
	}()
	//把错误消息投送到消息处理管道中
	go func() {
		for _, mail := range TestEnterGameRequestWrongSlice {
			time.Sleep(25 * time.Millisecond)
			room.TestRequestChan <- mail
		}
	}()

	time.Sleep(5 * time.Second)
	dup := make(map[int32]struct{})
	res := make(map[int32]struct{})
	for i := int32(1); i <= configs.MaxNormalGamePlayerNum; i++ {
		res[i] = struct{}{}
	}

	//查看分配的id有没有重复的
	for _, v := range TestEnterGameRequestRightSlice {
		pId := v.Msg.Request.PlayerEnterGameRequest.PlayerId
		hId := room.GetPlayerManager().GetPlayer(pId).HeroId
		if _, ok := dup[pId]; ok {
			t.Fatalf("pId:%d registered two times!", pId)
		}
		dup[pId] = struct{}{}
		delete(res, hId)
	}
	//查看有没有玩家被分到不正确的id
	if len(res) != 0 {
		t.Fatalf("not all heroId distributed!,left:%v", res)
	}

	//把正确的退出消息投送到消息处理管道中
	go func() {
		for _, mail := range TestQuitGameRequestRightSlice {
			time.Sleep(1000 * time.Millisecond)
			room.TestRequestChan <- mail
		}
	}()
	//把错误的退出消息投送到消息处理管道中
	go func() {
		for _, mail := range TestQuitGameRequestWrongSlice {
			time.Sleep(200 * time.Millisecond)
			room.TestRequestChan <- mail
		}
	}()
	time.Sleep(12 * time.Second)
}

func TestPlayerQuit(t *testing.T) {

}
