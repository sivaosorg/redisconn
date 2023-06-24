package redisconn

import (
	"sync"

	"github.com/go-redis/redis"
)

type RedisPubSubClient struct {
	redisConn       *redis.Client            `json:"-"`
	subscriptionMap map[string]*redis.PubSub `json:"-"`
}

type RedisMutex struct {
	mutexes map[string]*sync.Mutex
}
