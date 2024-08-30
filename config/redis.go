package config

import (
	"context"
	"fmt"
	"gin_demo/global"
	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		global.Logger.Error(fmt.Sprintln("redis connect err:", err))
		return nil, err
	}
	global.Logger.Info(fmt.Sprintf("redis connect success"))

	return rdb, nil
}
