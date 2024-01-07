package redisconn

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/sivaosorg/govm/dbx"
	"github.com/sivaosorg/govm/redisx"
)

type RedisPubSubClient struct {
	redisConn       *redis.Client
	subscriptionMap map[string]*redis.PubSub
}

type RedisMutex struct {
	mutexes map[string]*sync.Mutex
}

type Redis struct {
	conn   *redis.Client      `json:"-"`
	Config redisx.RedisConfig `json:"config,omitempty"`
	State  dbx.Dbx            `json:"state,omitempty"`
}
