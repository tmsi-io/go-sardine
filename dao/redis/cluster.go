package redis

import (
	"errors"
	_redis "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"time"
)

type Cluster struct {
	addr         []string
	host         string
	password     string
	db           int
	maxRetries   int
	poolSize     int
	client       *_redis.ClusterClient
	connOk       bool
	readOnly     bool
	logger       *logrus.Entry
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

var clusterKey = struct{}{}

func (cluster *Cluster) KEYS(key string) ([]string, error) {
	if cluster.client == nil {
		return nil, errors.New("Connection_Loss")
	}
	return cluster.client.Keys(key).Result()
}

func (cluster *Cluster) Get(key string) (string, error) {
	if cluster.client == nil {
		return "", errors.New("Connection_Loss")
	}
	return cluster.client.Get(key).Result()
}

func (cluster *Cluster) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	if cluster.client == nil {
		return "", errors.New("Connection Loss. ")
	}
	return cluster.client.Set(key, value, expiration).Result()
}

func (cluster *Cluster) Close() error {
	if cluster.client == nil {
		return errors.New("Connection Loss. ")
	}
	return cluster.client.Close()
}

func (cluster *Cluster) Del(keys ...string) (int64, error) {
	if cluster.client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return cluster.client.Del(keys...).Result()
}

func (cluster *Cluster) Exists(key ...string) (int64, error) {
	if cluster.client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return cluster.client.Exists(key...).Result()
}

func (cluster *Cluster) HGetAll(keys string) (map[string]string, error) {
	if cluster.client == nil {
		return nil, errors.New("Connection Loss. ")
	}
	return cluster.client.HGetAll(keys).Result()
}

func (cluster *Cluster) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if cluster.client == nil {
		return false, errors.New("Connection Loss. ")
	}
	return cluster.client.SetNX(key, value, expiration).Result()
}
func (cluster *Cluster) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	if cluster.client == nil {
		return nil, 0, errors.New("Connection Loss. ")
	}
	return cluster.client.HScan(key, cursor, match, count).Result()
}
func (cluster *Cluster) HSet(key string, filed string, value interface{}) (int64, error) {
	if cluster.client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return cluster.client.HSet(key, filed, value).Result()
}

func (cluster *Cluster) HGet(key string, filed string) (string, error) {
	if cluster.client == nil {
		return "", errors.New("Connection_Loss")
	}
	return cluster.client.HGet(key, filed).Result()
}

func (cluster *Cluster) HDel(key string, filed string) (int64, error) {
	if cluster.client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return cluster.client.HDel(key, filed).Result()
}

func (cluster *Cluster) PipeLineHSet(filed string, dataL map[string]interface{}) ([]_redis.Cmder, error) {
	if cluster.client == nil {
		return nil, errors.New("Connection Loss. ")
	}
	pipeLine := cluster.client.Pipeline()
	for key, value := range dataL {
		pipeLine.HSet(filed, key, value)
	}
	return pipeLine.Exec()
}

func (cluster *Cluster) LLen(key string) (int64, error) {
	if cluster.client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return cluster.client.LLen(key).Result()
}
