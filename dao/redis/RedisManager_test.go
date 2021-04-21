package redis

import (
	"fmt"
	"testing"
	"time"
)

func TestRedisManager(t *testing.T) {
	// Test for single redis
	Pool().SetConn("buy", "redis://:password@192.168.0.2:6379/0")
	_, errSet := Pool().GetConn("buy").Set("prod", "1", 1*time.Second)
	if errSet != nil {
		fmt.Printf(" set prod failed : %s \n", errSet)
	}
	val, errGet := Pool().GetConn("buy").Get("prod")
	if errGet != nil {
		fmt.Printf(" prod was : %s \n", val)
	}
	// Test for cluster redis
	cluster := Pool().SetConn("auth", "redis://:password@192.168.2.2:6379,192.168.2.3:6379/0")
	cluster.Set("user1", "something", 1*time.Minute)
	cluster.Get("user1")
	Pool().SetConn("cache", "redis-sentinel://:password@192.168.2.1:26379, 192.168.2.2:26379/mymaster/0")
	Pool().GetConn("cache").Set("a", "b", 0)
}
