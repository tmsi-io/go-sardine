package redis

import (
	"errors"
	_redis "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"time"
)

type Sentinel struct {
	addr         []string
	host         string
	name         string
	username     string
	password     string
	db           int
	maxRetries   int
	poolSize     int
	client       *_redis.Client
	connOk       bool
	readOnly     bool
	logger       *logrus.Entry
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (s *Sentinel) StatusOK() bool {
	panic("implement me")
}

func (s *Sentinel) Close() error {
	return s.client.Close()
}

func (s *Sentinel) Get(key string) (string, error) {
	if s.client == nil {
		return "", errors.New(ErrorConnLoss)
	}
	return s.client.Get(key).Result()
}

func (s *Sentinel) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	if s.client == nil {
		return "", errors.New(ErrorConnLoss)
	}
	return s.client.Set(key, value, expiration).Result()
}

func (s *Sentinel) Del(keys ...string) (int64, error) {
	if s.client == nil {
		return 0, errors.New(ErrorConnLoss)
	}
	return s.client.Del(keys...).Result()
}

func (s *Sentinel) HGet(key string, filed string) (string, error) {
	if s.client == nil {
		return "", errors.New(ErrorConnLoss)
	}
	return s.client.HGet(key, filed).Result()
}

func (s *Sentinel) HGetAll(keys string) (map[string]string, error) {
	if s.client == nil {
		return nil, errors.New(ErrorConnLoss)
	}
	return s.client.HGetAll(keys).Result()
}

func (s *Sentinel) HSet(key string, filed string, value interface{}) (int64, error) {
	if s.client == nil {
		return 0, errors.New(ErrorConnLoss)
	}
	return s.client.HSet(key, filed, value).Result()
}

func (s *Sentinel) HDel(key string, filed string) (int64, error) {
	if s.client == nil {
		return 0, errors.New(ErrorConnLoss)
	}
	return s.client.HDel(key, filed).Result()
}

func (s *Sentinel) LLen(key string) (int64, error) {
	if s.client == nil {
		return 0, errors.New(ErrorConnLoss)
	}
	return s.client.LLen(key).Result()
}

// can't run in this mode
func (s *Sentinel) KEYS(key string) ([]string, error) {
	if s.client == nil {
		return nil, errors.New(ErrorConnLoss)
	}
	return s.client.Keys(key).Result()
}

func (s *Sentinel) Exists(key ...string) (int64, error) {
	if s.client == nil {
		return 0, errors.New(ErrorConnLoss)
	}
	return s.client.Exists(key...).Result()
}

func (s *Sentinel) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if s.client == nil {
		return false, errors.New(ErrorConnLoss)
	}
	return s.client.SetNX(key, value, expiration).Result()
}

func (s *Sentinel) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	if s.client == nil {
		return nil, 0, errors.New(ErrorConnLoss)
	}
	return s.client.HScan(key, cursor, match, count).Result()
}

func (s *Sentinel) PipeLineHSet(filed string, dataL map[string]interface{}) ([]_redis.Cmder, error) {
	if s.client == nil {
		return nil, errors.New(ErrorConnLoss)
	}
	pipeLine := s.client.Pipeline()
	for key, value := range dataL {
		pipeLine.HSet(filed, key, value)
	}
	return pipeLine.Exec()
}
