package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type File struct {
	FileName string
	MaxSize  int   //日志文件最大容量，超过则切片
	SaveTime int   //日志文件保存时间
	size     int64 //文件大小
	file     *os.File
	lock     *sync.Mutex
}

func (f *File) Write(p []byte) (n int, err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.file == nil {

	}

	n, err = f.file.Write(p)
	return n, err
}

func NewFile(fileName string, maxSize int, saveTime int) *File {
	return &File{
		FileName: fileName,
		MaxSize:  maxSize,
		SaveTime: saveTime,
	}

}

//返回文件绝对路径
func (f *File) dir() string {
	basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return filepath.Join(basePath, f.FileName)
}

//传入文件名生成目录文件，可带相对路径,e.g name:"test.log",name:"log/test.log"
//若传入为空，默认生成log/项目名.log目录文件
func (f *File) createFile() {
	if f.FileName == "" {
		f.FileName = fmt.Sprintf("log/%s", getProjectName())
	}
}

func getProjectName() string {
	proPath, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		s := strings.Split(proPath, "\\")
		return s[len(s)-1]
	} else {
		s := strings.Split(proPath, "/")
		return s[len(s)-1]
	}
}

//创建文件目录，若路径不包含目录则返回空，若路径包含目录则校验并创建目录
func (f *File) createDir() error {
	pathDir := filepath.Dir(f.FileName)
	fmt.Println(pathDir)
	basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirPath := filepath.Join(basePath, pathDir)
	_, err := os.Stat(dirPath)
	if err != nil {
		return os.MkdirAll(dirPath, 0744)
	}
	return nil
}


func (f *File) close() error {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f == nil {
		return nil
	}
	err := f.file.Close()
	f.file = nil
	return err
}
