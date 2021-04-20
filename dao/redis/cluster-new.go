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
	if pUrl, err := url.Parse(_url); err != nil {
		return err
	} else {
		if args, err := url.ParseQuery(pUrl.RawQuery); err == nil {
			if maxRetry, ok := args["MaxRetries"]; ok {
				mR, _ := strconv.Atoi(maxRetry[0])
				cluster.maxRetries = mR
			}
			if PoolSize, ok := args["PoolSize"]; ok {
				pS, _ := strconv.Atoi(PoolSize[0])
				cluster.poolSize = pS
			}
		}
		if len(pUrl.Host) == 3 {
			return errors.New("Can't Parse HostPort Config. ")
		} else {
			cluster.host = pUrl.Host
			cluster.addr = strings.Split(pUrl.Host, ",")
			if db := strings.Split(pUrl.Path, "/"); len(db) < 2 {
				return errors.New("Can't Parse DB Config. ")
			} else {
				cluster.db, _ = strconv.Atoi(db[1])
			}
			cluster.password, _ = pUrl.User.Password()
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
		Addrs:      cluster.addr,     //set redis cluster url
		Password:   cluster.password, //set password
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
			cluster.logger.Error("Connection Loss, ReBuilding... ")
			if err2 := cluster.InitConn(); err2 != nil {
				cluster.logger.Error("Connection ReBuilding Failed: ", err2)
			} else {
				cluster.connOk = true
				cluster.logger.Info("Connection OK, ReBuilding ok ")
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
		DialTimeout: 10 * time.Second,
		PoolTimeout: 1 * time.Second,
		IdleTimeout: 20 * time.Second,
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
