package test

import (
	"fmt"
	"melee_game_server/internal/normal_game/game_room"
	"sync/atomic"
	"testing"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/18
*@Version 1.0
*@Description:
 */

func TestTimeEventController(t *testing.T) {
	const (
		CountPeopleTimeEventCode = iota
		PrintTimeTimeEventCode
		AddPeopleTimeEventCode
		DecPeopleTimeEventCode
		IdleEventCode
	)

	CountPeopleTimeEventSlice := 1 * time.Second
	PrintTimeTimeEventSlice := 1 * time.Second
	AddPeopleTimeEventSlice := 100 * time.Millisecond
	DecPeopleTimeEventSlice := 50 * time.Millisecond
	IdleEventSlice := 1 * time.Second

	Idle1EventCallback := func(gm *game_room.NormalGameRoom) {
		fmt.Printf("%v\n", "idle1...")
	}
	Idle2EventCallback := func(gm *game_room.NormalGameRoom) {
		fmt.Printf("%v\n", "idle2...")
	}

	AddPeopleTimeEventCallback := func(gm *game_room.NormalGameRoom) {
		atomic.AddInt32(&gm.PlayerNum, 1)
	}

	DecPeopleTimeEventCallback := func(gm *game_room.NormalGameRoom) {
		atomic.AddInt32(&gm.PlayerNum, -1)
	}

	CountPeopleTimeEventCallback := func(gm *game_room.NormalGameRoom) {
		fmt.Printf("current people in room:%v\n", gm.PlayerNum)
	}

	PrintTimeTimeEventCallback := func(gm *game_room.NormalGameRoom) {
		fmt.Printf("current time:%v\n", time.Now().String())
	}

	room := game_room.NormalGameRoom{PlayerNum: 100}
	countPeopleEvent := game_room.NewTimeEvent(CountPeopleTimeEventCode, CountPeopleTimeEventSlice, CountPeopleTimeEventCallback)
	printTimeTimeEvent := game_room.NewTimeEvent(PrintTimeTimeEventCode, PrintTimeTimeEventSlice, PrintTimeTimeEventCallback)
	addPeopleEvent := game_room.NewTimeEvent(AddPeopleTimeEventCode, AddPeopleTimeEventSlice, AddPeopleTimeEventCallback)
	decPeopleEvent := game_room.NewTimeEvent(DecPeopleTimeEventCode, DecPeopleTimeEventSlice, DecPeopleTimeEventCallback)
	idle1Event := game_room.NewTimeEvent(IdleEventCode, IdleEventSlice, Idle1EventCallback)
	idle2Event := game_room.NewTimeEvent(IdleEventCode, IdleEventSlice, Idle2EventCallback)
	c := game_room.NewTimeEventController(&room)
	c.AddEvent(countPeopleEvent)
	c.AddEvent(printTimeTimeEvent)
	c.AddEvent(addPeopleEvent)
	c.AddEvent(idle1Event)
	//测试删除不存在的编号
	c.CancelEvent(100)
	c.CancelEvent(100)
	c.CancelEvent(100)
	time.Sleep(5 * time.Second)

	//测试重复删除
	fmt.Printf("测试重复\n")
	c.CancelEvent(AddPeopleTimeEventCode)
	c.CancelEvent(AddPeopleTimeEventCode)
	c.CancelEvent(AddPeopleTimeEventCode)

	//测试使用新的id覆盖原来的id,使用idle1替代idle2
	fmt.Printf("使用idle2代替idle1\n")
	c.AddEvent(idle2Event)

	fmt.Printf("撤销:add people event\n")
	time.Sleep(3 * time.Second)
	c.AddEvent(decPeopleEvent)
	fmt.Printf("添加:idle event\n")
	fmt.Printf("添加:decPeopleEvent event\n")
	time.Sleep(5 * time.Second)
	c.Destroy()
	fmt.Printf("已删除所有的event\n")
	time.Sleep(5 * time.Second)
}
