package redis

import (
	"context"
	"errors"
	_redis "github.com/go-redis/redis/v7"
	"time"
)

func (s *Sentinel) BeforeProcess(ctx context.Context, cmd _redis.Cmder) (context.Context, error) {
	ctxNew := context.WithValue(ctx, clusterKey, time.Now())
	return ctxNew, nil
}

func (s *Sentinel) AfterProcess(ctx context.Context, cmd _redis.Cmder) error {
	err := cmd.Err()
	timeStart, ok := ctx.Value(clusterKey).(time.Time)
	if ok {
		_metricReqDur.Observe(time.Since(timeStart).Milliseconds(), s.host, cmd.Name())
	}
	if err != nil {
		if errors.Is(err, _redis.Nil) {
			_metricMisses.Inc(s.host)
		} else {
			_metricReqErr.Inc(s.host, cmd.Name(), err.Error())
		}
	} else {
		_metricHits.Inc(s.host)
	}
	return nil
}

func (s *Sentinel) BeforeProcessPipeline(ctx context.Context, cmds []_redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (s *Sentinel) AfterProcessPipeline(ctx context.Context, cmds []_redis.Cmder) error {
	return nil
}
