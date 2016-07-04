package mongo

import (
	"errors"
	"strings"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"../../common"
)

const (
	ErrNoHosts          string = "No database hosts given"
	ErrNoHostsReachable string = "no reachable servers"
)

type Database struct {
	session *mgo.Session
	hosts   []string
}

func Connect(hosts []string) (db *Database, err error) {
	if len(hosts) == 0 {
		err = errors.New(ErrNoHosts)
		return
	}

	session, err := mgo.Dial(strings.Join(hosts, ","))
	if err != nil {
		return
	}

	db = &Database{
		session: session,
		hosts:   hosts,
	}
	return
}

func (self Database) Write(msg *common.Packet) error {
	c := self.session.DB("gossip").C("messages")

	// TODO: final sanitation
	var err = c.Insert(msg)
	return err
}

func (self Database) LoadSince(since time.Time) (result []*common.Packet) {
	c := self.session.DB("gossip").C("messages")
	_ = c.Find(bson.M{"timestamp": map[string]interface{}{"$gte": since}}).All(&result)

	// TODO: handle error
	return result
}

func (self Database) LoadLast(count int64) (result []*common.Packet) {
	result = make([]*common.Packet, count)

	c := self.session.DB("gossip").C("messages")
	_ = c.Find(nil).Sort("-time").Limit(int(count)).All(&result)

	var length = len(result) - 1
	for i := 0; i < length/2; i++ {
		result[i], result[length-i] = result[length-i], result[i]
	}

	// TODO: handle error
	return result
}

func (self Database) Close() error {
	self.session.Close()
	return nil
}
