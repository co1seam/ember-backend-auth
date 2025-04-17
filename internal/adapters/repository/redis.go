package cache

import "github.com/redis/go-redis"

type Redis struct {
	redis *redis.Client
}

func NewRedis(host, port string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	})
	return &Redis{redis: client}
}
