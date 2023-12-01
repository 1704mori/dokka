package api

import (
	"encoding/json"
	"time"

	"github.com/1704mori/dokka/docker"
	"github.com/gin-gonic/gin"
)

func APIRoutes(router *gin.Engine, dockerClient *docker.Client) {
	v1 := router.Group("/v1")
	{
		v1.GET("/containers", func(ctx *gin.Context) {
			headers := map[string]string{
				"Content-Type":  "text/event-stream",
				"Cache-Control": "no-cache",
				"Connection":    "keep-alive",
			}

			for key, value := range headers {
				ctx.Writer.Header().Set(key, value)
			}

			sse := make(chan string)
			defer close(sse)

			go func() {
				for {
					containers, err := dockerClient.ListContainers()
					if err != nil {
						sse <- "error: " + err.Error()
					}

					marshal, err := json.Marshal(containers)
					if err != nil {
						sse <- "error: " + err.Error()
					}

					sse <- "data: " + string(marshal)
					time.Sleep(5 * time.Second)
				}
			}()

			for {
				msg := <-sse
				ctx.SSEvent("message", msg)
				ctx.Writer.Flush()
			}
		})
	}
}
