package vv104

import (
	"fmt"
	"log"
)

var (
	logStrings []string
	logDebug   *log.Logger
	logError   *log.Logger
	logInfo    *log.Logger

	vv104loggerDebug vv104logger
	vv104loggerError vv104logger
	vv104loggerInfo  vv104logger
	// e.g. for use with GUI: GUI can set its update function for the text box here
	LogCallBack func()
)

type vv104logger struct {
	logToBuffer bool
	logToStdOut bool
	logString   string
}

func addLogEntry(s string) {
	logStrings = append(logStrings, s)
}

// API to read from logs (concurrently?)
func ReadLogEntry() string {
	var s string = ""
	for _, logString := range logStrings {
		s += logString
		logStrings = logStrings[1:]

	}
	return s
}

func (vvlog vv104logger) Write(p []byte) (n int, err error) {

	if vvlog.logToBuffer {
		addLogEntry(string(p))
	}

	if vvlog.logToStdOut {
		n, err = fmt.Print(string(p))
	}

	if LogCallBack != nil {
		LogCallBack()
	}

	return n, err
}

func initLogger(config Config) {

	vv104loggerDebug = vv104logger{
		logToBuffer: config.LogToBuffer,
		logToStdOut: config.LogToStdOut,
	}
	vv104loggerError = vv104logger{
		logToBuffer: config.LogToBuffer,
		logToStdOut: config.LogToStdOut,
	}
	vv104loggerInfo = vv104logger{
		logToBuffer: config.LogToBuffer,
		logToStdOut: config.LogToStdOut,
	}

	// var output io.Writer
	// // var multiOutput io.MultiWriter() todo?

	// if config.LogToBuffer {
	// 	output = &LogBuf
	// } else {
	// 	output = os.Stdout
	// }

	// logDebug = log.New(output, "🪲 ", log.Ltime|log.Lshortfile)
	// logError = log.New(output, "❌ ", log.Ltime|log.Lshortfile)
	// logInfo = log.New(output, "🥨 ", log.Ltime)

	logDebug = log.New(vv104loggerDebug, "Debug: ", log.Ltime|log.Lshortfile)
	logError = log.New(vv104loggerError, "ERROR: ", log.Ltime|log.Lshortfile)
	logInfo = log.New(vv104loggerInfo, "", log.Ltime)

}

// func logInfof(v any) {
// 	logDebug.Println(v)

// 	// call callback function
// 	LogCallBack()

// }

// str, err := InfoLogBuf.ReadString(0x0a) // 0x0a = newline
