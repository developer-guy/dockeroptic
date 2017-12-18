package local

import (
	"github.com/developer-guy/dockeroptic/configs"
	"github.com/gin-gonic/gin"
	"github.com/fsouza/go-dockerclient"
	"github.com/developer-guy/dockeroptic/utils"
)

func GetAll(api configs.API) gin.HandlerFunc {
	return func(context *gin.Context) {
		images, err := api.DockerClient.ListImages(docker.ListImagesOptions{
			All: true,
		})

		if err != nil {
			context.JSON(500, utils.HandleResponse("error", err.Error()))
		}

		context.JSON(200, utils.HandleResponse("images", images))
	}
}

func GetImageDetail(api configs.API) gin.HandlerFunc {
	return func(context *gin.Context) {
		id := context.Params.ByName("id")
		image, err := api.DockerClient.InspectImage(id)

		if err != nil {
			context.JSON(500, utils.HandleResponse("error", err.Error()))
		}

		context.JSON(200, utils.HandleResponse("image", image))
	}
}

func GetImageHistory(api configs.API) gin.HandlerFunc {
	return func(context *gin.Context) {
		id := context.Params.ByName("id")
		imageHistory, err := api.DockerClient.ImageHistory(id)

		if err != nil {
			context.JSON(500, utils.HandleResponse("error", err.Error()))
		}

		context.JSON(200, utils.HandleResponse("history", imageHistory))
	}
}
