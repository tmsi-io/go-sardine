package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	//MB单位
	mByte      = 1024 * 1024
	backFormat = "0102_150405"
)

type File struct {
	FileName string
	MaxSize  int           //日志文件最大容量，超过则切片，为空默认200MB
	SaveTime int           //日志文件保存时间，为空默认保留三天
	Ext      func() string //日志文件备份后缀格式，可自定义方法返回后缀
	size     int64         //文件大小
	file     *os.File
	fileDir  string
	fileName string
	lock     sync.Mutex
}

func NewFile(fileName string, maxSize int, saveTime int) *File {
	return &File{
		FileName: fileName,
		MaxSize:  maxSize,
		SaveTime: saveTime,
	}

}

func (f *File) Write(p []byte) (n int, err error) {
	writeLen := int64(len(p))
	if writeLen > f.max() {
		return 0, fmt.Errorf(
			"write length %d more than max file size %d", writeLen, f.max(),
		)
	}

	if f.SaveTime == 0 {
		f.SaveTime = 3
	}
	if f.file == nil {
		f.existOrCreateFile()
	}
	if f.size+writeLen >= f.max() {
		err := f.rotate()
		return 0, err
	}
	n, err = f.file.Write(p)
	f.lock.Lock()
	defer f.lock.Unlock()
	f.size += int64(n)
	return n, err
}

//返回文件绝对路径
func (f *File) dir() string {
	basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return filepath.Join(basePath, f.FileName)
}

func absDir(name string) string {
	basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return filepath.Join(basePath, name)
}

func (f *File) max() int64 {
	if f.MaxSize == 0 {
		f.MaxSize = 200
	}
	return int64(f.MaxSize) * mByte
}

//检查文件是否存在，存在则写入文件对象，不存在则创建新对象
//传入文件名生成目录文件，可带相对路径,e.g name:"test.log",name:"log/test.log"
//若传入为空，默认生成log/项目名.log目录文件
func (f *File) existOrCreateFile() error {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.FileName == "" {
		f.FileName = fmt.Sprintf("log/%s", getProjectName())
	}
	proPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	f.fileDir = filepath.Join(proPath, filepath.Dir(f.FileName))
	f.fileName = filepath.Base(f.FileName)
	if err := f.createDir(); err != nil {
		return err
	}
	file := f.dir()
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return f.openFile()
	} else {
		f.openFile()
		f.size = info.Size()
	}
	return nil
}

func getProjectName() string {
	proPath, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		s := strings.Split(proPath, "/")
		return s[len(s)-1]
	} else {
		s := strings.Split(proPath, "/")
		return s[len(s)-1]
	}
}

//创建文件目录，若路径不包含目录则返回空，若路径包含目录则校验并创建目录
func (f *File) createDir() error {
	pathDir := filepath.Dir(f.FileName)
	basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirPath := filepath.Join(basePath, pathDir)
	_, err := os.Stat(dirPath)
	if err != nil {
		return os.MkdirAll(dirPath, 0744)
	}
	return nil
}

func (f *File) openFile() error {
	file, err := os.OpenFile(f.dir(), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		return err
	}
	f.file = file
	return nil
}

func (f *File) rotate() (err error) {
	if err = f.close(); err != nil {
		return err
	}
	if err = f.backup(); err != nil {
		return err
	}
	if err = f.existOrCreateFile(); err != nil {
		return err
	}
	return
}

func (f *File) backup() error {
	f.lock.Lock()
	defer f.lock.Unlock()
	newName := f.backupName()
	return os.Rename(f.dir(), absDir(newName))
}

//备份文件名，后缀可自定义，默认使用月日时分秒格式为后缀
func (f *File) backupName() string {
	ext := filepath.Ext(f.FileName)
	prefix := f.FileName[:len(f.FileName)-len(ext)]
	var bakExt string
	if f.Ext == nil {
		bakExt = time.Now().Format(backFormat)
	} else {
		bakExt = f.Ext()
	}
	return fmt.Sprintf("%s-%s%s", prefix, bakExt, ext)
}

func (f *File) close() error {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f == nil {
		return nil
	}
	err := f.file.Close()
	f.file = nil
	f.size = 0
	return err
}
