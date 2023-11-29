package redisconn

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/sivaosorg/govm/dbx"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/redisx"
	"github.com/sivaosorg/govm/utils"
)

var (
	instance *Redis
	_logger  = logger.NewLogger()
)

func NewRedis() *Redis {
	return &Redis{}
}

func (r *Redis) SetConn(value *redis.Client) *Redis {
	r.conn = value
	return r
}

func (r *Redis) SetConfig(value redisx.RedisConfig) *Redis {
	r.Config = value
	return r
}

func (r *Redis) SetState(value dbx.Dbx) *Redis {
	r.State = value
	return r
}

func (r *Redis) Json() string {
	return utils.ToJson(r)
}

func (r *Redis) GetConn() *redis.Client {
	return r.conn
}

func NewClient(config redisx.RedisConfig) (*Redis, dbx.Dbx) {
	s := dbx.NewDbx().SetDatabase(config.Database)
	if !config.IsEnabled {
		s.SetConnected(false).
			SetMessage("Redis unavailable").
			SetError(fmt.Errorf(s.Message))
		instance = NewRedis().SetState(*s)
		return instance, *s
	}
	if instance != nil {
		s.SetConnected(true).SetNewInstance(false)
		instance.SetState(*s)
		return instance, *s
	}
	client := redis.NewClient(&redis.Options{
		Addr:        config.UrlConn,
		Password:    config.Password,
		DialTimeout: config.Timeout,
		DB:          0,
	})
	err := client.Ping().Err()
	if err != nil {
		s.SetConnected(false).SetError(err).SetMessage(err.Error())
		instance = NewRedis().SetState(*s)
		return instance, *s
	}
	if config.DebugMode {
		_logger.Info(fmt.Sprintf("Redis client connection:: %s", config.Json()))
		_logger.Info(fmt.Sprintf("Connected successfully to redis cache:: %s/%s", config.UrlConn, config.Database))
	}
	s.SetConnected(true).SetMessage("Connected successfully").SetPid(os.Getpid()).SetNewInstance(true)
	instance = NewRedis().SetConn(client).SetState(*s)
	return instance, *s
}
