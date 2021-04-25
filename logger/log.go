package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var _log logrus.Logger

type Logger struct {
	logrus.Logger
	Opts *Option
	lock sync.RWMutex
}

func New() *Logger {
	var l = new(Logger)
	options := DefaultOptions()
	l.SetOutput(os.Stderr)
	l.SetLevel(options.Level)
	l.SetReportCaller(true)
	l.ExitFunc = os.Exit
	l.ReportCaller = false
	l.Hooks = make(logrus.LevelHooks)
	l.Opts = options
	return l
}

func DefaultOptions() *Option {
	return &Option{
		SaveTime: 3,
		MaxSize:  200,
		Level: logrus.InfoLevel,
	}
}

func (l *Logger) SetOptions(opts ...OptionFunc) {
	for _, opt := range opts {
		opt(l.Opts)
	}
	l.SetLevel(l.Opts.Level)
	if l.Opts.IsGradeOutput{
		l.AddLfsHook()
	}
	//l.Opts = options
}


func (l *Logger) Remove() {

}
