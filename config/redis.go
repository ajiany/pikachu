package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

var redisCli redis.UniversalClient
var redisOnce sync.Once

// Redis redis服务初始化
func Redis() redis.UniversalClient {
	if redisCli != nil {
		return redisCli
	}

	redisOnce.Do(func() {
		redisCli = redistrace.NewClient(&redis.Options{
			Network:      "tcp",
			Addr:         fmt.Sprintf("%s:%d", Cfg.RedisCacheHost, Cfg.RedisCachePort),
			PoolSize:     50,
			MinIdleConns: 50,
		})
	})

	if err := redisCli.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return redisCli
}
