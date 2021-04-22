package logger

import "testing"

func TestCreateFile(t *testing.T) {
	file := NewFile("log/test.log", 1, 1)
	file.createFile()
}
