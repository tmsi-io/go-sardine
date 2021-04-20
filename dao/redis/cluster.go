package redis

import (
	"context"
	"errors"
	_redis "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/tmsi-io/go-sardine/logger"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Cluster struct {
	Addr         []string
	Host         string
	Password     string
	DB           int
	MaxRetries   int
	PoolSize     int
	Client       *_redis.ClusterClient
	ConnOk       bool
	ReadOnly     bool
	Logger       *logrus.Entry
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var clusterkey = struct{}{}

func (cluster *Cluster) BeforeProcess(ctx context.Context, cmd _redis.Cmder) (context.Context, error) {
	ctxNew := context.WithValue(ctx, clusterkey, time.Now())
	return ctxNew, nil
}

func (cluster *Cluster) AfterProcess(ctx context.Context, cmd _redis.Cmder) error {
	err := cmd.Err()
	timeStart, ok := ctx.Value(clusterkey).(time.Time)
	if ok {
		_metricReqDur.Observe(time.Since(timeStart).Milliseconds(), cluster.Host, cmd.Name())
	}
	if err != nil {
		if errors.Is(err, _redis.Nil) {
			_metricMisses.Inc(cluster.Host)
		} else {
			_metricReqErr.Inc(cluster.Host, cmd.Name(), err.Error())
		}
		return nil
	}
	_metricHits.Inc(cluster.Host)
	return nil
}

func (cluster *Cluster) BeforeProcessPipeline(ctx context.Context, cmds []_redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (cluster *Cluster) AfterProcessPipeline(ctx context.Context, cmds []_redis.Cmder) error {
	return nil
}

func (cluster *Cluster) ExtractURL(_url string) error {
	if pUrl, err := url.Parse(_url); err != nil {
		return err
	} else {
		if args, err := url.ParseQuery(pUrl.RawQuery); err == nil {
			if maxRetry, ok := args["MaxRetries"]; ok {
				mR, _ := strconv.Atoi(maxRetry[0])
				cluster.MaxRetries = mR
			}
			if PoolSize, ok := args["PoolSize"]; ok {
				pS, _ := strconv.Atoi(PoolSize[0])
				cluster.PoolSize = pS
			}
		}
		if len(pUrl.Host) == 3 {
			return errors.New("Can't Parse HostPort Config. ")
		} else {
			cluster.Host = pUrl.Host
			cluster.Addr = strings.Split(pUrl.Host, ",")
			if db := strings.Split(pUrl.Path, "/"); len(db) < 2 {
				return errors.New("Can't Parse DB Config. ")
			} else {
				cluster.DB, _ = strconv.Atoi(db[1])
			}
			cluster.Password, _ = pUrl.User.Password()
		}
	}
	return nil
}

func (cluster *Cluster) Conn(_url string) error {
	cluster.Logger = logger.GetLogger().WithFields(logrus.Fields{"URL": _url})
	cluster.MaxRetries = 3
	cluster.DialTimeout = 3 * time.Second
	cluster.ReadTimeout = time.Second
	cluster.WriteTimeout = time.Second
	if err := cluster.ExtractURL(_url); err != nil {
		return err
	}
	cluster.Client = _redis.NewClusterClient(&_redis.ClusterOptions{
		Addrs:      cluster.Addr,     //set redis cluster url
		Password:   cluster.Password, //set password
		PoolSize:   cluster.PoolSize,
		MaxRetries: cluster.MaxRetries,
	})
	cluster.Client.AddHook(cluster)
	if _, err := cluster.Client.Ping().Result(); err != nil {
		return err
	} else {
		cluster.ConnOk = true
		go cluster.goKeepAlive()
		return nil
	}
}

func (cluster *Cluster) goKeepAlive() {
	for {
		if _, err := cluster.Client.Ping().Result(); err != nil {
			cluster.ConnOk = false
			cluster.Logger.Error("Connection Loss, ReBuilding... ")
			if err2 := cluster.InitConn(); err2 != nil {
				cluster.Logger.Error("Connection ReBuilding Failed: ", err2)
			} else {
				cluster.ConnOk = true
				cluster.Logger.Info("Connection OK, ReBuilding ok ")
			}
		}
		time.Sleep(time.Second * 3)
	}
}

func (cluster *Cluster) InitConn() error {
	cluster.Client = _redis.NewClusterClient(&_redis.ClusterOptions{
		Addrs:       cluster.Addr,       //
		Password:    cluster.Password,   // no password set
		MaxRetries:  cluster.MaxRetries, //
		PoolSize:    cluster.PoolSize,   //
		DialTimeout: 10 * time.Second,
		PoolTimeout: 1 * time.Second,
		IdleTimeout: 20 * time.Second,
	})
	if _, err := cluster.Client.Ping().Result(); err != nil {
		return err
	} else {
		return nil
	}
}

func (cluster *Cluster) StatusOK() bool {
	return cluster.ConnOk
}

func (cluster *Cluster) KEYS(key string) ([]string, error) {
	if cluster.Client == nil {
		return nil, errors.New("Connection_Loss")
	}
	return cluster.Client.Keys(key).Result()
}

func (cluster *Cluster) Get(key string) (string, error) {
	if cluster.Client == nil {
		return "", errors.New("Connection_Loss")
	}
	return cluster.Client.Get(key).Result()
}

func (cluster *Cluster) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	if cluster.Client == nil {
		return "", errors.New("Connection Loss. ")
	}
	return cluster.Client.Set(key, value, expiration).Result()
}

func (cluster *Cluster) Close() error {
	if cluster.Client == nil {
		return errors.New("Connection Loss. ")
	}
	return cluster.Client.Close()
}

func (cluster *Cluster) Del(keys ...string) (int64, error) {
	if cluster.Client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return cluster.Client.Del(keys...).Result()
}

func (cluster *Cluster) Exists(key ...string) (int64, error) {
	if cluster.Client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return cluster.Client.Exists(key...).Result()
}

func (cluster *Cluster) HGetAll(keys string) (map[string]string, error) {
	if cluster.Client == nil {
		return nil, errors.New("Connection Loss. ")
	}
	return cluster.Client.HGetAll(keys).Result()
}

func (cluster *Cluster) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if cluster.Client == nil {
		return false, errors.New("Connection Loss. ")
	}
	return cluster.Client.SetNX(key, value, expiration).Result()
}
func (cluster *Cluster) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	if cluster.Client == nil {
		return nil, 0, errors.New("Connection Loss. ")
	}
	return cluster.Client.HScan(key, cursor, match, count).Result()
}
func (cluster *Cluster) HSet(key string, filed string, value interface{}) (int64, error) {
	if cluster.Client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return cluster.Client.HSet(key, filed, value).Result()
}

func (cluster *Cluster) HGet(key string, filed string) (string, error) {
	if cluster.Client == nil {
		return "", errors.New("Connection_Loss")
	}
	return cluster.Client.HGet(key, filed).Result()
}

func (cluster *Cluster) HDel(key string, filed string) (int64, error) {
	if cluster.Client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return cluster.Client.HDel(key, filed).Result()
}

func (cluster *Cluster) PipeLineHSet(filed string, datas map[string]interface{}) ([]_redis.Cmder, error) {
	if cluster.Client == nil {
		return nil, errors.New("Connection Loss. ")
	}
	pipeLine := cluster.Client.Pipeline()
	for key, value := range datas {
		pipeLine.HSet(filed, key, value)
	}
	return pipeLine.Exec()
}

func (cluster *Cluster) GetClient() (interface{}, error) {
	if cluster.Client == nil {
		return nil, errors.New("Connection Loss. ")
	}
	return cluster.Client, nil
}

func (cluster *Cluster) LLen(key string) (int64, error) {
	if cluster.Client == nil {
		return 0, errors.New("Connection Loss. ")
	}
	return cluster.Client.LLen(key).Result()
}
