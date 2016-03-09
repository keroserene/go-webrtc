package webrtc

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	INFO  *log.Logger
	WARN  *log.Logger
	ERROR *log.Logger
	TRACE *log.Logger
)

// Logging verbosity level, from 0 (nothing) upwards.
func SetLoggingVerbosity(level int) {
	// handle io.Writer
	infoOut := ioutil.Discard
	warnOut := ioutil.Discard
	errOut := ioutil.Discard
	traceOut := ioutil.Discard

	// TODO: Better logging levels
	if level > 0 {
		errOut = os.Stdout
	}
	if level > 1 {
		warnOut = os.Stdout
	}
	if level > 2 {
		infoOut = os.Stdout
	}
	if level > 3 {
		traceOut = os.Stdout
	}

	INFO = log.New(infoOut,
		"INFO: ",
		// log.Ldate|log.Ltime|log.Lshortfile)
		log.Lshortfile)
	WARN = log.New(warnOut,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	ERROR = log.New(errOut,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	TRACE = log.New(traceOut,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
