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
	PrintTimeTimeEventSlice := 5 * time.Second
	AddPeopleTimeEventSlice := 100 * time.Millisecond
	DecPeopleTimeEventSlice := 50 * time.Millisecond
	IdleEventSlice := 1 * time.Second

	IdleEventCallback := func(gm *game_room.NormalGameRoom) {
		fmt.Printf("%v\n", "idle...")
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
	countPeopleEvent := game_room.NewTimeEvent(CountPeopleTimeEventCode, CountPeopleTimeEventSlice, CountPeopleTimeEventCallback, &room)
	printTimeTimeEvent := game_room.NewTimeEvent(PrintTimeTimeEventCode, PrintTimeTimeEventSlice, PrintTimeTimeEventCallback, &room)
	addPeopleEvent := game_room.NewTimeEvent(AddPeopleTimeEventCode, AddPeopleTimeEventSlice, AddPeopleTimeEventCallback, &room)
	decPeopleEvent := game_room.NewTimeEvent(DecPeopleTimeEventCode, DecPeopleTimeEventSlice, DecPeopleTimeEventCallback, &room)
	idleEvent := game_room.NewTimeEvent(IdleEventCode, IdleEventSlice, IdleEventCallback, &room)
	c := game_room.NewTimeEventController()
	c.AddEvent(countPeopleEvent)
	c.AddEvent(printTimeTimeEvent)
	c.AddEvent(addPeopleEvent)
	time.Sleep(10 * time.Second)
	c.CancelEvent(AddPeopleTimeEventCode)
	fmt.Printf("canceled:add people event\n")
	time.Sleep(3 * time.Second)
	c.AddEvent(decPeopleEvent)
	c.AddEvent(idleEvent)
	fmt.Printf("add:idle event\n")
	fmt.Printf("add:decPeopleEvent event\n")
	time.Sleep(20 * time.Second)
}
