package redis

import (
	"github.com/go-redis/redis"
	"github.com/developer-guy/dockeroptic/environments"
)

func Connect() *redis.Client {
	options := &redis.Options{
		Addr:     environments.RedisConnectionBaseUri,
		Password: "",
		DB:       0,
	}
	client := redis.NewClient(options)
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client
}
