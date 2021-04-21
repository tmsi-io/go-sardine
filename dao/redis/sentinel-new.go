package redis

import (
	"errors"
	_redis "github.com/go-redis/redis/v7"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ExtractURL
// extract redis url
func (s *Sentinel) ExtractURL(_url string) error {
	if _url, err := url.Parse(_url); err != nil {
		return err
	} else {
		if args, err := url.ParseQuery(_url.RawQuery); err == nil {
			if maxRetry, ok := args[ConfigMaxRetries]; ok {
				s.maxRetries, _ = strconv.Atoi(maxRetry[0])
			}
			if PoolSize, ok := args[ConfigPoolSize]; ok {
				s.poolSize, _ = strconv.Atoi(PoolSize[0])
			}
		}
		if len(_url.Host) == 3 {
			return errors.New(ErrorConfigHost)
		} else {
			s.host = _url.Host
			s.addr = strings.Split(_url.Host, ",")
			if args := strings.Split(_url.Path, "/"); len(args) < 2 {
				return errors.New(ErrorConfigDB)
			} else {
				s.db, _ = strconv.Atoi(args[2])
				s.name = args[1]
				if s.name == "" {
					s.name = DefaultSentinel
				}
			}
			s.password, _ = _url.User.Password()
			s.username = _url.User.Username()
			s.name = _url.Path
		}
	}
	return nil
}

func (s *Sentinel) Conn(_url string) error {
	if err := s.ExtractURL(_url); err != nil {
		return err
	}
	_ = s.InitConn()
	s.client.AddHook(s)
	if _, err := s.client.Ping().Result(); err != nil {
		return err
	} else {
		s.connOk = true
		go s.goKeepAlive()
		return nil
	}
}

func (s *Sentinel) goKeepAlive() {
	for {
		if _, err := s.client.Ping().Result(); err != nil {
			s.connOk = false
			s.logger.Error(ErrorConnLoss)
			if err2 := s.InitConn(); err2 != nil {
				s.logger.Errorf(ErrorConnLoss)
			} else {
				s.connOk = true
				s.logger.Info("connection ok. ")
			}
		}
		time.Sleep(time.Second * 3)
	}
}

func (s *Sentinel) InitConn() error {
	s.client = _redis.NewFailoverClient(&_redis.FailoverOptions{
		MasterName:       s.name,
		SentinelAddrs:    s.addr,
		SentinelUsername: s.username,
		SentinelPassword: s.password,
		PoolSize:         s.poolSize,
		DialTimeout:      CfgDialTimeout * time.Second,
		PoolTimeout:      CfgPoolTimeout * time.Second,
		IdleTimeout:      CfgIdleTimeout * time.Second,
	})
	if _, err := s.client.Ping().Result(); err != nil {
		return err
	}
	return nil
}
