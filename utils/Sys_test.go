package utils

import (
	"fmt"
	"testing"
)

func TestGetHostName(t *testing.T) {
	host := GetHostName()
	fmt.Println(host)
}
