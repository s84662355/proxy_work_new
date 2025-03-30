package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func RedisIncrByExpirationOnPipe(
	ctx context.Context,
	rdb *redis.Client,
	key string,
	increment int64,
	expiration time.Duration,
) error {
	_, err := rdb.Pipelined(
		ctx,
		func(pipe redis.Pipeliner) error {
			pipe.IncrBy(ctx, key, increment)
			pipe.Expire(ctx, key, expiration)
			return nil
		})
	if err != nil {
		err = fmt.Errorf("RedisIncrByExpirationOnPipe  err:%+v", err)
		return err
	}

	return nil
}
