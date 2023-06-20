package redisconn

import "github.com/go-redis/redis"

type RedisPubSubClient struct {
	redisConn       *redis.Client            `json:"-"`
	subscriptionMap map[string]*redis.PubSub `json:"-"`
}
