package logger

import (
	"github.com/sirupsen/logrus"
	"strings"
)

type Option struct {
	FileName      string
	SaveTime      int
	MaxSize       int
	Format        func() string
	Level         logrus.Level
	Formatter     logrus.Formatter
	//IsSplit       bool
	IsGradeOutput bool
}

type OptionFunc func(opt *Option)

//日志文件备份后缀格式，可自定义方法返回后缀
func RotateFormat(f func() string) OptionFunc {
	if strings.Contains(f(), ":") {
		panic("ext format cannot contain ':'")
	}
	return func(opt *Option) {
		opt.Format = f
	}
}

func FileName(name string) OptionFunc {
	return func(opt *Option) {
		opt.FileName = name
	}
}

//日志文件保存时间，默认保留三天
func SaveTime(s int) OptionFunc {
	return func(opt *Option) {
		opt.SaveTime = s
	}
}

//日志切割大小，默认200M
func MaxSize(m int) OptionFunc {
	return func(opt *Option) {
		opt.MaxSize = m
	}
}

//是否日志分片，默认false
//func IsSplit(b bool) OptionFunc {
//	return func(opt *Option) {
//		opt.IsSplit = b
//	}
//}

//是否按照日志等级分别输出文件，默认false
func IsGradeOutput(b bool) OptionFunc {
	return func(opt *Option) {
		opt.IsGradeOutput = b
	}
}

//设置日志等级
func Level(lv logrus.Level) OptionFunc{
	return func(opt *Option) {
		opt.Level = lv
	}
}

func Formater(format logrus.Formatter) OptionFunc{
	return func(opt *Option) {
		opt.Formatter = format
	}
}