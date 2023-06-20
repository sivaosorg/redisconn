package redisconn

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sivaosorg/govm/dbx"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/redisx"
)

var (
	instance *redis.Client
	_logger  = logger.NewLogger()
)

func NewClient(config redisx.RedisConfig) (*redis.Client, dbx.Dbx) {
	s := dbx.NewDbx().SetDatabase(config.Database)
	if !config.IsEnabled {
		s.SetConnected(false).
			SetMessage("Redis unavailable").
			SetError(fmt.Errorf(s.Message))
		return &redis.Client{}, *s
	}
	if instance != nil {
		s.SetConnected(true)
		return instance, *s
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.UrlConn,
		Password: config.Password,
		DB:       0,
	})
	err := client.Ping().Err()
	if err != nil {
		s.SetConnected(false).SetError(err).SetMessage(err.Error())
		return &redis.Client{}, *s
	}
	if config.DebugMode {
		_logger.Info(fmt.Sprintf("Redis client connection:: %s", config.Json()))
		_logger.Info(fmt.Sprintf("Connected successfully to redis cache:: %s/%s", config.UrlConn, config.Database))
	}
	instance = client
	s.SetConnected(true).SetMessage("Connection established")
	return instance, *s
}
