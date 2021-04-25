package logger

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	log := New()
	f := NewFile("log/a.log", 10, 1)
	log.L.SetOutput(f)
	log.L.Infof("asdfadf")
	//basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	f1 := NewFile("log/b.log", 1,1)
	f2 := NewFile("log/c.log", 1,1)
	//writer1, _ := rotatelogs.New(
	//	"log/a1.log",
	//	rotatelogs.WithLinkName("log/b.log"),      // 生成软链，指向最新日志文件
	//	rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
	//	rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	//)
	//writer2, _ := rotatelogs.New(
	//	"log/a2.log",
	//	rotatelogs.WithLinkName("log/c.log"),      // 生成软链，指向最新日志文件
	//	rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
	//	rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	//)
	lfs := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel: f1,
		logrus.ErrorLevel: f2,
	}, &logrus.TextFormatter{})
	log.L.AddHook(lfs)
	go func(){
		for{
			log.L.Infof("1111111111111111")
			time.Sleep(time.Second)
		}
	}()
	go func(){
		log.L.Errorf("1111111111111111")
		time.Sleep(time.Second)
	}()
	var ss = make(chan int)
	<- ss
}