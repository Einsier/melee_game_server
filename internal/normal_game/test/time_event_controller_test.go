package test

/**
*@Author Sly
*@Date 2022/1/18
*@Version 1.0
*@Description:
 */

/*func TestTimeEventController(t *testing.T) {
	const(
		CountPeopleTimeEventCode = iota
		PrintTimeTimeEventCode
		AddPeopleTimeEventCode
		DecPeopleTimeEventCode
		IdleEventCode
	)

	type GameRoomForTest struct {
		peopleNum	int32
	}

	CountPeopleTimeEventSlice := 1 * time.Second
	PrintTimeTimeEventSlice := 5 * time.Second
	AddPeopleTimeEventSlice := 100 * time.Millisecond
	DecPeopleTimeEventSlice := 50 * time.Millisecond
	IdleEventSlice	:= 1 * time.Second

	IdleEventCallback := func (gm *GameRoomForTest){
		fmt.Printf("%v\n","idle...")
	}

	AddPeopleTimeEventCallback := func(gm *GameRoomForTest){
		atomic.AddInt32(&gm.peopleNum,1)
	}

	DecPeopleTimeEventCallback := func(gm *GameRoomForTest){
		atomic.AddInt32(&gm.peopleNum,-1)
	}

	CountPeopleTimeEventCallback := func(gm *GameRoomForTest){
		fmt.Printf("current people in room:%v\n",gm.peopleNum)
	}

	PrintTimeTimeEventCallback := func(gm *GameRoomForTest){
		fmt.Printf("current time:%v\n",time.Now().String())
	}

	room := GameRoomForTest{peopleNum: 100}
	countPeopleEvent := gc.NewTimeEvent(CountPeopleTimeEventCode,CountPeopleTimeEventSlice,CountPeopleTimeEventCallback,&room)
	printTimeTimeEvent := gc.NewTimeEvent(PrintTimeTimeEventCode,PrintTimeTimeEventSlice,PrintTimeTimeEventCallback,&room)
	addPeopleEvent := gc.NewTimeEvent(AddPeopleTimeEventCode,AddPeopleTimeEventSlice,AddPeopleTimeEventCallback,&room)
	decPeopleEvent := gc.NewTimeEvent(DecPeopleTimeEventCode,DecPeopleTimeEventSlice,DecPeopleTimeEventCallback,&room)
	idleEvent := gc.NewTimeEvent(IdleEventCode,IdleEventSlice,IdleEventCallback,&room)
	c := gc.NewTimeEventController()
	c.AddEvent(countPeopleEvent)
	c.AddEvent(printTimeTimeEvent)
	c.AddEvent(addPeopleEvent)
	time.Sleep(10 * time.Second)
	c.CancelEvent(AddPeopleTimeEventCode)
	fmt.Printf("canceled:add people event\n")
	time.Sleep(3*time.Second)
	c.AddEvent(decPeopleEvent)
	c.AddEvent(idleEvent)
	fmt.Printf("add:idle event\n")
	fmt.Printf("add:decPeopleEvent event\n")
	time.Sleep(20 * time.Second)
}*/
