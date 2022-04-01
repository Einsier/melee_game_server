/*
 *mailbox包定义了出版的网络收发组件
 */
package mailbox

import (
	"errors"
	"log"
	pb "melee_game_server/api/client/proto"
	mail "melee_game_server/framework/game_net/api"
	adapter "melee_game_server/plugins/kcp/adapter"
	mq "melee_game_server/plugins/kcp/messageQueue"
	"net"
	"sync"
)

/*
 *名称：Mailbox
 *功能：通过KCP协议对网络消息进行收发
 		使用这个结构体的程序,可以通过Init()函数完成Mailbox初始化
		再通过Start()启动mailbox进行工作
		然后通过Receive()和Send()进行消息的接受和发送
 *注意：Receive()和Send()都是同步阻塞的,在未完成消息的投递前不会返回
*/
type Mailbox struct {
	//用于收发消息的网络套接字地址
	addr string
	//接收消息队列和发送消息队列
	receiveMQ mq.MsgQueue
	sendMQ    mq.MsgQueue
	//记录客户端Conn和与之对应chan的map。
	//将要从指定Conn发送的消息将会首先被放入对应的chan中
	mutex    sync.Mutex
	channels map[net.Conn]chan *pb.TopMessage
}

/*
 *函数名:Init
 *功能:完成初始化mailbox套接字地址,初始化接收消息队列和发送消息队列的大小
 */
func (box *Mailbox) Init(addr string, recvSize uint32, sendSize uint32) {
	box.addr = addr
	box.receiveMQ.Init(recvSize)
	box.sendMQ.Init(sendSize)
	box.channels = make(map[net.Conn]chan *pb.TopMessage)
}

/*
 *函数名:Start()
 *功能:
 *	启动mailbox,为程序提供收发消息的服务
 *工作流程：
 *	Start()函数的main goroutine负责监听端口,接收客户端连接
 *	每当一个客户端连接建立成功,Start()将会为这个连接做如下配置
 *		1.建立注册一个容量为1的通道,负责存放将要从这个连接发送出去的消息
 *		2.两个处理goroutine：sendHandler和receive负责这个连接上的消息的收发
 *  另有一个sendMQHandler()goroutine负责将sendMQ中的消息分发给各个连接进行发送
 */
func (box *Mailbox) Start() error {
	//基于KCP协议监听套接字端口
	//listener, err := kcp.Listen(box.addr)
	listener, err := net.Listen("tcp", box.addr)
	if err != nil {
		log.Println(err)
		return errors.New("启动mailbox失败:" + err.Error())
	}
	//启动sendMQHandler()goroutine
	//负责将box.sendMQ中消息转发给各个连接进行发送
	go box.sendMQHandler()
	//负责监听端口，接受连接
	//为每个连接初始化配置
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("接收到新连接")
			//为新建立的连接分配一个容量为1的通道
			//存放将要从这个连接发送出去的消息
			channel := make(chan *pb.TopMessage, 1)
			//在map中注册连接和通道的对应关系
			box.mutex.Lock()
			box.channels[conn] = channel
			box.mutex.Unlock()
			//为新建立的连接启动收发消息的goroutine
			go box.receiveHandler(conn)
			go box.sendHandler(conn)
		}
	}()
	return nil
}

/*
 *函数名:
 *功能:
 */
func (box *Mailbox) Shutdown() {

}

/*
 *函数名:
 *	Receive()
 *功能:
 *	使用mailbox的程序通过Receive()函数取出一个消息
 *注意:
 *	Receive()是同步阻塞的,当receiveMQ中没有消息时，会阻塞调用线程
 */
func (box *Mailbox) Receive() *mail.Mail {
	msg, err := box.receiveMQ.Get()
	if err != nil {
		return nil
	}
	mailPtr := msg.(*mail.Mail)
	return mailPtr
}

/*
 *函数名:
 *	Send()
 *功能:
 *	使用mailbox的程序可以通过Send()发送一个消息
 *注意:
 *	Send()是同步阻塞的,当sendMQ中消息满时，会阻塞调用线程
 *	Send()并不负责真正的发送工作,它仅仅是将消息投递到发送消息队列中
 */
func (box *Mailbox) Send(replyPtr *mail.ReplyMail) {
	//将发送消息存放到sendMQ中
	box.sendMQ.Put(replyPtr)
}

/*
 *函数名:
 *	sendMQHandler
 *功能:
 *	负责将sendMQ中的消息派发给"将要发送消息"的连接的通道中
 */
func (box *Mailbox) sendMQHandler() {
	for {
		//当sendMQ为空时，将会阻塞
		//当sendMQ中有消息时,取出消息进行投递
		msg, err := box.sendMQ.Get()
		if err != nil {
			log.Println(err)
			continue
		}
		replyPtr := msg.(*mail.ReplyMail)
		//投递给所有接受消息的连接
		for _, connPtr := range replyPtr.ConnSlice {
			box.mutex.Lock()
			channel := box.channels[connPtr]
			box.mutex.Unlock()
			channel <- replyPtr.Msg
		}
	}
}

/*
 *函数名:
 *	receiveHandler
 *功能:
 *	接受连接中的消息,封装成为Mail对象，存放到box.receiveMQ中,如果接收到错误的消息,返回nil并断开连接
 */
func (box *Mailbox) receiveHandler(conn net.Conn) *pb.TopMessage {
	for {
		msg := adapter.Receive(conn)
		m := mail.Mail{Conn: conn, Msg: msg}
		box.receiveMQ.Put(&m)
		if msg == nil {
			//如果收到下层返回的空包,那么本层放到mq中,通知上层之后,自己退出
			conn.Close()
			return nil
		}
	}
}

/*
 *函数名:
 *	sendHandler
 *功能:
 *	将派送到通道中的消息真正地从连接中发送出去
 */
func (box *Mailbox) sendHandler(conn net.Conn) {
	box.mutex.Lock()
	channel := box.channels[conn]
	box.mutex.Unlock()
	for {
		reply := <-channel
		adapter.Send(conn, reply)
	}
}

/*
 *函数名:
 *	broadcast
 *功能:
 *	用于开发阶段测试服务器和客户端之间的发送功能
 */
func (box *Mailbox) Broadcast(msg *pb.TopMessage) {
	s := make([]net.Conn, 0, len(box.channels))
	for conn := range box.channels {
		s = append(s, conn)
	}
	var reply mail.ReplyMail
	reply.ConnSlice = s
	reply.Msg = msg
	box.Send(&reply)
}
