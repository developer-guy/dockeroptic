package registry

import (
	"github.com/gin-gonic/gin"
	"github.com/developer-guy/dockeroptic/utils"
	"io/ioutil"
	"github.com/developer-guy/dockeroptic/configs"
	"github.com/developer-guy/dockeroptic/environments"
)

func GetRepositories(api configs.API) gin.HandlerFunc {
	return func(context *gin.Context) {
		var headers utils.Headers = make(map[string]string)
		currentUser, _ := api.RedisClient.Get("currentUser").Result()
		token, _ := api.RedisClient.Get(currentUser).Result()
		headers.Add("Authorization", "Basic "+token)

		resp, err := utils.HandleRequest("GET", environments.DockerRegistryBaseUri+"/v2/_catalog", nil, headers)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)

		context.Header("Content-Type", "application/json")
		context.String(200, string(body))
	}
}
