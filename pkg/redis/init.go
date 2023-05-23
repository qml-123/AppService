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
		Addr:     "114.116.15.130:6379",
		Password: "123456",
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

func HSetTask(key string, value interface{}) (bool, error) {
	return _client.HSet(key, "args", value).Result()
}

func ZAddTask(score float64, member string) (int64, error) {
	return _client.ZAdd("task", redis.Z{Score: score, Member: member}).Result()
}

func HGetTask(key string) (string, error) {
	return _client.HGet(key, "args").Result()
}

func HDelTask(key string) (int64, error) {
	return _client.HDel(key, "args").Result()
}

func ZPopMinTask() ([]redis.Z, error) {
	return _client.ZPopMin("task", 1).Result()
}

func SAdd(key string, value ...string) (int64, error) {
	return _client.SAdd(key, value).Result()
}

func SMember(key string) ([]string, error) {
	return _client.SMembers(key).Result()
}

func SRem(key string, value ...string) (int64, error) {
	return _client.SRem(key, value).Result()
}
