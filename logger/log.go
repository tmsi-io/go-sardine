package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var _log logrus.Logger

func New() *logrus.Logger{
	return &logrus.Logger{
		Out:          os.Stderr,
	}
}

func GetLogger() *logrus.Logger {
	return &_log
}
