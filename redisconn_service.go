package redisconn

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/utils"
)

type RedisService interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
	ListKeys() map[string]string
	ListKeysNearExpired() []string
	Increase(key string) error
	Decrease(key string) error
	Lock(key string, value interface{}, expiration time.Duration) error
	Unlock(key string, value interface{}) error
	Handler() *redis.Client

	// sync pubsub hook
	SyncPubSub() RedisPubSubService
}

type redisServiceImpl struct {
	redisConn *redis.Client
	mutex     *RedisMutex
	pubsub    RedisPubSubService
}

func newRedisMutex() *RedisMutex {
	return &RedisMutex{make(map[string]*sync.Mutex)}
}

func NewRedisService(redisConn *redis.Client) RedisService {
	s := &redisServiceImpl{
		redisConn: redisConn,
		mutex:     newRedisMutex(),
		pubsub:    NewRedisPubSub(redisConn),
	}
	return s
}

func (r *RedisMutex) getMutex(key string) *sync.Mutex {
	mutex, ok := r.mutexes[key]
	if !ok {
		mutex = &sync.Mutex{}
		r.mutexes[key] = mutex
	}
	return mutex
}

func (r *redisServiceImpl) Get(key string) (interface{}, error) {
	if utils.IsEmpty(key) {
		return nil, fmt.Errorf("invalid key")
	}
	r.mutex.getMutex(key).Lock()
	defer r.mutex.getMutex(key).Unlock()
	// Get the type of the key
	keyType, err := r.redisConn.Type(key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s not found", key)
	}
	if err != nil {
		return nil, err
	}
	// Check if the key is a hash
	if keyType == "hash" {
		hash, err := r.redisConn.HGetAll(key).Result()
		if err != nil {
			return nil, err
		}
		return hash, nil
	}
	// Otherwise, get the value of the key
	value, err := r.redisConn.Get(key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s not found", key)
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (r *redisServiceImpl) Set(key string, value interface{}, expiration time.Duration) error {
	if utils.IsEmpty(key) {
		return fmt.Errorf("invalid key")
	}
	json, err := utils.MarshalToString(value)
	if err != nil {
		return err
	}
	r.mutex.getMutex(key).Lock()
	defer r.mutex.getMutex(key).Unlock()
	return r.redisConn.Set(key, json, expiration).Err()
}

func (r *redisServiceImpl) Delete(key string) error {
	if utils.IsEmpty(key) {
		return fmt.Errorf("invalid key")
	}
	r.mutex.getMutex(key).Lock()
	defer r.mutex.getMutex(key).Unlock()
	return r.redisConn.Del(key).Err()
}

func (r *redisServiceImpl) ListKeys() map[string]string {
	r.mutex.getMutex("*").Lock()
	defer r.mutex.getMutex("*").Unlock()
	keys := make(map[string]string)
	var cursor uint64
	for {
		var batchKeys []string
		var err error
		batchKeys, cursor, err = r.redisConn.Scan(cursor, "*", 10).Result()
		if err != nil {
			logger.Errorf("ListKeys has an error occurred: %s", err, err.Error())
			break
		}
		for _, key := range batchKeys {
			keyType, err := r.redisConn.Type(key).Result()
			if err != nil {
				logger.Errorf("ListKeys has an error occurred: %s", err, err.Error())
				break
			}
			keys[key] = keyType
		}
		if cursor == 0 {
			break
		}
	}
	return keys
}

func (r *redisServiceImpl) ListKeysNearExpired() []string {
	r.mutex.getMutex("*").Lock()
	defer r.mutex.getMutex("*").Unlock()
	keys := []string{}
	var cursor uint64
	for {
		var batchKeys []string
		var err error
		batchKeys, cursor, err = r.redisConn.Scan(cursor, "*", 10).Result()
		if err != nil {
			logger.Errorf("ListKeysNearExpired has an error occurred: %s", err, err.Error())
			break
		}
		for _, key := range batchKeys {
			ttl, err := r.redisConn.TTL(key).Result()
			if err != nil {
				logger.Errorf("ListKeysNearExpired has an error occurred: %s", err, err.Error())
				break
			}
			// If TTL is negative, ignore the key
			if ttl < 0 {
				continue
			}
			// If TTL is zero or less than a minute, add to the expired keys list
			if ttl == 0 || ttl.Minutes() < 1 {
				keys = append(keys, key)
			}
		}
		// Exit the loop once cursor gets to zero
		if cursor == 0 {
			break
		}
	}
	return keys
}

func (r *redisServiceImpl) Increase(key string) error {
	r.mutex.getMutex(key).Lock()
	defer r.mutex.getMutex(key).Unlock()
	_, err := r.redisConn.Incr(key).Result()
	return err
}

func (r *redisServiceImpl) Decrease(key string) error {
	r.mutex.getMutex(key).Lock()
	defer r.mutex.getMutex(key).Unlock()
	_, err := r.redisConn.Decr(key).Result()
	return err
}

func (r *redisServiceImpl) Lock(key string, value interface{}, expiration time.Duration) error {
	r.mutex.getMutex(key).Lock()
	defer r.mutex.getMutex(key).Unlock()
	v := utils.ToJson(value)
	return r.redisConn.SetNX(key, v, expiration).Err()
}

func (r *redisServiceImpl) Unlock(key string, value interface{}) error {
	r.mutex.getMutex(key).Lock()
	defer r.mutex.getMutex(key).Unlock()

	// Get the current value of the key
	currentValue, err := r.redisConn.Get(key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key %s not found", key)
	}
	if err != nil {
		return err
	}
	// Check if the key is locked by the given value
	if currentValue != utils.ToJson(value) {
		return fmt.Errorf("key %s is not locked by the given value", key)
	}
	// Delete the key from redis
	return r.redisConn.Del(key).Err()
}

func (r *redisServiceImpl) Handler() *redis.Client {
	return r.redisConn
}

func (r *redisServiceImpl) SyncPubSub() RedisPubSubService {
	return r.pubsub
}
