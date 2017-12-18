package configs

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/go-redis/redis"
)

type API struct {
	DockerClient *docker.Client
	RedisClient  *redis.Client
}
