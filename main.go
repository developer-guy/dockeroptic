package main

import (
	"github.com/gin-gonic/gin"
	"github.com/developer-guy/dockeroptic/configs"
	"github.com/developer-guy/dockeroptic/redis"
	localContainers "github.com/developer-guy/dockeroptic/containers/local"
	registryContainers "github.com/developer-guy/dockeroptic/containers/registry"
	localImages "github.com/developer-guy/dockeroptic/images/local"
	"github.com/fsouza/go-dockerclient"
	"github.com/googollee/go-socket.io"
	"time"
	"github.com/itsjamie/gin-cors"
	"github.com/developer-guy/dockeroptic/socket"
)

func main() {

	server, errz := socketio.NewServer(nil)
	if errz != nil {
		panic(errz)
	}

	var (
		redisClient       = redis.Connect()
		router            = gin.Default()
		dockerClient, err = docker.NewClientFromEnv()
	)


	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))


	if err != nil {
		panic(err)
	}

	api := configs.API{DockerClient: dockerClient, RedisClient: redisClient}

	registryGroup := router.Group("/registry", registryContainers.AuthRequired(api))
	registryGroup.GET("/containers", registryContainers.GetRepositories(api))

	localGroup := router.Group("/local")
	localGroup.GET("/images", localImages.GetAll(api))
	localGroup.GET("/images/:id/detail", localImages.GetImageDetail(api))
	localGroup.GET("/images/:id/history", localImages.GetImageHistory(api))
	localGroup.GET("/containers", localContainers.GetAllOrFilterByState(api))
	localGroup.GET("/containers/:id/detail", localContainers.GetContainerDetail(api))
	localGroup.GET("/containers/:id/action/:action", localContainers.DoProcessOnContainer(api))

	router.GET("/socket.io/", socket.HandleSocketConnection(server, api))
	router.POST("/socket.io/", socket.HandleSocketConnection(server, api))

	router.Run()
}
