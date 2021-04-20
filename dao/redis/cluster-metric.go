package redis

import (
	"context"
	"errors"
	_redis "github.com/go-redis/redis/v7"
	"time"
)

func (cluster *Cluster) BeforeProcess(ctx context.Context, cmd _redis.Cmder) (context.Context, error) {
	ctxNew := context.WithValue(ctx, clusterKey, time.Now())
	return ctxNew, nil
}

func (cluster *Cluster) AfterProcess(ctx context.Context, cmd _redis.Cmder) error {
	err := cmd.Err()
	timeStart, ok := ctx.Value(clusterKey).(time.Time)
	if ok {
		_metricReqDur.Observe(time.Since(timeStart).Milliseconds(), cluster.host, cmd.Name())
	}
	if err != nil {
		if errors.Is(err, _redis.Nil) {
			_metricMisses.Inc(cluster.host)
		} else {
			_metricReqErr.Inc(cluster.host, cmd.Name(), err.Error())
		}
		return nil
	}
	_metricHits.Inc(cluster.host)
	return nil
}

func (cluster *Cluster) BeforeProcessPipeline(ctx context.Context, cmds []_redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (cluster *Cluster) AfterProcessPipeline(ctx context.Context, cmds []_redis.Cmder) error {
	return nil
}
