package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host   string `toml:"host"`
	Port   string `toml:"port"`
	Pwd    string `toml:"pwd"`
	Prefix string `toml:"rediskey"`
}

type RedisClient struct {
	Client *redis.Client
	Conf   RedisConfig
}

var rdb RedisClient

func InitCache(redisConfig RedisConfig) error {
	rdb = RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
			Password: redisConfig.Pwd, // no password set
			DB:       0,               // use default DB
		}),
		Conf: redisConfig,
	}
	return nil
}

//增加 redis 前缀
func addPrefix(key string) string {
	return fmt.Sprintf("%s:%s", rdb.Conf.Prefix, key)
}

//expiration = 0表示没有过期时间
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := rdb.Client.Set(ctx, addPrefix(key), value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get(ctx context.Context, key string) (string, error) {
	key = addPrefix(key)
	res, err := rdb.Client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			return "", err
		}
		return "", errors.New("key " + key + "not exist")
	}
	return res, nil
}

func GetV2(ctx context.Context, key string) (string, error) {
	res, err := rdb.Client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			return "", err
		}
		return "", errors.New("key " + key + "not exist")
	}
	return res, nil
}

func GetBytes(ctx context.Context, key string) ([]byte, error) {
	res, err := rdb.Client.Get(ctx, addPrefix(key)).Bytes()
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
		return nil, errors.New("key " + key + "not exist")
	}
	return res, nil
}

func SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	res, err := rdb.Client.SetEX(ctx, addPrefix(key), value, expiration).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	ok, err := rdb.Client.SetNX(ctx, addPrefix(key), value, expiration).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func Del(ctx context.Context, key string) (int64, error) {
	i, err := rdb.Client.Del(ctx, addPrefix(key)).Result()
	if err != nil {
		return 0, err
	}
	return i, err
}

func PTTL(ctx context.Context, key string) (time.Duration, error) {
	res, err := rdb.Client.PTTL(ctx, addPrefix(key)).Result()
	if err != nil {
		if err != redis.Nil {
			return 0, err
		}
		return 0, errors.New("key " + key + "not exist")
	}
	return res, nil
}

func Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	ok, err := rdb.Client.Expire(ctx, addPrefix(key), expiration).Result()
	if err != nil {
		if err != redis.Nil {
			return false, err
		}
		return false, errors.New("key " + key + "not exist")
	}
	return ok, nil
}

// HSet accepts values in following formats:
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
func HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	i, err := rdb.Client.HSet(ctx, addPrefix(key), values).Result()
	if err != nil {
		return 0, err
	}
	return i, err
}

func HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	i, err := rdb.Client.HDel(ctx, addPrefix(key), fields...).Result()
	if err != nil {
		return 0, err
	}
	return i, err
}

func HGet(ctx context.Context, key, field string) (string, error) {
	res, err := rdb.Client.HGet(ctx, addPrefix(key), field).Result()
	if err != nil {
		if err != redis.Nil {
			return "", err
		}
		return "", errors.New("key " + key + "not exist")
	}
	return res, err
}

func HGetBytes(ctx context.Context, key, field string) ([]byte, error) {
	res, err := rdb.Client.HGet(ctx, addPrefix(key), field).Bytes()
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
		return nil, errors.New("key " + key + "not exist")
	}
	return res, err
}

func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	res, err := rdb.Client.HGetAll(ctx, addPrefix(key)).Result()
	if err != nil {
		if err != redis.Nil {
			return res, err
		}
		return res, errors.New("key " + key + "not exist")
	}
	return res, err
}

//expiration = 0表示没有过期时间
func SetData(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	b, _ := json.Marshal(value)
	err := rdb.Client.Set(ctx, addPrefix(key), string(b), expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get 从redis中读取指定值，使用json的反序列化方式
func GetData(ctx context.Context, key string, value interface{}) error {
	bytes, err := rdb.Client.Get(ctx, addPrefix(key)).Bytes()
	if err != nil {
		if err != redis.Nil {
			return err
		}
		return nil
	}
	err = jsoniter.Unmarshal(bytes, &value)
	if err != nil {
		return err
	}
	return nil
}
