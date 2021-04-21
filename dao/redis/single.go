package redis

import (
	"errors"
	redisv7 "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"time"
)

type Single struct {
	addr       string
	password   string
	db         int
	maxRetries int
	poolSize   int
	client     *redisv7.Client
	connOk     bool
	ReadOnly   bool
	logger     *logrus.Entry
}

var singleRedis Single

func GetSingleRedis() *Single {
	return &singleRedis
}

var key = struct{}{}

func (single *Single) KEYS(key string) ([]string, error) {
	if single.client == nil {
		return nil, errors.New(ErrorConnLoss)
	}
	return single.client.Keys(key).Result()
}

func (single *Single) Get(key string) (string, error) {
	if single.client == nil {
		return "", errors.New(ErrorConnLoss)
	}
	return single.client.Get(key).Result()
}

func (single *Single) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	if single.client == nil {
		return "", errors.New(ErrorConnLoss)
	}
	return single.client.Set(key, value, expiration).Result()
}

func (single *Single) Del(keys ...string) (int64, error) {
	if single.client == nil {
		return 0, errors.New(ErrorConnLoss)
	}
	return single.client.Del(keys...).Result()
}

func (single *Single) Close() error {
	if single.client == nil {
		return errors.New(ErrorConnLoss)
	}
	return single.client.Close()
}

func (single *Single) Exists(key ...string) (int64, error) {
	if single.client == nil {
		return 0, errors.New(ErrorConnLoss)
	}
	return single.client.Exists(key...).Result()
}

func (single *Single) HGetAll(keys string) (map[string]string, error) {
	if single.client == nil {
		return nil, errors.New(ErrorConnLoss)
	}
	return single.client.HGetAll(keys).Result()
}

func (single *Single) LLen(key string) (int64, error) {
	if single.client == nil {
		return 0, errors.New(ErrorConnLoss)
	}
	return single.client.LLen(key).Result()
}

func (single *Single) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if single.client == nil {
		return false, errors.New(ErrorConnLoss)
	}
	return single.client.SetNX(key, value, expiration).Result()
}

func (single *Single) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	if single.client == nil {
		return nil, 0, errors.New(ErrorConnLoss)
	}
	return single.client.HScan(key, cursor, match, count).Result()
}

func (single *Single) HSet(key string, filed string, value interface{}) (int64, error) {
	if single.client == nil {
		return 0, errors.New(ErrorConnLoss)
	}
	return single.client.HSet(key, filed, value).Result()
}

func (single *Single) HGet(key string, filed string) (string, error) {
	if single.client == nil {
		return "", errors.New(ErrorConnLoss)
	}
	return single.client.HGet(key, filed).Result()
}

func (single *Single) HDel(key string, filed string) (int64, error) {
	if single.client == nil {
		return 0, errors.New(ErrorConnLoss)
	}
	return single.client.HDel(key, filed).Result()
}

func (single *Single) PipeLineHSet(filed string, dataL map[string]interface{}) ([]redisv7.Cmder, error) {
	if single.client == nil {
		return nil, errors.New(ErrorConnLoss)
	}
	pipeLine := single.client.Pipeline()
	for key, value := range dataL {
		pipeLine.HSet(filed, key, value)
	}
	return pipeLine.Exec()
}

func (single *Single) GetClient() (interface{}, error) {
	if single.client == nil {
		return nil, errors.New(ErrorConnLoss)
	}
	return single.client, nil
}
