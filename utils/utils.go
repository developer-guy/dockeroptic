package utils

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)


func HandleResponse(key string, obj interface{}) gin.H {
	return gin.H{
		key: obj,
	}
}

func HandleRequest(method string, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	client := http.DefaultClient
	request, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	return client.Do(request)
}
