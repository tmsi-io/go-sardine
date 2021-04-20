package logger

import (
	"github.com/sirupsen/logrus"
)

var _log logrus.Logger

func GetLogger() *logrus.Logger {
	return &_log
}
