package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type logger interface {

}

type Logger struct {
	L *logrus.Logger
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


func (l *Logger) Remove() {

}
