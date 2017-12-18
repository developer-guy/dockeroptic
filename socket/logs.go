package socket

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/googollee/go-socket.io"
	"github.com/developer-guy/dockeroptic/configs"
	"github.com/gin-gonic/gin"
	"time"
	"log"
)

type logStream struct {
	Socket socketio.Socket
}

func (logStream logStream) Write(p []byte) (n int, error error) {
	time.Sleep(1 * time.Second)
	logStream.Socket.Emit("log_stream", string(p))
	return len(p), nil
}

func HandleSocketConnection(server *socketio.Server, api configs.API) gin.HandlerFunc {
	return func(context *gin.Context) {
		server.On("connection", func(s socketio.Socket) {
			log.Println("Connection established. Socket id :", s.Id())
			s.On("container_log", func(id string) {
				handleContainerLog(id, s, api)
			})
		})

		server.On("error", func(so socketio.Socket, err error) {
			log.Fatal("[ WebSocket ] Error : %v", err.Error())
			panic(err)
		})

		server.ServeHTTP(context.Writer, context.Request)
	}
}

func handleContainerLog(id string, socket socketio.Socket, api configs.API) {
	err := api.DockerClient.Logs(docker.LogsOptions{
		Container: id,
		Follow:    true,
		OutputStream: logStream{
			Socket: socket,
		},
		Stdout:      true,
		RawTerminal: true,
	})

	if err != nil {
		panic(err)
	}
}
