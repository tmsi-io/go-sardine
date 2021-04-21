package redis

import (
	"errors"
	redisv7 "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/tmsi-io/go-sardine/logger"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (single *Single) Conn(_url string) error {
	single.logger = logger.GetLogger().WithFields(logrus.Fields{
		"Config": _url,
	})
	if pUrl, err := url.Parse(_url); err != nil {
		return err
	} else {
		if args, err := url.ParseQuery(pUrl.RawQuery); err == nil { // 解析url对象
			if maxRetry, ok := args[ConfigMaxRetries]; ok {
				mR, _ := strconv.Atoi(maxRetry[0])
				single.maxRetries = mR
			}
			if PoolSize, ok := args[ConfigPoolSize]; ok {
				pS, _ := strconv.Atoi(PoolSize[0])
				single.poolSize = pS
			}
		}
		if len(pUrl.Host) == 3 {
			return errors.New(ErrorConfigHost)
		} else {
			single.addr = pUrl.Host
			if db := strings.Split(pUrl.Path, "/"); len(db) < 2 {
				return errors.New(ErrorConfigDB)
			} else {
				single.db, _ = strconv.Atoi(db[1])
			}
			single.password, _ = pUrl.User.Password()
		}
	}
	if err := single.InitConn(); err != nil {
		return err
	} else {
		single.connOk = true
		go single.goKeepAlive() // 发送心跳保活
		return nil
	}
}

func (single *Single) InitConn() error {
	client := redisv7.NewClient(&redisv7.Options{
		Addr:        single.addr,       //
		Password:    single.password,   // no password set
		DB:          single.db,         // use default DB
		MaxRetries:  single.maxRetries, //
		PoolSize:    single.poolSize,   //
		DialTimeout: 10 * time.Second,
		PoolTimeout: 1 * time.Second,
		IdleTimeout: 20 * time.Second,
	})
	single.client = client
	single.client.AddHook(single)
	if _, err := single.client.Ping().Result(); err != nil {
		single.logger.Error(err)
		return err
	} else {
		single.logger.Info("Conn To Redis OK.")
		return nil
	}
}

func (single *Single) StatusOK() bool {
	return single.connOk
}

func (single *Single) goKeepAlive() {
	for {
		if _, err := single.client.Ping().Result(); err != nil {
			single.connOk = false
			single.logger.Error("Connection Loss, ReBuilding... ")
			if err2 := single.InitConn(); err2 != nil {
				single.logger.Error("Connection ReBuilding Failed: ", err2)
			} else {
				single.connOk = true
				single.logger.Info("Connection OK, ReBuilding ok ")
			}
			time.Sleep(3 * time.Second)
		}
		time.Sleep(time.Second * 10)
	}
}
