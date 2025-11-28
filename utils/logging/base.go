package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

var IsVerbose = false
var IsDebug = false

func Init(base string, suffix string) {
	execPath, _ := os.Executable()
	// Setup lumberjack
	logFile := &lumberjack.Logger{
		Filename:   base + "/" + filepath.Base(execPath) + "_" + suffix + ".log",
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
