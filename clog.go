// Package clog implements an alternative logger to the one found in the standard
// library with support for more logging levels and a different output format.
// It also has support for splitting log files on daily boundaries.
//
// Author: Clint Caywood
//
// https://github.com/cratonica/clog
package clog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Represents how critical the logged
// message is.
type Level uint8

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

var levelStrings = map[Level]string{
	LevelTrace:   "TRACE",
	LevelDebug:   "DEBUG",
	LevelInfo:    "INFO",
	LevelWarning: "WARNING",
	LevelError:   "ERROR",
	LevelFatal:   "FATAL",
}

func (this Level) String() string {
	result := levelStrings[this]
	if len(result) == 0 {
		return fmt.Sprintf("Unknown Level: %d", this)
	}
	return result
}

type output struct {
	writer   io.Writer
	levelMin Level
	levelMax Level
}

// The Logger
type Clog struct {
	mtx     sync.Mutex
	outputs []output
	// Terminate the program on Fatal, false by default.
	ExitOnFatal bool
}

// Instantiate a new Clog
func NewClog() *Clog {
	return &Clog{sync.Mutex{}, make([]output, 0), false}
}

// Adds an ouput, specifying the minimum log Level
// you want to be written to this output. For instance,
// if you pass Warning for level, all logs of type
// Warning, Error, and Fatal would be logged to this output.
func (this *Clog) AddOutput(writer io.Writer, level Level) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.outputs = append(this.outputs, output{writer, level, LevelFatal})
}

// Adds an ouput, specifying the minimum and maximum log Level
// you want to be written to this output. For instance,
// if you pass Trace for levelMin and Warning for levelMax, all logs of type
// Trace, Debug, Info and Warning would be logged to this output.
func (this *Clog) AddOutputRange(writer io.Writer, levelMin, levelMax Level) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.outputs = append(this.outputs, output{writer, levelMin, levelMax})
}

// Convenience function
func (this *Clog) Trace(format string, v ...interface{}) {
	this.Log(LevelTrace, format, v...)
}

// Convenience function
func (this *Clog) Debug(format string, v ...interface{}) {
	this.Log(LevelDebug, format, v...)
}

// Convenience function
func (this *Clog) Info(format string, v ...interface{}) {
	this.Log(LevelInfo, format, v...)
}

// Convenience function
func (this *Clog) Warning(format string, v ...interface{}) {
	this.Log(LevelWarning, format, v...)
}

// Convenience function
func (this *Clog) Error(format string, v ...interface{}) {
	this.Log(LevelError, format, v...)
}

// Convenience function, will terminate the program if you set
// ExitOnFatal to true.
func (this *Clog) Fatal(format string, v ...interface{}) {
	this.Log(LevelFatal, format, v...)
	if this.ExitOnFatal {
		os.Exit(1)
	}
}

// Logs a message for the given level. Most callers will likely
// prefer to use one of the provided convenience functions.
func (this *Clog) Log(level Level, format string, v ...interface{}) {
	message := fmt.Sprintf(format+"\n", v...)
	strTimestamp := getTimestamp()
	strFinal := fmt.Sprintf(
		"%s [%-7s] %s", strTimestamp, levelStrings[level], message)
	bytes := []byte(strFinal)
	this.mtx.Lock()
	defer this.mtx.Unlock()
	for _, output := range this.outputs {
		if output.levelMin <= level && level <= output.levelMax {
			output.writer.Write(bytes)
		}
	}
}

// Gets the timestamp string
func getTimestamp() string {
	now := time.Now()
	return fmt.Sprintf(
		"%v-%02d-%02d %02d:%02d:%02d.%03d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000000)
}
