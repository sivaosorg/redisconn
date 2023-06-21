package redisconn

import (
	"fmt"

	"github.com/go-redis/redis"
)

type RedisPubSubService interface {
	Publish(channel string, message interface{}) error
	Subscribe(channels ...string) ([]*redis.PubSub, error)
	Unsubscribe(channels ...string) error
	Close() error
}

func NewRedisPubSub(redisConn *redis.Client) RedisPubSubService {
	s := &RedisPubSubClient{
		redisConn:       redisConn,
		subscriptionMap: make(map[string]*redis.PubSub),
	}
	return s
}

func (c *RedisPubSubClient) Publish(channel string, message interface{}) error {
	return c.redisConn.Publish(channel, message).Err()
}

func (c *RedisPubSubClient) Subscribe(channels ...string) ([]*redis.PubSub, error) {
	pubSubs := []*redis.PubSub{}
	for _, channel := range channels {
		pubsub := c.redisConn.Subscribe(channel)
		c.subscriptionMap[channel] = pubsub
		pubSubs = append(pubSubs, pubsub)
	}
	return pubSubs, nil
}

func (c *RedisPubSubClient) Unsubscribe(channels ...string) error {
	for _, channel := range channels {
		pubsub, ok := c.subscriptionMap[channel]
		if !ok {
			return fmt.Errorf("not subscribed to channel %s", channel)
		}
		err := pubsub.Unsubscribe(channel)
		if err != nil {
			return err
		}
		delete(c.subscriptionMap, channel)
	}
	return nil
}

func (c *RedisPubSubClient) Close() error {
	for channel, pubsub := range c.subscriptionMap {
		err := pubsub.Unsubscribe(channel)
		if err != nil {
			return err
		}
		delete(c.subscriptionMap, channel)
	}
	return c.redisConn.Close()
}
