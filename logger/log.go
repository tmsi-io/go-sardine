package logger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/buffer"
	"io"
	"os"
	"reflect"
	"sync"
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
	Out       io.Writer
	Level     Level
	Data      Fields
	Opts      *Option
	lock      sync.RWMutex
	Hooks     LevelHooks
	Formatter Formatter
	Buffer    *bytes.Buffer
	err       string
}

func New() *Logger {
	var l Logger
	options := DefaultOptions()
	l.Opts = options
	return &l
}

func DefaultOptions() *Option {
	return &Option{
		SaveTime: 3,
		MaxSize:  200,
		Level:    logrus.InfoLevel,
	}
}

func (l *Logger) level() Level {
	return l.Level
}

func (l *Logger) SetOptions(opts ...OptionFunc) {
	for _, opt := range opts {
		opt(l.Opts)
	}
	l.L.SetLevel(l.Opts.Level)
	if l.Opts.IsGradeOutput {
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

func (l *Logger) WithFiled(key string, value interface{}) *Logger {
	return l.WithFields(Fields{key: value})
}

func (l *Logger) WithFields(fields Fields) *Logger {
	data := make(Fields, len(l.Data)+len(fields))
	for k, v := range l.Data {
		data[k] = v
	}
	fieldErr := l.err
	for k, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch {
			case t.Kind() == reflect.Func, t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Func:
				isErrField = true
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", k)
			if fieldErr != "" {
				fieldErr = l.err + ", " + tmp
			} else {
				fieldErr = tmp
			}
		} else {
			data[k] = v
		}
	}
	l.Data = data
	return l
}

func (l *Logger) log(level Level, msg string) {
	var buffer *bytes.Buffer
	l.fireHooks()
	buffer = getBuffer()
	defer func() {
		l.Buffer = nil
		putBuffer(buffer)
	}()
	buffer.Reset()
	l.Buffer = buffer
	l.write()

}

func (l *Logger) write() {
	serialized, err := l.Formatter.Format(l)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}
	.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	if _, err := entry.Logger.Out.Write(serialized); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}

func (l *Logger) fireHooks() {
	var tmpHooks LevelHooks
	l.lock.Lock()
	tmpHooks = make(LevelHooks, len(l.Hooks))
	for k, v := range l.Hooks {
		tmpHooks[k] = v
	}
	l.lock.Unlock()

	err := tmpHooks.Fire(l.Level, l)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
	}
}

func (l *Logger) Remove() {

}
