// modules/logger/logger.go

package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	defaultDepth = 2
)

type Level int

const (
	None     Level = iota
	Debug          = iota + 10
	Info           = iota + 10
	Warning        = iota + 10
	Error          = iota + 10
	Critical       = iota + 10
)

// Logger implements a logging mechanism
type Logger interface {
	GetLevel() Level
	Output(level Level, depth int, s string) error
	Log(level Level, v ...interface{}) error
	Logf(level Level, format string, v ...interface{}) error
	Logd(level Level, delta int, v ...interface{}) error
	Logdf(level Level, delta int, format string, v ...interface{}) error
}

var logger Logger

// Log represents the logging object
type Log struct {
	owner string
	level Level
	pid   int
}

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
}

// NewLogger creates a new logger
func NewLogger(owner string, level Level) Logger {
	if level < 0 {
		level = 0
	} else if level > 99 {
		level = 99
	}
	if level < Info {
		log.SetFlags(log.Flags() | log.Lshortfile)
	}
	logger = &Log{owner, level, os.Getpid()}
	return logger
}

// GetLogger returns the current logger
func GetLogger() Logger {
	return logger
}

func (level Level) ToString() string {
	var levelName string
	switch level {
	case Critical:
		levelName = "CRITICAL"
	case Error:
		levelName = "ERROR"
	case Warning:
		levelName = "WARNING"
	case Info:
		levelName = "INFO"
	case Debug:
		levelName = "DEBUG"
	case None:
		levelName = "NONE"
	default:
		levelName = fmt.Sprintf("CUSTOM%02d", level)
	}
	return levelName
}

// GetLevel returns the log level
func (l *Log) GetLevel() Level {
	return l.level
}

// Output writes the output for a logging event
func (l *Log) Output(level Level, depth int, s string) error {
	if level < l.level {
		return nil
	}
	return log.Output(depth+1, fmt.Sprintf("%-8s %s %v %s", level.ToString(), l.owner, l.pid, s))
}

// Log calls Output to print to the logger
func (l *Log) Log(level Level, v ...interface{}) error {
	return l.Output(level, defaultDepth, fmt.Sprint(v...))
}

// Logf calls Output to print to the logger
func (l *Log) Logf(level Level, format string, v ...interface{}) error {
	return l.Output(level, defaultDepth, fmt.Sprintf(format, v...))
}

// Logd calls Output with variable depth to print to the logger
func (l *Log) Logd(level Level, delta int, v ...interface{}) error {
	return l.Output(level, defaultDepth+delta, fmt.Sprint(v...))
}

// Logfd calls Output with variable depth to print to the logger
func (l *Log) Logdf(level Level, delta int, format string, v ...interface{}) error {
	return l.Output(level, defaultDepth+delta, fmt.Sprintf(format, v...))
}
