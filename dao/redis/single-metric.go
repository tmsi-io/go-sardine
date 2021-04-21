package redis

import (
	"context"
	"errors"
	redisv7 "github.com/go-redis/redis/v7"
	"time"
)

func (single *Single) BeforeProcess(ctx context.Context, cmd redisv7.Cmder) (context.Context, error) {
	ctxNew := context.WithValue(ctx, key, time.Now())
	return ctxNew, nil
}

func (single *Single) AfterProcess(ctx context.Context, cmd redisv7.Cmder) error {
	err := cmd.Err()
	tStart, ok := ctx.Value(key).(time.Time)
	if ok {
		_metricReqDur.Observe(time.Since(tStart).Milliseconds(), single.addr, cmd.Name())
	}
	if err != nil {
		if errors.Is(err, redisv7.Nil) {
			_metricMisses.Inc(single.addr)
		} else {
			_metricReqErr.Inc(single.addr, cmd.Name(), err.Error())
		}
		return nil
	}
	_metricHits.Inc(single.addr)
	return nil
}

func (single *Single) BeforeProcessPipeline(ctx context.Context, cmds []redisv7.Cmder) (context.Context, error) {
	return ctx, nil
}

func (single *Single) AfterProcessPipeline(ctx context.Context, cmds []redisv7.Cmder) error {
	return nil
}
