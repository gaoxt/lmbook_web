package singletonRedis

import (
	"sync"

	"github.com/go-redis/redis"
)

type Redis struct {
	*redis.Client
}

var redisInstance *Redis
var once sync.Once

func GetRedis() *Redis {
	once.Do(func() {
		redisInstance = &Redis{redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		})}
	})
	return redisInstance
}
