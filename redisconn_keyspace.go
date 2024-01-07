package redisconn

import (
	"sync"

	"github.com/go-redis/redis"
)

type RedisKeySpaceService interface {
	AddCallback(c func(message *redis.Message, err error)) RedisKeySpaceService
	Register()
	Unregister()
	PatternKeySpace() string
	PatternKeyEvent() string
}

type redisKeySpaceServiceImpl struct {
	client   *redis.Client
	callback []func(message *redis.Message, err error)
	stop     chan struct{}
	mutex    sync.Mutex
}

func NewRedisKeySpaceService(client *redis.Client) RedisKeySpaceService {
	return &redisKeySpaceServiceImpl{
		client:   client,
		callback: make([]func(event *redis.Message, err error), 0),
		stop:     make(chan struct{}),
	}
}

func (s *redisKeySpaceServiceImpl) AddCallback(c func(message *redis.Message, err error)) RedisKeySpaceService {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.callback = append(s.callback, c)
	return s
}

func (s *redisKeySpaceServiceImpl) Register() {
	go s.run()
}

func (s *redisKeySpaceServiceImpl) Unregister() {
	close(s.stop)
}

func (s *redisKeySpaceServiceImpl) PatternKeySpace() string {
	return "__keyspace@*__:*"
}

func (s *redisKeySpaceServiceImpl) PatternKeyEvent() string {
	return "__keyevent@*__:*"
}

func (s *redisKeySpaceServiceImpl) run() {
	pattern := []string{s.PatternKeySpace(), s.PatternKeyEvent()}
	pubsub := s.client.PSubscribe(pattern...)
	for {
		select {
		case <-s.stop:
			_ = pubsub.Close()
			return
		default:
			msg, err := pubsub.ReceiveMessage()
			s.hook(msg, err)
		}
	}
}

func (s *redisKeySpaceServiceImpl) hook(message *redis.Message, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, callback := range s.callback {
		callback(message, err)
	}
}
