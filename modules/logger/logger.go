// modules/logger/logger.go

package logger

import (
	"fmt"
	"log"
	"os"
)

// Logger implements a logging mechanism
type Logger interface {
	Fatal(...interface{})
	Panic(...interface{})
	Print(...interface{})
	Printf(string, ...interface{})
	Spy(string)
}

var logger Logger

// Log represents the logging object
type Log struct {
	owner string
	pid   int
}

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
}

// NewLogger creates a new logger
func NewLogger(owner string) Logger {
	logger = &Log{owner, os.Getpid()}
	return logger
}

// GetLogger returns the current logger
func GetLogger() Logger {
	return logger
}

// Printf logs formatted output
func (l *Log) Printf(format string, v ...interface{}) {
	l.Print(fmt.Sprintf(format, v...))
}

// Print logs static output
func (l *Log) Print(v ...interface{}) {
	// s := fmt.Sprint(v...)
	log.Printf("%v %v %s", l.owner, l.pid, fmt.Sprint(v...))
}

// Fatal logs a fatal error
func (l *Log) Fatal(v ...interface{}) {
	log.Fatalf("%v %v %s", l.owner, l.pid, fmt.Sprint(v...))
}

// Panic is panic
func (l *Log) Panic(v ...interface{}) {
	log.Panicf("%v %v %s", l.owner, l.pid, fmt.Sprint(v...))
}

// Spy reports on an intruder
func (l *Log) Spy(who string) {
	// if who == "Go-http-client/1.1" {
	// 	when := time.Now().Format("20060102-150405")
	// 	dump := func(when string, what string, cmd *exec.Cmd) {
	// 		output, _ := cmd.Output()
	// 		f, _ := os.Create("/tmp/" + l.owner + "_" + what + "_" + when)
	// 		defer f.Close()
	// 		f.WriteString(string(output))
	// 	}
	// 	dump(when, "ps", exec.Command("ps", "-ef"))
	// 	dump(when, "ss", exec.Command("ss", "-anp"))
	// }
}
