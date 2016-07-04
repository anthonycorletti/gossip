package chat

import (
	"encoding/json"
	"io"
	"time"

	"./common"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"
	"github.com/russross/blackfriday"
)

const (
	DEFAULT_CLIENT_BUFFER int   = 10
	DEFAULT_SYNC          uint8 = 1
)

type Client struct {
	ID       string
	Username string

	Socket    *websocket.Conn
	Receive   chan *common.Packet
	Hangup    chan bool
	Heartbeat chan *common.Packet
	SyncLock  uint8
}

type ClientSync struct {
	Members  map[string]*common.User `json:"members"`
	Messages []*common.Packet        `json:"messages"`
}

func (self *Client) Listen(lock chan error, done chan *Client) {
	go func() {
		for {
			hangup := <-self.Hangup
			if hangup {
				done <- self
				return
			}
		}
	}()

	for {
		self.Read(lock)
	}
}

func (self *Client) Read(lock chan error) {
	var msg common.Packet
	var err = websocket.JSON.Receive(self.Socket, &msg)

	if err != nil && err.Error() != ErrClosedNetwork {
		if err == io.EOF {
			self.Hangup <- true
			return
		}

		lock <- err
		return
	}

	msg.Time = time.Now()
	msg.Sender.ID = self.ID // verify and sanitize

	switch msg.Action {
	case "heartbeat":
		logger.Printf("Received Heartbeat Pong From: %s\n", self.ID)
		if self.ShouldSync(DEFAULT_SYNC) { // NOTE: "emergent bug" if you change the SynLock value in NewClient()
			self.Sync(lock, server, false) // TODO: do not use module global
		}
		self.Heartbeat <- &msg
		break

	default:
		logger.Printf("Received Message From: %s\n", self.ID)
		self.Receive <- &msg
		break
	}
}

func (self *Client) Write(lock chan error, msg *common.Packet) {
	var err = websocket.JSON.Send(self.Socket, msg)

	if err != nil && err.Error() != ErrClosedNetwork {
		logger.Errorln(err)
		lock <- err
		return
	}
}

func (self *Client) Ping(lock chan error) {
	self.Write(lock, &common.Packet{
		Action: "heartbeat",
		Sender: &common.SERVER_ALIAS,
		Body:   "ping",
	})
}

func (self *Client) PingPong(lock chan error, serv *Server) {
	go func(server *Server) {
		for _ = range time.Tick(5 * time.Second) {
			self.Ping(lock)

			msg := <-self.Heartbeat

			if msg.Body != "pong" {
				logger.Errorf("Invalid Heartbeat From: %s\n\tReceived: %s\n\tExpected: pong\n", self.ID, msg.Body)
				server.Leave <- self
			}
		}
	}(serv)
}

func (self *Client) Sync(lock chan error, serv *Server, fullSync bool) {
	var messages []*common.Packet = nil
	if fullSync {
		messages = serv.Backing.LoadLast(50)
		for _, msg := range messages {
			var temp = blackfriday.MarkdownCommon([]byte(msg.Body))
			msg.Body = string(temp)
			msg.Sanitize()
		}
	}

	var syncData, err = json.Marshal(ClientSync{
		Members:  serv.UserList(),
		Messages: messages,
	})
	if err != nil {
		lock <- err
		return
	}

	self.Write(lock, &common.Packet{
		Action: "sync",
		Sender: &common.SERVER_ALIAS,
		Body:   string(syncData),
	})
}

func (self *Client) ShouldSync(updated uint8) bool {
	self.SyncLock -= 1
	if self.SyncLock < 1 {
		self.SyncLock = updated
		return true
	}

	return false
}

func NewClient(sock *websocket.Conn, pipe chan *common.Packet) *Client {
	if sock == nil {
		logger.Errorln("Received nil socket")
		return nil
	}

	return &Client{
		ID:     uuid.New(),
		Socket: sock,

		Receive:   pipe,
		Hangup:    make(chan bool, 0),           // blocks until read
		Heartbeat: make(chan *common.Packet, 0), // TODO: test conditions of blocking and missing heartbeat
		SyncLock:  DEFAULT_SYNC,                 // sync = 1 yields updated users every heartbeat
	}
}
