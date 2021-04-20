package redis

import (
	redisv7 "github.com/go-redis/redis/v7"
	"time"
)

type Interface interface {
	Conn(_url string) error
	StatusOK() bool
	Close() error
	HSet(key string, filed string, value interface{}) (int64, error)
	HDel(key string, filed string) (int64, error)
	HGet(key string, filed string) (string, error)
	HGetAll(keys string) (map[string]string, error)
	LLen(key string) (int64, error)
	KEYS(key string) ([]string, error)
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) (string, error)
	Del(keys ...string) (int64, error)
	Exists(key ...string) (int64, error)
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error)
	PipeLineHSet(filed string, datas map[string]interface{}) ([]redisv7.Cmder, error)
	GetClient() (interface{}, error)
}
