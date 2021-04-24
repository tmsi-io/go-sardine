package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)


var _log logrus.Logger



func GetLogger(opts ...OptionFunc) *logrus.Logger {
	l := logrus.New()
	options := DefaultOptions()
	for _, opt := range opts{
		opt(options)
	}
	if options.FileName != ""{
		os.Open()
		l.SetOutput()
	}
}

func DefaultOptions() *Option {
	return &Option{
		SaveTime: 3,
		MaxSize: 200,
	}
}

