package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/russross/blackfriday"

	"../logging"
	"./common"
	"./databases/file"
	"./databases/mongo"
)

const (
	DEFAULT_ENTRANCE_BUFFER  int = 10
	DEFAULT_MESSAGE_BUFFER   int = 50
	DEFAULT_HEARTBEAT_BUFFER int = 100

	DEFAULT_STATUS string = "available"
)

var server *Server
var logger *logging.Router

var ErrClosedNetwork = "use of closed network connection"

type Server struct {
	Clients map[string]*Client
	Users   map[string]*common.User
	Join    chan *Client
	Leave   chan *Client
	Receive chan *common.Packet
	Error   chan error
	Backing Database
}

func Initialize(lock chan error, log *logging.Router, port int, dbHosts []string) Database {
	logger = log

	// first, attempt to connect to mongo
	var persistence = getDataBacking(lock, dbHosts)

	server = &Server{
		Clients: make(map[string]*Client),
		Users:   make(map[string]*common.User),
		Join:    make(chan *Client, DEFAULT_ENTRANCE_BUFFER),
		Leave:   make(chan *Client, DEFAULT_ENTRANCE_BUFFER),
		Receive: make(chan *common.Packet, DEFAULT_MESSAGE_BUFFER),
		Error:   lock,
		Backing: persistence,
	}

	http.Handle("/chat", websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				server.Error <- err
			}
		}()

		logger.Println("Got chat connection")

		var client = NewClient(ws, server.Receive)
		server.Join <- client
		client.Listen(lock, server.Leave)
	}))

	logger.Banner("Starting Chat Server")
	go server.Listen(lock, port)
	return persistence
}

func (self *Server) Listen(lock chan error, port int) {
	go func() {
		for {
			select {
			case joining := <-self.Join:
				self.Clients[joining.ID] = joining
				logger.Printf("New client: %s\n", joining.ID)
				joining.PingPong(lock, self)
				joining.Sync(lock, self, false)
				break

			case leaving := <-self.Leave:
				logger.Printf("Signing off: %s\n", leaving.ID)
				self.AnnounceLeaveToRoom(leaving)
				delete(self.Clients, leaving.ID)
				delete(self.Users, leaving.Username)
				break

			case data := <-self.Receive:
				if data.Action == "join" {
					self.HandleJoin(data)
					continue
				}
				self.Backing.Write(data)
				self.Broadcast(data)
				break

			case err := <-self.Error:
				lock <- err
				break
			}
		}
	}()

	lock <- http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (self *Server) Broadcast(message *common.Packet) {
	// NOTE: here is where secondary santitation should happen
	//	   The message.Sanitize() should be called, but if
	//	   more sanitation is needed do it here

	var output = blackfriday.MarkdownCommon([]byte(message.Body))

	message.Body = string(output)
	message.Sanitize()

	for _, client := range self.Clients {
		client.Write(self.Error, message)
	}
}

func (self *Server) AnnounceJoinToRoom(data *common.Packet, client *Client) {
	self.Broadcast(&common.Packet{
		Action: "joined",
		Sender: &common.SERVER_ALIAS,
		Body:   fmt.Sprintf("<i>%s</i> has joined the room.", data.Sender.Username),
		Time:   time.Now(),
	})
}

func (self *Server) AnnounceLeaveToRoom(leaving *Client) {
	self.Broadcast(&common.Packet{
		Action: "left",
		Sender: &common.SERVER_ALIAS,
		Body:   fmt.Sprintf("<i>%s</i> has left the room.", leaving.Username),
		Time:   time.Now(),
	})
}

func (self *Server) HandleJoin(data *common.Packet) {
	var client = self.Clients[data.Sender.ID]
	var created = self.AddUser(data.Sender.Username, client)
	if created {
		self.AnnounceJoinToRoom(data, client)
		client.Sync(self.Error, self, true)
	}
}

func (self *Server) UserList() map[string]*common.User {
	return self.Users
}

func (self *Server) AddUser(name string, client *Client) (created bool) {
	if _, exists := self.Users[name]; exists {
		client.Write(self.Error, &common.Packet{
			Sender: &common.SERVER_ALIAS,
			Action: "ACK",
			Body:   "exists",
			Time:   time.Now(),
		})
		return false
	}

	client.Username = strings.TrimSpace(name)
	self.Users[name] = &common.User{
		ID:       client.ID,
		Username: client.Username,
		Status:   DEFAULT_STATUS,
	}

	var userData, err = json.Marshal(self.Users[name])
	if err != nil {
		self.Error <- err
		return false
	}

	client.Write(self.Error, &common.Packet{
		Sender: &common.SERVER_ALIAS,
		Action: "ACK",
		Body:   string(userData),
		Time:   time.Now(),
	})
	return true
}

func getDataBacking(lock chan error, hosts []string) Database {
	if mongoBack, err := mongo.Connect(hosts); err != nil {
		lock <- err
	} else {
		return mongoBack
	}

	// last resort -- file backing
	fileBacking, err := file.Open(filepath.Join(logger.Directory, "messages.store"))
	if err != nil {
		lock <- err
		return nil
	}
	return *fileBacking
}
