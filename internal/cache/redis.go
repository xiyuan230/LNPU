package cache

import (
	"MyLNPU/conf"
	"MyLNPU/internal/logger"
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var rdb *redis.Client

func Init() {
	redisConf := conf.GetConfig().Redis
	addr := redisConf.Host + ":" + redisConf.Port
	password := redisConf.Password
	DB := redisConf.DB
	ctx := context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       DB,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Errorf("redis初始化出错... ERROR: %s", err)
		os.Exit(-1)
	}
	logger.Println("redis初始化成功...")
}

func GetRDB() *redis.Client {
	return rdb
}

func Set(key string, value any, expiration time.Duration) error {
	ctx := context.Background()
	return rdb.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	ctx := context.Background()
	return rdb.Get(ctx, key).Result()
}

func Del(key string) error {
	ctx := context.Background()
	return rdb.Del(ctx, key).Err()
}

func HSet(key string, value ...any) error {
	ctx := context.Background()
	return rdb.HSet(ctx, key, value).Err()
}

func HGet(key, field string) (string, error) {
	ctx := context.Background()
	return rdb.HGet(ctx, key, field).Result()
}

func HMGet(key string, fields ...string) ([]any, error) {
	ctx := context.Background()
	return rdb.HMGet(ctx, key, fields...).Result()
}
