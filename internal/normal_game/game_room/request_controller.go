package game_room

import (
	"melee_game_server/api/client/proto"
	"melee_game_server/framework/game_net/api"
	"melee_game_server/plugins/logger"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/19
*@Version 1.0
*@Description:
 */

type RequestController struct {
	register *RequestHandlerRegister
}

func NewRequestController() (c *RequestController) {
	r := NewRequestHandlerRegister()
	r.NormalGameInit()
	return &RequestController{register: r}
}

//Work1 需要使用go Work1调用,接收消息并根据消息类型开go程执行回调函数
func (c *RequestController) Work1(room *NormalGameRoom) {
	for {
		//todo 测试用
		//mail := <-room.TestRequestChan
		mail, ok := room.netServer.Receive()
		if !ok {
			//防止go程泄露
			return
		}
		if mail.Msg == nil || mail.Msg.Request == nil {
			//判断消息是否合法
			logger.TestErrf("收到了错误的请求:%v", mail)
			continue
		}
		h := c.register.GetHandler(int32(mail.Msg.Request.RequestCode)) //从回调函数注册中心根据消息类型取出注册函数
		if h != nil {
			go h(mail, room)
		}
	}
}

//Work2 需要使用go Work2调用,接收消息并查看消息合法性,如果合法,投放到管道中让消费者go程执行回调函数
func (c *RequestController) Work2(room *NormalGameRoom) {
	mailQueue := make(chan *api.Mail, 10000)

	for i := 0; i < 1; i++ {
		//开10个消费者go程
		go func() {
			for {
				mail, ok := <-mailQueue
				if !ok {
					//防止go程泄露
					return
				}

				h := c.register.GetHandler(int32(mail.Msg.Request.RequestCode)) //从回调函数注册中心根据消息类型取出注册函数
				if h != nil {
					h(mail, room)
				}
			}
		}()
	}

	//父go程检测消息合法性,如果合法,投入到mailQueue中让子go程消费
	for {
		mail, ok := room.netServer.Receive()
		if !ok {
			//防止go程泄露
			close(mailQueue)
			return
		}
		if mail.Msg == nil || mail.Msg.Request == nil {
			//判断消息是否合法
			logger.TestErrf("收到了错误的请求:%v", mail)
			continue
		}
		mailQueue <- mail
	}
}

type RequestHandler func(s *api.Mail, room *NormalGameRoom)

type RequestHandlerRegister struct {
	handlers map[int32]RequestHandler
	lock     sync.RWMutex
}

func NewRequestHandlerRegister() *RequestHandlerRegister {
	return &RequestHandlerRegister{
		handlers: make(map[int32]RequestHandler),
		lock:     sync.RWMutex{},
	}
}

func (r *RequestHandlerRegister) Register(reqId int32, handler RequestHandler) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.handlers[reqId] = handler
}

func (r *RequestHandlerRegister) Delete(reqId int32) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.handlers, reqId)
}

func (r *RequestHandlerRegister) GetHandler(reqId int32) RequestHandler {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.handlers[reqId]
}

func (r *RequestHandlerRegister) NormalGameInit() {
	r.Register(int32(proto.RequestCode_PlayerEnterGameRequestCode), PlayerEnterGameRequestCallback)
	r.Register(int32(proto.RequestCode_PlayerQuitGameRequestCode), PlayerQuitGameRequestCallback)
	r.Register(int32(proto.RequestCode_HeroPositionReportRequestCode), HeroPositionReportRequestCallback)
	r.Register(int32(proto.RequestCode_HeroMovementChangeRequestCode), HeroMovementChangeRequestCallback)
	r.Register(int32(proto.RequestCode_HeroBulletLaunchRequestCode), HeroBulletLaunchRequestCallback)
	r.Register(int32(proto.RequestCode_HeroSwordAttackRequestCode), HeroSwordAttackRequestCallback)
	r.Register(int32(proto.RequestCode_HeroGetPropRequestCode), HeroGetPropRequestCallback)
	r.Register(int32(proto.RequestCode_HeroBulletColliderHeroRequestCode), HeroBulletColliderHeroRequestCallback)
	r.Register(int32(proto.RequestCode_PlayerHeartBeatRequestCode), PlayerHeartBeatRequestCallback)
}
