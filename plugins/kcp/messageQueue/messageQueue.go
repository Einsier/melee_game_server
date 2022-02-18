package messageQueue

import (
	"errors"
	"log"
)

type message interface{}

type MsgQueue struct {
	msgQ     chan message
	capacity uint32
}

func (mq *MsgQueue) Init(capacity uint32) {
	mq.msgQ = make(chan message, capacity)
	mq.capacity = capacity
	log.Println("MsgQueue Init()初始化完成.")
}

func (mq *MsgQueue) Put(msg message) error {
	if mq.msgQ == nil {
		return errors.New("尝试向未初始化的MsgQueue进行放置消息")
	}
	mq.msgQ <- msg
	return nil
}

func (mq *MsgQueue) Get() (message, error) {
	if mq.msgQ == nil {
		log.Println("尝试向未初始化的MsgQueue提取消息")
		return nil, errors.New("尝试向未初始化的MsgQueue提取消息")
	}
	msg := <-mq.msgQ
	return msg, nil
}

func (mq *MsgQueue) Free() {
	close(mq.msgQ)
}
