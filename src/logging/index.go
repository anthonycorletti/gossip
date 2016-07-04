package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	// show the date, time, microseconds, and full filename path + line number
	LOG_PROPERTIES     int = log.Ldate | log.Ltime
	LOG_PROPERTIES_WRN int = log.Ldate | log.Ltime | log.Lshortfile
	LOG_PROPERTIES_ERR int = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
)

var openFiles []*os.File
var handler *Router

func Initialize(path string, quiet bool) (logger *Router, err error) {
	// TODO: accept non-fs path strings (i.e. allow network paths)
	// TODO: once we accpet those path, feed them to modular logging backends

	openFiles = make([]*os.File, 0)

	// load log file only if path is actually given
	// this should default to "./logs" but can be overriden to silence
	var file, errFile *os.File
	if path != "" {
		// create standard log
		file, err = OpenFile(filepath.Join(path, "gossip.log"))
		if err != nil {
			return
		}

		// create error log
		errFile, err = OpenFile(filepath.Join(path, "error.log"))
		if err != nil {
			return
		}
	}

	// if the quiet flag is false, write to both the STDOUT and file
	// otherwise, only log to file
	if quiet == false && file != nil && errFile != nil {
		var multi = io.MultiWriter(file, os.Stdout)
		var multiErr = io.MultiWriter(file, errFile, os.Stderr)
		handler = &Router{
			out:  log.New(multi, "     ", LOG_PROPERTIES),
			warn: log.New(multi, "WRN: ", LOG_PROPERTIES_WRN),
			err:  log.New(multiErr, "ERR: ", LOG_PROPERTIES_ERR),
			web:  log.New(multi, "WEB: ", LOG_PROPERTIES),
		}
	} else {
		handler = &Router{
			out:  log.New(file, "     ", LOG_PROPERTIES),
			warn: log.New(file, "WRN: ", LOG_PROPERTIES_WRN),
			err:  log.New(file, "ERR: ", LOG_PROPERTIES_ERR),
			web:  log.New(file, "WEB: ", LOG_PROPERTIES),
		}
	}

	logger = handler
	logger.Directory = path
	return
}

func OpenFile(logPath string) (file *os.File, err error) {
	// make sure the owning path exists -- if not, attempt to create it
	err = os.MkdirAll(filepath.Dir(logPath), 0774)
	if err != nil {
		return
	}

	// now, open the file in RW mode (creating if necessary)
	file, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		return
	}

	// add the file to the openFiles list
	openFiles = append(openFiles, file)
	return
}

func Close() {
	for _, f := range openFiles {
		f.Close()
	}
}
