package redis

import (
	"time"

	"github.com/go-redis/redis"
)

var (
	_client *redis.Client
)

func InitRedis() error {
	_client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // 如果有设置密码，需要填写密码
		DB:       0,
	})
	return nil
}

func Set(key string, value interface{}, expire time.Duration) error {
	return _client.Set(key, value, expire).Err()
}

func SetNX(key string, value interface{}, expire time.Duration) (bool, error) {
	return _client.SetNX(key, value, expire).Result()
}

func Del(key string) (int64, error) {
	return _client.Del(key).Result()
}

func Get(key string) (string, error) {
	return _client.Get(key).Result()
}

func Incr(key string) (int64, error) {
	return _client.Incr(key).Result()
}
