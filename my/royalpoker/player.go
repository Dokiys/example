package main

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"sync"
)

type Player interface {
	GetId() int
	GetName() string
	Send(ctx context.Context, data []byte)
	Receive(ctx context.Context) ([]byte, error)
	Close(ctx context.Context)
	SetConn(ctx context.Context, conn *websocket.Conn)
}
type PlayerWs struct {
	Id   int
	Name string
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send    chan []byte
	receive chan []byte
	close   chan struct{}

	isWaiting bool
	l         sync.Mutex
}

func NewPlayerWs(conn *websocket.Conn, id int, name string) *PlayerWs {
	player := &PlayerWs{
		Id:      id,
		Name:    name,
		conn:    conn,
		send:    make(chan []byte),
		receive: make(chan []byte),
		close:   make(chan struct{}),
	}
	player.conn.SetCloseHandler(func(code int, text string) error {
		player.close <- struct{}{}
		return nil
	})
	go player.startSend()
	go player.startReceive()
	go player.startClose()

	return player
}

func (self *PlayerWs) GetId() int {
	return self.Id
}

func (self *PlayerWs) GetName() string {
	return self.Name
}

func (self *PlayerWs) Send(ctx context.Context, data []byte) {
	self.send <- data
}

func (self *PlayerWs) Receive(ctx context.Context) ([]byte, error) {
	self.l.Lock()
	self.isWaiting = true
	defer func() {
		self.isWaiting = false
		self.l.Unlock()

	}()
	select {
	case data := <-self.receive:
		return data, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-self.close:
		return nil, errors.New("玩家退出！")
	}
}

func (self *PlayerWs) Close(ctx context.Context) {
	self.close <- struct{}{}
	//self.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "关闭连接"))
	err := self.conn.Close()
	if err != nil {
		logrus.Errorf("关闭链接失败：%s", err.Error())
	}
}

func (self *PlayerWs) SetConn(ctx context.Context, conn *websocket.Conn) {
	//self.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "关闭连接"))
	self.conn.Close()
	self.conn = conn
	self.startWatch()
}

// ============================================
func (self *PlayerWs) startWatch() {
	go self.startSend()
	go self.startReceive()
	go self.startClose()
}

var sendLock sync.Mutex
func (self *PlayerWs) startSend() {
	for {
		select {
		case msg := <-self.send:
			sendLock.Lock()
			err := self.conn.WriteMessage(websocket.TextMessage, msg)
			sendLock.Unlock()
			if err != nil {
				sendLock.Unlock()
				logrus.Errorf("推送玩家消息[%s]错误：%s", msg, err.Error())
				return
			}
		case <-self.close:
			return
		}
	}
}

func (self *PlayerWs) startReceive() {
	for {
		select {
		case <-self.close:
			return
		default:
			_, msg, err := self.conn.ReadMessage()
			if err != nil {
				logrus.Errorf("主动接收玩家消息错误：%s", err.Error())
				return
			}
			if !self.isWaiting {
				continue
			}
			self.receive <- msg
		}
	}
}

func (self *PlayerWs) startClose() {
	select {
	case <-self.close:
		close(self.close)
		close(self.receive)
		self.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "关闭连接"))
		self.conn.Close()
	}
}