package redisconn

import (
	"fmt"

	"github.com/go-redis/redis"
)

type RedisPubSubService interface {
	Publish(channel string, message interface{}) error
	Subscribe(channels ...string) ([]*redis.PubSub, error)
	SubscribeWith(callback func(message *redis.Message, err error), close bool, channels ...string) error
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

func (c *RedisPubSubClient) SubscribeWith(callback func(message *redis.Message, err error), close bool, channels ...string) error {
	subs, err := c.Subscribe(channels...)
	if err != nil {
		return err
	}
	go func() {
		for _, sub := range subs {
			go func(ps *redis.PubSub) {
				for {
					msg, err := ps.ReceiveMessage()
					callback(msg, err)
				}
			}(sub)
		}
	}()
	if !close {
		return nil
	}
	err = c.Unsubscribe(channels...)
	if err != nil {
		return err
	}
	err = c.Close()
	return err
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
