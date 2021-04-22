package utils

import (
	"os"
)

type IPList struct {
	IPv4 []string //
	IPv6 []string //	暂不处理
}

// GetHostName
// Get system Host name
func GetHostName() string {
	if name, err := os.Hostname(); err != nil {
		// can't get hostname, because in docker?
		ips, errIP := GetLocalIPList()
		if errIP != nil {
			return "error"
		} else {
			return ips.IPv4[0]
		}
	} else {
		return name
	}
}
