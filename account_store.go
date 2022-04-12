package aegirdungeons

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ErrUserNotRegistered = errors.New("accounts: user not registered")

// InMemoryStore ...
type InMemoryStore struct {
	data map[string]string
}

// NewInMemoryStore ...
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: map[string]string{}}
}

// Set ...
func (ims *InMemoryStore) Set(_ context.Context, key string, value string) error {
	ims.data[key] = value

	return nil
}

// Get ...
func (ims *InMemoryStore) Get(_ context.Context, key string) (string, error) {
	v, ok := ims.data[key]
	if !ok {
		return "", ErrUserNotRegistered
	}

	return v, nil
}

// Close ...
func (ims *InMemoryStore) Close() error {
	return nil
}

// RedisStore ...
type RedisStore struct {
	rdb *redis.Client
}

// NewRedisStore ...
func NewRedisStore(url string) (*RedisStore, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("redis: parse_url: %v", err)
	}

	if opts.TLSConfig != nil {
		opts.TLSConfig.InsecureSkipVerify = true
	}

	rdb := redis.NewClient(opts)
	if _, err := rdb.Ping(context.TODO()).Result(); err != nil {
		return nil, fmt.Errorf("redis: ping: %v", err)
	}

	return &RedisStore{rdb: rdb}, nil
}

// Set ...
func (r *RedisStore) Set(ctx context.Context, key string, value string) error {
	return r.rdb.Set(ctx, key, value, 0).Err()
}

// Get ...
func (r *RedisStore) Get(ctx context.Context, key string) (string, error) {
	v, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrUserNotRegistered
		}
		return "", err
	}

	return v, nil
}

// Close ...
func (r *RedisStore) Close() error {
	return r.rdb.Close()
}
