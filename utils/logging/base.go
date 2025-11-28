package logging

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var IsVerbose = false
var IsDebug = false

func Init(base string) {
	// Setup lumberjack
	logFile := &lumberjack.Logger{
		Filename:   base + "/BinanceBuddy.log",
		MaxSize:    200,
		MaxBackups: 50,
		MaxAge:     14,
		Compress:   true,
	}

	// Output to both stdout and file.
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// setup system log to use lumberjack
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags) // print date.

	log.Println("[logging] Tide logging initialized. path: " + base)
}

func Verbosef(format string, a ...any) {
	if IsVerbose {
		log.Printf("[VERBOSE] "+format, a...)
	}
}

func Debugf(format string, a ...any) {
	if IsDebug {
		log.Printf("[DEBUG] "+format, a...)
	}
}

func Logf(format string, a ...any) {
	log.Printf(format, a...)
}

func Errorf(format string, a ...any) {
	log.Panicf(format, a...)
}
