package logger

import (
	"github.com/sirupsen/logrus"
)

var Loggerlevels = [7]logrus.Level{logrus.TraceLevel, logrus.DebugLevel, logrus.InfoLevel,
	logrus.DebugLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}

func (l *Logger) AddLfsHook(){
//	var lvIdx int
//	for i, lv := range Loggerlevels{
//		if l.Level == lv{
//			lvIdx = i
//		}
//	}
//	lvs := Loggerlevels[lvIdx:]
//	hookMap := make(lfshook.WriterMap)
//	for _, lv := range lvs{
//		hookMap[lv] = io.Writer()
//	}
//	lfHook := lfshook.NewHook(hookMap, l.Formatter)
//	l.Hooks.Add(lfHook)
}
