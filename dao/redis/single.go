package redis

import (
	"context"
	"errors"
	redisv7 "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"go-sardine/logger"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Single struct {
	addr       string
	Password   string
	DB         int
	MaxRetries int
	PoolSize   int
	Client     *redisv7.Client
	ConnOk     bool
	ReadOnly   bool
	logger     *logrus.Entry
}

var singleRedis Single

func GetSingleRedis() *Single {
	return &singleRedis
}

var key = struct{}{}

func (rConn *Single) BeforeProcess(ctx context.Context, cmd redisv7.Cmder) (context.Context, error) {
	ctxNew := context.WithValue(ctx, key, time.Now())
	return ctxNew, nil
}

func (rConn *Single) AfterProcess(ctx context.Context, cmd redisv7.Cmder) error {
	err := cmd.Err()
	timeStart, ok := ctx.Value(key).(time.Time)
	if ok {
		_metricReqDur.Observe(time.Since(timeStart).Milliseconds(), rConn.addr, cmd.Name())
	}
	if err != nil {
		if errors.Is(err, redisv7.Nil) {
			_metricMisses.Inc(rConn.addr)
		} else {
			_metricReqErr.Inc(rConn.addr, cmd.Name(), err.Error())
		}
		return nil
	}
	_metricHits.Inc(rConn.addr)
	return nil
}

func (rConn *Single) BeforeProcessPipeline(ctx context.Context, cmds []redisv7.Cmder) (context.Context, error) {
	return ctx, nil
}

func (rConn *Single) AfterProcessPipeline(ctx context.Context, cmds []redisv7.Cmder) error {
	return nil
}

func (rConn *Single) Conn(_url string) error {
	rConn.logger = logger.GetLogger().WithFields(logrus.Fields{
		"SourceURL": _url,
	})
	if pUrl, err := url.Parse(_url); err != nil {
		return err
	} else {
		if args, err := url.ParseQuery(pUrl.RawQuery); err == nil { // 解析url对象
			if maxRetry, ok := args["MaxRetries"]; ok {
				mR, _ := strconv.Atoi(maxRetry[0])
				rConn.MaxRetries = mR
			}
			if PoolSize, ok := args["PoolSize"]; ok {
				pS, _ := strconv.Atoi(PoolSize[0])
				rConn.PoolSize = pS
			}
		}
		if len(pUrl.Host) == 3 {
			return errors.New("Can't Parse HostPort Config. ")
		} else {
			rConn.addr = pUrl.Host
			if db := strings.Split(pUrl.Path, "/"); len(db) < 2 {
				return errors.New("Can't Parse DB Config. ")
			} else {
				rConn.DB, _ = strconv.Atoi(db[1])
			}
			rConn.Password, _ = pUrl.User.Password()
		}
	}
	if err := rConn.InitConn(); err != nil {
		return err
	} else {
		rConn.ConnOk = true
		go rConn.goKeepAlive() // 发送心跳保活
		return nil
	}
}

func (rConn *Single) InitConn() error {
	client := redisv7.NewClient(&redisv7.Options{
		Addr:        rConn.addr,       //
		Password:    rConn.Password,   // no password set
		DB:          rConn.DB,         // use default DB
		MaxRetries:  rConn.MaxRetries, //
		PoolSize:    rConn.PoolSize,   //
		DialTimeout: 10 * time.Second,
		PoolTimeout: 1 * time.Second,
		IdleTimeout: 20 * time.Second,
	})
	rConn.Client = client
	rConn.Client.AddHook(rConn)
	if _, err := rConn.Client.Ping().Result(); err != nil {
		rConn.logger.Error(err)
		return err
	} else {
		rConn.logger.Info("Conn To Redis OK.")
		return nil
	}
}

func (rConn *Single) StatusOK() bool {
	return rConn.ConnOk
}

func (rConn *Single) goKeepAlive() {
	for {
		if _, err := rConn.Client.Ping().Result(); err != nil {
			rConn.ConnOk = false
			rConn.logger.Error("Connection Loss, ReBuilding... ")
			if err2 := rConn.InitConn(); err2 != nil {
				rConn.logger.Error("Connection ReBuilding Failed: ", err2)
			} else {
				rConn.ConnOk = true
				rConn.logger.Info("Connection OK, ReBuilding ok ")
			}
			time.Sleep(3 * time.Second)
		}
		time.Sleep(time.Second * 10)
	}
}

func (rConn *Single) KEYS(key string) ([]string, error) {
	if rConn.Client == nil {
		return nil, errors.New("Connection_Loss")
	}
	return rConn.Client.Keys(key).Result()
}

func (rConn *Single) Get(key string) (string, error) {
	if rConn.Client == nil {
		return "", errors.New("Connection_Loss")
	}
	return rConn.Client.Get(key).Result()
}

func (rConn *Single) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	if rConn.Client == nil {
		return "", errors.New("Connection Loss. ")
	}
	return rConn.Client.Set(key, value, expiration).Result()
}

func (rConn *Single) Close() error {
	if rConn.Client == nil {
		return errors.New("Connection Loss. ")
	}
	return rConn.Client.Close()
}

func (rConn *Single) Del(keys ...string) (int64, error) {
	if rConn.Client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return rConn.Client.Del(keys...).Result()
}

func (rConn *Single) Exists(key ...string) (int64, error) {
	if rConn.Client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return rConn.Client.Exists(key...).Result()
}

func (rConn *Single) HGetAll(keys string) (map[string]string, error) {
	if rConn.Client == nil {
		return nil, errors.New("Connection Loss. ")
	}
	return rConn.Client.HGetAll(keys).Result()
}

func (rConn *Single) LLen(key string) (int64, error) {
	if rConn.Client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return rConn.Client.LLen(key).Result()
}

func (rConn *Single) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if rConn.Client == nil {
		return false, errors.New("Connection Loss. ")
	}
	return rConn.Client.SetNX(key, value, expiration).Result()
}

func (rConn *Single) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	if rConn.Client == nil {
		return nil, 0, errors.New("Connection Loss. ")
	}
	return rConn.Client.HScan(key, cursor, match, count).Result()
}

func (rConn *Single) HSet(key string, filed string, value interface{}) (int64, error) {
	if rConn.Client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return rConn.Client.HSet(key, filed, value).Result()
}

func (rConn *Single) HGet(key string, filed string) (string, error) {
	if rConn.Client == nil {
		return "", errors.New("Connection_Loss")
	}
	return rConn.Client.HGet(key, filed).Result()
}

func (rConn *Single) HDel(key string, filed string) (int64, error) {
	if rConn.Client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return rConn.Client.HDel(key, filed).Result()
}

func (rConn *Single) PipeLineHSet(filed string, datas map[string]interface{}) ([]redisv7.Cmder, error) {
	if rConn.Client == nil {
		return nil, errors.New("Connection Loss. ")
	}
	pipeLine := rConn.Client.Pipeline()
	for key, value := range datas {
		pipeLine.HSet(filed, key, value)
	}
	return pipeLine.Exec()
}

func (rConn *Single) GetClient() (interface{}, error) {
	if rConn.Client == nil {
		return nil, errors.New("Connection Loss. ")
	}
	return rConn.Client, nil
}
