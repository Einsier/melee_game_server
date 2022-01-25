package game_room

import (
	"time"
)

/**
*@Author Sly
*@Date 2022/1/18
*@Version 1.0
*@Description:用于注册房间中的定时事件
 */

type TimeEventCallback func(gm *NormalGameRoom)

//TimeEvent 需要通过框架定时的操作,Id最好通过const注册到全局中,以防重复(暂时没写错误处理,错误会宕机),slice为调用的间隔,callback为处理事件
//room为绑定的room,可以通过room来对房间内网络/hero/道具等内容进行定时的管理
type TimeEvent struct {
	Id       int32             //编号
	slice    time.Duration     //间隔的时间
	callback TimeEventCallback //执行的回调函数
	//room     *NormalGameRoom   //作为callBack的参数
}

func NewTimeEvent(id int32, slice time.Duration, callback TimeEventCallback, room *NormalGameRoom) *TimeEvent {
	return &TimeEvent{
		Id:       id,
		slice:    slice,
		callback: callback,
		//room:     room,
	}
}

type TimeEventController struct {
	addPipe           chan *TimeEvent //存放要加入的timeEvent
	cancelPipe        chan int32      //存放要取消的timeEvent的编号
	timeEventCloseMap map[int32]chan interface{}
	room              *NormalGameRoom
}

//NewTimeEventController 创建(并运行)一个NewTimeEventController
func NewTimeEventController(room *NormalGameRoom) *TimeEventController {
	c := &TimeEventController{
		addPipe:           make(chan *TimeEvent),
		cancelPipe:        make(chan int32),
		timeEventCloseMap: make(map[int32]chan interface{}),
		room:              room,
	}
	go c.manageEventPipe()
	return c
}

//AddEvent 添加定时事件,注意一定不能重复添加...否则宕机
func (c *TimeEventController) AddEvent(event *TimeEvent) {
	c.addPipe <- event
}

//CancelEvent 通过注册过的某个event的id来删除该事件,注意一定不能重复删除...否则宕机
func (c *TimeEventController) CancelEvent(eventId int32) {
	c.cancelPipe <- eventId
}

func (c *TimeEventController) manageEventPipe() {
	for {
		select {
		case event := <-c.addPipe:
			closer := make(chan interface{})
			c.timeEventCloseMap[event.Id] = closer
			go c.handleEvent(event, closer)
		case id := <-c.cancelPipe:
			closer, ok := c.timeEventCloseMap[id]
			if ok {
				close(closer)
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
