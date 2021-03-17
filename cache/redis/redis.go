package redis

import (
	"time"

	r "github.com/go-redis/redis"
)

// RedisStorage is an option for caching
type RedisStorage struct {
	client *r.Client
}

// NewStorage sets a new instance of redis
func NewStorage(url string) (*RedisStorage, error) {
	var opts *r.Options
	var err error

	if opts, err = r.ParseURL(url); err != nil {
		return nil, err
	}

	return &RedisStorage{
		client: r.NewClient(opts),
	}, nil
}

// Get implements getting something from cache
func (s RedisStorage) Get(key string) []byte {
	val, _ := s.client.Get(key).Bytes()
	return val
}

// Set implements setting something to cache
func (s *RedisStorage) Set(key string, content []byte, duration time.Duration) {
	s.client.Set(key, content, duration)
}
