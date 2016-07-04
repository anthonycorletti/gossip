package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

const (
	DEFAULT_WEB_PORT  int    = 8765
	DEFAULT_CHAT_PORT int    = 5432
	DEFAULT_LOG_PATH  string = "./logs"
)

var Opts Options
var Config Configuration
var Verbosity int

type Options struct {
	Verbose   []bool `short:"v"   long:"verbose"        description:"Show verbose log information. Supports -v[vvv] syntax."`
	Quiet     bool   `short:"q"   long:"quiet"          description:"Do not output any text to STDOUT."`
	LogPath   string `short:"l"   long:"log"            description:"Path to folder where logs should be saved" default:"./logs"`
	WebPort   int    `short:"w"   long:"web-port"       description:"Port to listen for HTTP traffic" default:"8765"`
	ChatPort  int    `short:"c"   long:"chat-port"      description:"Port to listen for Chat traffic" default:"7654"`
	Directory string `short:"d"   long:"directory"      description:"Directory containing gossip data/assets" default:"/etc/gossip"`
}

type Configuration struct {
	Verbosity int    `json:"verbosity"`
	Quiet     bool   `json:"quiet"`
	LogPath   string `json:"logPath"`

	Web struct {
		Port int `json:"port"`
	}

	Chat struct {
		Port int `json:"port"`
	}

	Database struct {
		Hosts []string `json:"hosts"`
	}
}

func cliInit() (err error) {
	if _, help := flags.Parse(&Opts); help != nil {
		os.Exit(1)
	}

	// load the base configuration file if it exists
	err = loadConfigFile(Opts.Directory)
	if err != nil {
		if err == os.ErrNotExist {
			logger.Warn("No config file in given directory")
		} else {
			return
		}
	}

	// begin overwriting based on command line

	if len(Opts.Verbose) > 0 {
		Config.Verbosity = len(Opts.Verbose)
	}

	if Opts.Quiet {
		Config.Quiet = Opts.Quiet
	}

	if Opts.LogPath != DEFAULT_LOG_PATH {
		Config.LogPath = Opts.LogPath
	}

	if Opts.WebPort != DEFAULT_WEB_PORT {
		Config.Web.Port = Opts.WebPort
	}

	if Opts.WebPort != DEFAULT_CHAT_PORT {
		Config.Chat.Port = Opts.ChatPort
	}

	return
}

func loadConfigFile(path string) error {
	bytes, err := ioutil.ReadFile(filepath.Join(path, "config.json"))
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &Config)
	return err
}
