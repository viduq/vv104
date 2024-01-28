package vv104

import (
	"fmt"
	"log"
)

var (
	LogChan chan string

	logDebug *log.Logger
	logError *log.Logger
	logInfo  *log.Logger

	vv104loggerDebug vv104logger
	vv104loggerError vv104logger
	vv104loggerInfo  vv104logger
)

type vv104logger struct {
	logToChan   bool
	logToStdOut bool
}

func (vvlog vv104logger) Write(p []byte) (n int, err error) {

	if vvlog.logToChan {
		// addLogEntry(string(p))
		LogChan <- string(p)
	}

	if vvlog.logToStdOut {
		n, err = fmt.Print(string(p))
	}

	return n, err
}

func NewLogChan() chan string {
	LogChan = make(chan string, 50)
	return LogChan
}

func initLogger(config Config) {

	vv104loggerDebug = vv104logger{
		logToChan:   config.LogToBuffer,
		logToStdOut: config.LogToStdOut,
	}
	vv104loggerError = vv104logger{
		logToChan:   config.LogToBuffer,
		logToStdOut: config.LogToStdOut,
	}
	vv104loggerInfo = vv104logger{
		logToChan:   config.LogToBuffer,
		logToStdOut: config.LogToStdOut,
	}

	logDebug = log.New(vv104loggerDebug, "Debug: ", log.Ltime|log.Lshortfile)
	logError = log.New(vv104loggerError, "ERROR: ", log.Ltime|log.Lshortfile)
	logInfo = log.New(vv104loggerInfo, "", log.Ltime)

}
