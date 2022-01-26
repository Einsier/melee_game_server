package game_room

import (
	"melee_game_server/api/proto"
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

func (c *RequestController) Work(room *NormalGameRoom) {
	for {
		mail := <-room.TestRequestChan
		if mail.Msg == nil || mail.Msg.Request == nil {
			logger.TestErrf("收到了错误的请求:%v", mail)
			continue
		}
		h := c.register.GetHandler(int32(mail.Msg.Request.RequestCode))
		if h != nil {
			go h(mail, room)
		}
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
