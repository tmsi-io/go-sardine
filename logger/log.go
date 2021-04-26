package logger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

type Fields map[string]interface{}

type logger interface {
	WithField(key string, value interface{}) *logger
	WithFields(fields Fields) *logger
	WithError(err error) *logger

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

type Logger struct {
	Out io.Writer
	L *logrus.Logger
	Level Level
	Opts *Option
	lock sync.RWMutex
}

func New() *Logger {
	var l Logger
	l.L = logrus.New()
	options := DefaultOptions()
	l.L.SetOutput(os.Stderr)
	l.L.SetLevel(options.Level)
	l.L.SetReportCaller(true)
	l.L.ExitFunc = os.Exit
	l.L.ReportCaller = false
	l.L.Hooks = make(logrus.LevelHooks)
	l.Opts = options
	return &l
}

func DefaultOptions() *Option {
	return &Option{
		SaveTime: 3,
		MaxSize:  200,
		Level: logrus.InfoLevel,
	}
}

func (l *Logger) level ()Level{
	return l.le
}

func (l *Logger) SetOptions(opts ...OptionFunc) {
	for _, opt := range opts {
		opt(l.Opts)
	}
	l.L.SetLevel(l.Opts.Level)
	if l.Opts.IsGradeOutput{
		l.AddLfsHook()
	}
	//l.Opts = options
}

func (l *Logger) Log(level Level, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		l.log(level, fmt.Sprint(args...))
	}
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (l *Logger) IsLevelEnabled(level Level) bool {
	return logger.level() >= level
}

func (l *Logger) log(level Level, msg string) {
	var buffer *bytes.Buffer

	newEntry := entry.Dup()

	if newEntry.Time.IsZero() {
		newEntry.Time = time.Now()
	}

	newEntry.Level = level
	newEntry.Message = msg

	newEntry.Logger.mu.Lock()
	reportCaller := newEntry.Logger.ReportCaller
	newEntry.Logger.mu.Unlock()

	if reportCaller {
		newEntry.Caller = getCaller()
	}

	newEntry.fireHooks()

	buffer = getBuffer()
	defer func() {
		newEntry.Buffer = nil
		putBuffer(buffer)
	}()
	buffer.Reset()
	newEntry.Buffer = buffer

	newEntry.write()

	newEntry.Buffer = nil

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(newEntry)
	}
}


func (l *Logger) Remove() {

}
