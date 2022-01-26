package game_room

import (
	"sync"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/18
*@Version 1.0
*@Description:用于注册房间中的定时事件
 */

type TimeEventCallback func(gm *NormalGameRoom)

//TimeEvent 需要通过框架定时的操作,Id最好通过const注册到全局中,以防重复,如果先注册a再注册重复的b,b会将a覆盖.slice为调用的间隔,callback为处理事件
type TimeEvent struct {
	Id       int32             //编号
	slice    time.Duration     //间隔的时间
	callback TimeEventCallback //执行的回调函数
}

func NewTimeEvent(id int32, slice time.Duration, callback TimeEventCallback) *TimeEvent {
	return &TimeEvent{
		Id:       id,
		slice:    slice,
		callback: callback,
	}
}

type TimeEventController struct {
	addPipe           chan *TimeEvent //存放要加入的timeEvent
	cancelPipe        chan int32      //存放要取消的timeEvent的编号
	timeEventCloseMap sync.Map        //key为event的int32类型的id,value为chan interface{}
	room              *NormalGameRoom
}

//NewTimeEventController 创建(并运行)一个NewTimeEventController
func NewTimeEventController(room *NormalGameRoom) *TimeEventController {
	c := &TimeEventController{
		addPipe:           make(chan *TimeEvent),
		cancelPipe:        make(chan int32),
		timeEventCloseMap: sync.Map{},
		room:              room,
	}
	go c.manageEventPipe()
	return c
}

//AddEvent 添加定时事件,如果重复添加id相同的事件,后来的事件会覆盖先前的事件
func (c *TimeEventController) AddEvent(event *TimeEvent) {
	//如果注册的eventId当前已经有别的事件使用了,则应该删除原来的事件
	if _, ok := c.timeEventCloseMap.Load(event.Id); ok {
		c.CancelEvent(event.Id)
	}
	c.addPipe <- event
}

//CancelEvent 通过注册过的某个event的id来删除该事件,可重复删除
func (c *TimeEventController) CancelEvent(eventId int32) {
	if _, ok := c.timeEventCloseMap.Load(eventId); ok {
		c.cancelPipe <- eventId
	}
}

func (c *TimeEventController) manageEventPipe() {
	for {
		select {
		case event := <-c.addPipe:
			closer := make(chan interface{})
			c.timeEventCloseMap.Store(event.Id, closer)
			go c.handleEvent(event, closer)
		case id := <-c.cancelPipe:
			closer, ok := c.timeEventCloseMap.Load(id)
			if ok {
				close(closer.(chan interface{}))
				c.timeEventCloseMap.Delete(id)
			}
		}
	}
}

func (c *TimeEventController) handleEvent(event *TimeEvent, closer chan interface{}) {
	ticker := time.Tick(event.slice)
	for {
		select {
		case <-ticker:
			event.callback(c.room)
		case <-closer:
			return
		}
	}
}

//Destroy 取消所有的现有的事件
func (c *TimeEventController) Destroy() {
	c.timeEventCloseMap.Range(func(id, _ interface{}) bool {
		c.CancelEvent(id.(int32))
		return true
	})
}
