package example

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/sivaosorg/govm/dbx"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/redisx"
	"github.com/sivaosorg/redisconn"
)

func createConn() (*redisconn.Redis, dbx.Dbx) {
	return redisconn.NewClient(*redisx.GetRedisConfigSample().SetDebugMode(true).SetPassword("Tm5@P@ssw0rd"))
}

func TestConn(t *testing.T) {
	_, s := createConn()
	logger.Infof("Redis connection status: %v", s)
}

func TestGetKeys(t *testing.T) {
	r, _ := createConn()
	svc := redisconn.NewRedisService(r.GetConn())
	keys := svc.ListKeys()
	logger.Infof("Current all keys: %v", keys)
}

func TestGetKeysNearExpired(t *testing.T) {
	r, _ := createConn()
	svc := redisconn.NewRedisService(r.GetConn())
	keys := svc.ListKeysNearExpired()
	logger.Infof("Current all keys near expiring: %v", keys)
}

func TestSetKey(t *testing.T) {
	r, _ := createConn()
	svc := redisconn.NewRedisService(r.GetConn())
	err := svc.Set("_key_123", 123, time.Second*10)
	if err != nil {
		logger.Errorf("Setting key on redis got an error", err)
		return
	}
	logger.Successf("Set key on redis successfully")
}

func TestGetKey(t *testing.T) {
	r, _ := createConn()
	svc := redisconn.NewRedisService(r.GetConn())
	value, err := svc.Get("_key_123")
	if err != nil {
		logger.Errorf("Getting key on redis got an error", err)
		return
	}
	logger.Infof("Value: %v", value)
}

func TestRemoveKey(t *testing.T) {
	r, _ := createConn()
	svc := redisconn.NewRedisService(r.GetConn())
	err := svc.Delete("_key_123")
	if err != nil {
		logger.Errorf("Removing key on redis got an error", err)
		return
	}
	logger.Successf("Removed key on redis successfully")
}

func TestPublish(t *testing.T) {
	r, _ := createConn()
	svc := redisconn.NewRedisPubSub(r.GetConn())
	err := svc.Publish("channel_1", 123)
	if err != nil {
		logger.Errorf("Publishing message on redis got an error", err)
		return
	}
	logger.Infof("Published message on redis successfully")
}

func TestConsume(t *testing.T) {
	r, _ := createConn()
	svc := redisconn.NewRedisService(r.GetConn())
	subs, err := svc.SyncPubSub().Subscribe("channel_1")
	if err != nil {
		logger.Errorf("Subscribing message on redis got an error", err)
		return
	}
	go func() {
		for _, sub := range subs {
			go func(ps *redis.PubSub) {
				for {
					msg, err := ps.ReceiveMessage()
					if err != nil {
						logger.Errorf("Subscribing message on channel %v got an error", err, msg.Channel)
						return
					}
					logger.Infof("%s: %s", msg.Channel, msg.Payload)
				}
			}(sub)
		}
	}()

	// Unsubscribe from multiple channels and close the connection
	err = svc.SyncPubSub().Unsubscribe("channel_1")
	if err != nil {
		logger.Errorf("Unsubscribing message on redis got an error", err)
		return
	}
	err = svc.SyncPubSub().Close()
	if err != nil {
		logger.Errorf("Closing message on redis got an error", err)
		return
	}
}

func TestConsumeCallback(t *testing.T) {
	r, _ := createConn()
	svc := redisconn.NewRedisPubSub(r.GetConn())
	callback := func(msg *redis.Message, err error) {
		if err != nil {
			logger.Errorf("Subscribing message on channel %v got an error", err, msg.Channel)
			return
		}
		logger.Infof("%s: %s", msg.Channel, msg.Payload)
	}
	err := svc.SubscribeWith(callback, false, "channel_1")
	if err != nil {
		logger.Errorf("Subscribing message on redis got an error", err)
		return
	}
}
