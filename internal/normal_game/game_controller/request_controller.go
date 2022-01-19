package game_controller

import (
	"melee_game_server/api/proto"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/19
*@Version 1.0
*@Description:
 */

type RequestHandler interface {
	Handle(s *proto.Request)
}

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