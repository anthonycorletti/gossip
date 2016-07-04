package file

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"../../common"
)

type Database struct {
	path   string
	file   *os.File
	input  *bufio.Reader
	output *bufio.Writer
}

func Open(path string) (result *Database, err error) {
	var desc *os.File
	desc, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	result = &Database{
		path:   path,
		file:   desc,
		input:  bufio.NewReader(desc),
		output: bufio.NewWriter(desc),
	}

	return
}

func (self Database) Write(msg *common.Packet) (err error) {
	var out = fmt.Sprintf("%s | %s | %s\n", msg.Sender, msg.Time, msg.Body)
	_, err = self.output.WriteString(out)
	return
}

func (self Database) LoadLast(count int64) []*common.Packet {
	return []*common.Packet{
		{
			Sender: &common.SERVER_ALIAS,
			Action: "message",
			Body:   "We're sorry, but the file backing does not currently support back-tracking",
		},
		{
			Sender: &common.SERVER_ALIAS,
			Action: "message",
			Body:   "Want to help? Submit a PR @ <a href=\"https://github.com/anthcor/gossip\">the project page</a>",
		},
	}
}

func (self Database) LoadSince(since time.Time) []*common.Packet {
	return []*common.Packet{
		{
			Sender: &common.SERVER_ALIAS,
			Action: "message",
			Body:   "We're sorry, but the file backing does not currently support back-tracking",
		},
		{
			Sender: &common.SERVER_ALIAS,
			Action: "message",
			Body:   "Want to help? Submit a PR @ <a href=\"https://github.com/anthcor/gossip\">the project page</a>",
		},
	}
}

func (self Database) Close() (err error) {
	err = self.output.Flush()
	return
}

func (self Database) GetPath() string {
	return self.path
}

func (self Database) GetFile() *os.File {
	return self.file
}
