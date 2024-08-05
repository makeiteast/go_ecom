package cachestore

import (
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var (
	// ErrCacheMiss is the error returned when the requested item is not available in cache
	ErrCacheMiss = errors.New("not found in cache")
	// ErrCacheNotInitialized is the error returned when the cache handler is not initialized
	ErrCacheNotInitialized = errors.New("not initialized")
)

// Config holds all the configuration required for this package
type Config struct {
	Host string
	Port string

	StoreName string
	Username  string
	Password  string

	PoolSize     int
	MinIdleConns int
	PoolTimeout  int
}

// Returns new redis client
func NewRedisClient(cfg *Config) (*redis.Client, error) {
	redisHost := cfg.Host + ":" + cfg.Port
	logrus.Info(redisHost)
	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		Password:     cfg.Password,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
