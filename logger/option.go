package logger

import "strings"

type option struct {
	FileName      string
	SaveTime      int
	MaxSize       int
	Format        func() string
	IsSplit       bool
	IsGradeOutput bool
}

type Option func(opt *option)

//日志文件备份后缀格式，可自定义方法返回后缀
func RotateFormat(f func() string) Option {
	if strings.Contains(f(), ":") {
		panic("ext format cannot contain ':'")
	}
	return func(opt *option) {
		opt.Format = f
	}
}

func FileName(name string) Option {
	return func(opt *option) {
		opt.FileName = name
	}
}

//日志文件保存时间，默认保留三天
func SaveTime(s int) Option {
	return func(opt *option) {
		opt.SaveTime = s
	}
}

//日志切割大小，默认200M
func MaxSize(m int) Option {
	return func(opt *option) {
		opt.MaxSize = m
	}
}

//是否日志分片，默认false
func IsSplit(b bool) Option {
	return func(opt *option) {
		opt.IsSplit = b
	}
}

//是否按照日志等级分别输出文件，默认false
func IsGradeOutput(b bool) Option {
	return func(opt *option) {
		opt.IsGradeOutput = b
	}
}
