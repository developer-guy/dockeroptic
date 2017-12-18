package local

import (
	"github.com/gin-gonic/gin"
	"github.com/fsouza/go-dockerclient"
	"github.com/developer-guy/dockeroptic/configs"
	"github.com/developer-guy/dockeroptic/utils"
)

func GetAllOrFilterByState(api configs.API) gin.HandlerFunc {
	return func(context *gin.Context) {
		state := context.Query("state")
		containers, err := api.DockerClient.ListContainers(docker.ListContainersOptions{All: true})
		if err != nil {
			context.JSON(500, utils.HandleResponse("error", err))
		}
		if state != "" {
			containers = filterBy(containers, func(container docker.APIContainers) bool {
				return container.State == state
			})
		}
		context.JSON(200, utils.HandleResponse("containers", containers))
	}
}

func GetContainerDetail(api configs.API) gin.HandlerFunc {
	return func(context *gin.Context) {
		id := context.Params.ByName("id")
		container, err := api.DockerClient.InspectContainer(id)
		if err != nil {
			context.JSON(500, utils.HandleResponse("error", err))
		}
		context.JSON(200, utils.HandleResponse("container", container))
	}
}

func DoProcessOnContainer(api configs.API) gin.HandlerFunc {
	return func(context *gin.Context) {
		action := context.Params.ByName("action")
		id := context.Params.ByName("id")
		if action == "start" {
			err := api.DockerClient.StartContainer(id, &docker.HostConfig{})
			if err != nil {
				context.JSON(500, utils.HandleResponse("result", err))
			} else {
				context.JSON(200, utils.HandleResponse("result", "OK"))
			}
		} else if action == "stop" {
			err := api.DockerClient.StopContainer(id, 120)
			if err != nil {
				context.JSON(500, utils.HandleResponse("result", err))
			} else {
				context.JSON(200, utils.HandleResponse("result", "OK"))
			}
		} else {
			context.JSON(500, utils.HandleResponse("result", "Action undefined.Supported actions start,stop"))
		}
	}
}

func filterBy(containers []docker.APIContainers, filterFunction func(container docker.APIContainers) bool) []docker.APIContainers {
	var filteredAPIContainers = []docker.APIContainers{}
	for _, container := range containers {
		if filterFunction(container) {
			filteredAPIContainers = append(filteredAPIContainers, container)
		}
	}
	return filteredAPIContainers
}
