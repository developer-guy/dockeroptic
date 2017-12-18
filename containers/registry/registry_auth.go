package registry

import (
	"github.com/gin-gonic/gin"
	"github.com/developer-guy/dockeroptic/utils"
	"encoding/base64"
	"strings"
	"github.com/developer-guy/dockeroptic/configs"
	"github.com/developer-guy/dockeroptic/environments"
)

type loginCredentials struct {
	username string `json:"username"`
	password string `json:"password"`
}

func basicAuth(credentials loginCredentials) string {
	auth := credentials.username + ":" + credentials.password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func AuthRequired(api configs.API) gin.HandlerFunc {
	return func(context *gin.Context) {
		var headers utils.Headers = make(map[string]string)
		authorization := context.GetHeader("Authorization")

		if authorization == "" {
			context.JSON(500, gin.H{
				"error": "Basic Auth Required",
			})
			context.Abort()
		}

		authorization = authorization[6:]

		decodedCredential, err := base64.StdEncoding.DecodeString(authorization)

		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			context.Abort()
		}

		credentials := strings.Split(string(decodedCredential), ":")

		var loginCredentials = loginCredentials{
			username: credentials[0],
			password: credentials[1],
		}
		token := basicAuth(loginCredentials)

		api.RedisClient.Set("currentUser", loginCredentials.username, 0).Result()
		api.RedisClient.Set(loginCredentials.username, token, 0)

		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			context.Abort()
		}

		headers.Add("Authorization", "Basic "+token)

		resp, err := utils.HandleRequest("GET", environments.DockerRegistryBaseUri+"/v2/", nil, headers)

		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			context.Abort()
		}

		status := resp.Status[:3]

		if status != "200" {
			context.JSON(500, gin.H{
				"error": "Authorization Failed",
			})
			context.Abort()
		}
	}
}
