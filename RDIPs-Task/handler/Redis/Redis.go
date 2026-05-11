package handler

import (
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

var redisLocker *redislock.Client

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	})

	redisLocker = redislock.New(client)
	return client
}

func GetRedisLockerClient() *redislock.Client {
	return redisLocker
}
