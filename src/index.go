package main

import (
	"fmt"
	"os"
	"os/signal"

	"./chat"
	"./logging"
	"./web"
)

const (
	EXIT_NORMAL      int = 0
	EXIT_CONFIG_ERR  int = 100
	EXIT_STARTUP_ERR int = 110
	EXIT_GENERAL_ERR int = 200
)

var logger *logging.Router
var persistence chat.Database

func main() {
	// initialize the CLI options
	// returns any errors to be caugh by the checkError function
	checkError(cliInit(), "There was a command-line error", EXIT_CONFIG_ERR)

	// start logging how that the system is initialized
	// currently have no use for the actual logger object
	var err error
	logger, err = logging.Initialize(Config.LogPath, Config.Quiet)
	checkError(err, "Could not load log file", EXIT_STARTUP_ERR)

	// setup watchers for termination signals
	var _, lock = initSignalWatchers()

	// start the web
	web.Initialize(lock, logger, Config.Web.Port, Opts.Directory)

	// start the chat server
	logger.Println(Config.Database.Hosts)
	persistence = chat.Initialize(lock, logger, Config.Chat.Port, Config.Database.Hosts)

	for {
		err := <-lock

		if err == nil {
			// didn't actually get an error
			continue
		} else if err.Error() == "broken pipe" {
			// if there's a broken pipe.... it sucks
			// we can't crash the server though, so
			// ditch the pipe and keep going
			logger.Println("ERR: A websocket pipe was broken\n")
			continue
		}

		logger.Println("Error received on channel")
		checkError(err, "Received error from server", EXIT_GENERAL_ERR)
	}
}

func checkError(actual error, text string, code int) {
	if actual != nil {
		if logger != nil {
			logger.Errorf("%s:\n%s\n", text, actual)
		} else {
			fmt.Printf("Fatal error killed logger: %#v\n", actual)
		}

		Cleanup(code)
	}
}

func initSignalWatchers() (sigLock chan os.Signal, deadLock chan error) {
	// make the signal channel
	sigLock = make(chan os.Signal, 10)
	signal.Notify(sigLock, os.Interrupt)

	// make the global error channel
	deadLock = make(chan error, 10)

	// this functions spins up in a separate context
	// blocks at the first line until a signal is passed on channel
	go func() {
		// we don't care about value here, just that a signal was caught
		<-sigLock
		Cleanup(EXIT_NORMAL)
	}()

	// this functions spins up in a separate context
	// blocks at first line of every iteration of for loop until an error is passed
	// the for loop is included in the case on non-fatal errors -- keeps the system going

	return sigLock, deadLock
}

func Cleanup(code int) {
	persistence.Close()
	logging.Close()
	os.Exit(code)
}
