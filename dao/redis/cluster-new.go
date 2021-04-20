package redis

import (
	"errors"
	_redis "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"go-sardine/logger"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ExtractURL
// extract redis url
func (cluster *Cluster) ExtractURL(_url string) error {
	if _url, err := url.Parse(_url); err != nil {
		return err
	} else {
		if args, err := url.ParseQuery(_url.RawQuery); err == nil {
			if maxRetry, ok := args[ConfigMaxRetries]; ok {
				cluster.maxRetries, _ = strconv.Atoi(maxRetry[0])
			}
			if PoolSize, ok := args[ConfigPoolSize]; ok {
				cluster.poolSize, _ = strconv.Atoi(PoolSize[0])
			}
		}
		if len(_url.Host) == 3 {
			return errors.New(ErrorConfigHost)
		} else {
			cluster.host = _url.Host
			cluster.addr = strings.Split(_url.Host, ",")
			if db := strings.Split(_url.Path, "/"); len(db) < 2 {
				return errors.New(ErrorConfigDB)
			} else {
				cluster.db, _ = strconv.Atoi(db[1])
			}
			cluster.password, _ = _url.User.Password()
		}
	}
	return nil
}

func (cluster *Cluster) Conn(_url string) error {
	cluster.logger = logger.GetLogger().WithFields(logrus.Fields{"URL": _url})
	cluster.maxRetries = 3
	cluster.dialTimeout = 3 * time.Second
	cluster.readTimeout = time.Second
	cluster.writeTimeout = time.Second
	if err := cluster.ExtractURL(_url); err != nil {
		return err
	}
	cluster.client = _redis.NewClusterClient(&_redis.ClusterOptions{
		Addrs:      cluster.addr,
		Password:   cluster.password,
		PoolSize:   cluster.poolSize,
		MaxRetries: cluster.maxRetries,
	})
	cluster.client.AddHook(cluster)
	if _, err := cluster.client.Ping().Result(); err != nil {
		return err
	} else {
		cluster.connOk = true
		go cluster.goKeepAlive()
		return nil
	}
}

func (cluster *Cluster) goKeepAlive() {
	for {
		if _, err := cluster.client.Ping().Result(); err != nil {
			cluster.connOk = false
			cluster.logger.Error(ErrorConnLoss)
			if err2 := cluster.InitConn(); err2 != nil {
				cluster.logger.Errorf("connection failed:[%v] ", err2)
			} else {
				cluster.connOk = true
				cluster.logger.Info("connection ok. ")
			}
		}
		time.Sleep(time.Second * 3)
	}
}

func (cluster *Cluster) InitConn() error {
	cluster.client = _redis.NewClusterClient(&_redis.ClusterOptions{
		Addrs:       cluster.addr,       //
		Password:    cluster.password,   // no password set
		MaxRetries:  cluster.maxRetries, //
		PoolSize:    cluster.poolSize,   //
		DialTimeout: CfgDialTimeout * time.Second,
		PoolTimeout: CfgPoolTimeout * time.Second,
		IdleTimeout: CfgIdleTimeout * time.Second,
	})
	if _, err := cluster.client.Ping().Result(); err != nil {
		return err
	} else {
		return nil
	}
}

func (cluster *Cluster) StatusOK() bool {
	return cluster.connOk
}
