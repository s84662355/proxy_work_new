package common

import (
	"sync"

	"github.com/redis/go-redis/v9"
	"mproxy/config"
)

var GetRedisDB = sync.OnceValue[*redis.Client](func() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.ConfData.Redis.Addr,
		Password:     config.ConfData.Redis.Password,
		DB:           config.ConfData.Redis.DB,
		MinIdleConns: config.ConfData.Redis.MinIdleConns,
		MaxIdleConns: config.ConfData.Redis.MaxIdleConns,
		PoolSize:     config.ConfData.Redis.PoolSize,
	})
	return rdb
})
