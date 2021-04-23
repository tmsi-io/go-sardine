package logger

import (
	"github.com/sirupsen/logrus"
)


var _log logrus.Logger



func GetLogger(ops Option) *logrus.Logger {
	l := logrus.New()
	l.Formatter
}

