package logger

import "testing"

func TestCreateFile(t *testing.T) {
	file := NewFile("log/test.log", 1, 1)
	if err := file.existOrCreateFile();err != nil{
		t.Errorf("createFile err:%v", err)
	}
}

func TestBackup(t *testing.T){
	file := NewFile("log/test.log", 1, 1)
	if err := file.backup(); err != nil{
		t.Errorf("backup err:%v", err)
	}
}

func TestRotate(t *testing.T){
	file := NewFile("log/test.log", 1, 1)
	if err := file.existOrCreateFile(); err != nil{
		t.Errorf("createFile err:%v", err)
	}
	if err := file.rotate();err != nil{
		t.Errorf("rotate err:%v", err)
	}
}
