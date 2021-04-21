package redis

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tmsi-io/go-sardine/logger"
	"net/url"
	"strings"
	"sync"
	"time"
)

// RedisBiz
type redis struct {
	rMap  map[string]Interface
	rLock sync.RWMutex
}

var bizMap redis

func init() {
	bizMap.rMap = make(map[string]Interface)
}

// RedisManager
// manager all redis by biz name
func Pool() *redis {
	return &bizMap
}

func (manager *redis) SetConn(business string, URL string) Interface {
	/*
		新建业务链接
	*/
	manager.rLock.RLock()
	if existRedis, ok := manager.rMap[business]; ok {
		manager.rLock.RUnlock()
		return existRedis
	} else {
		manager.rLock.RUnlock()
		log := logger.GetLogger().WithFields(logrus.Fields{"Redis": URL, "Business": business})
		log.Info("Start Add New Redis Conn")
		newRedis := manager.getNewRedisObj(URL)
		manager.rLock.Lock()
		manager.rMap[business] = newRedis
		manager.rLock.Unlock()
		go manager.dialWithURL(URL, newRedis)
		return newRedis
	}
}

func (manager *redis) GetConn(business string) Interface {
	/*
		获取业务链接
	*/
	manager.rLock.RLock()
	if existRedis, ok := manager.rMap[business]; ok {
		manager.rLock.RUnlock()
		return existRedis
	}
	manager.rLock.RUnlock()
	return nil
}

func (manager *redis) DelCon(business string) {
	/*
		从业务Map中删除并关闭Redis链接
	*/
	manager.rLock.Lock()
	if existRedis, ok := manager.rMap[business]; ok {
		_ = existRedis.Close()
		delete(manager.rMap, business)
	}
	manager.rLock.Unlock()
}

func (manager *redis) getNewRedisObj(url string) Interface {
	if manager.URlWasCluster(url) {
		return new(Cluster)
	} else {
		return new(Single)
	}
}

func (manager *redis) dialWithURL(URL string, Redis Interface) {
	log := logger.GetLogger().WithFields(logrus.Fields{"Redis": URL})
	for {
		if err := Redis.Conn(URL); err != nil {
			log.Errorf("Dali Redis with URL Failed: [%v]", err)
			time.Sleep(3 * time.Second)
		} else {
			log.Info("Connect to Redis OK.")
			break
		}
	}
}

func (manager *redis) URlWasCluster(_url string) bool {
	if pUrl, err := url.Parse(_url); err != nil {
		fmt.Println(err)
		return false
	} else {
		if len(strings.Split(pUrl.Host, ",")) > 1 {
			return true
		} else {
			return false
		}
	}
}

func (manager *redis) URLWasSentinel(_url string) bool {
	if pUrl, err := url.Parse(_url); err != nil {
		fmt.Println(err)
		return false
	} else {
		if pUrl.Scheme == SentinelAddr1 || pUrl.Scheme == SentinelAddr2 {
			return true
		} else {
			return false
		}
	}
}
