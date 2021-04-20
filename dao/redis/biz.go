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

var businessMap redis

func init() {
	businessMap.rMap = make(map[string]Interface)
}

func RedisMap() *redis {
	return &businessMap
}

func (manager *redis) SetBusinessRedis(business string, URL string) Interface {
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

func (manager *redis) GetBusinessRedis(business string) Interface {
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

func (manager *redis) DelBusinessRedis(business string) {
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
	logger := logger.GetLogger().WithFields(logrus.Fields{"Redis": URL})
	for {
		if err := Redis.Conn(URL); err != nil {
			logger.Errorf("Dali Redis with URL Failed: [%v]", err)
			time.Sleep(3 * time.Second)
		} else {
			logger.Info("Connect to Redis OK.")
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
