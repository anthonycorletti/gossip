package common

import (
	// "html"
	"time"
)

var SERVER_ALIAS = User{
	ID:       "0",
	Username: "server",
	Status:   "healthy",
}

type Packet struct {
	// message data
	Sender   *User  `json:"sender"`
	Receiver *User  `json:"receiver"`
	Body     string `json:"body"`
	Action   string `json:"action"`

	// metadata
	Time time.Time `json:"timestamp"`
}

func (self *Packet) String() string {
	if self.Sender == nil {
		// TODO: report error
		return ""
	}
	return self.Sender.Username + " says " + self.Body
}

func (self *Packet) Sanitize() {
	// TODO: full sanitization
	// TODO: error on no sender?
	// TODO: what to do about empty bodies, actions

	if self.Time.IsZero() {
		self.Time = time.Now()
	}

	if self.Receiver == nil {
		self.Receiver = &SERVER_ALIAS
	}
}
