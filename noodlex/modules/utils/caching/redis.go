
package caching

import (
	"time"

	"github.com/NoodleSoup/NoodleX/noodlex"
	"github.com/go-redis/redis"
)

var REDIS *redis.Client

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:         noodlex.BotConfig.RedisAddress,
		Password:     noodlex.BotConfig.RedisPassword,
		DB:           0,
		DialTimeout:  time.Second,
		MinIdleConns: 0,
	})
	REDIS = client
}
