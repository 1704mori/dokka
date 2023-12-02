package api

import (
	"encoding/json"
	"fmt"
	"time"

	env "github.com/1704mori/dokka/backend"
	"github.com/1704mori/dokka/backend/docker"
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

					sse <- "data: " + string(marshal) + "\n\n"
					time.Sleep(time.Duration(env.Args.EventStreamInterval) * time.Second)
				}
			}()

			for {
				msg := <-sse
				ctx.SSEvent("message", msg)
				ctx.Writer.Flush()
			}
		})
		v1.POST("/container/:action/:id", func(ctx *gin.Context) {
			action := ctx.Param("action")
			id := ctx.Param("id")

			if id == "" {
				ctx.JSON(400, gin.H{
					"result":  "error",
					"message": "missing container id parameter",
				})
				return
			}

			if action == "" {
				ctx.JSON(400, gin.H{
					"result":  "error",
					"message": "missing action parameter",
				})
				return
			}

			container, err := dockerClient.FindContainer(id)
			if err != nil {
				ctx.JSON(400, gin.H{
					"result":  "error",
					"message": fmt.Sprintf("container with id %s not found", id)})
				return
			}

			if container.ID == "" {
				ctx.JSON(400, gin.H{
					"result":  "error",
					"message": fmt.Sprintf("container with id %s not found", id)})
				return
			}

			switch action {
			case "stop":
				_, err = dockerClient.StopContainer(id)
				if err != nil {
					ctx.JSON(500, gin.H{
						"result":  "error",
						"message": fmt.Sprintf("failed to stop container: %v", err),
					})
					return
				}
				ctx.JSON(200, gin.H{
					"result":  "ok",
					"message": fmt.Sprintf("container %s stopped successfully", id),
				})
			case "start":
				_, err = dockerClient.StartContainer(id)
				if err != nil {
					ctx.JSON(500, gin.H{
						"result":  "error",
						"message": fmt.Sprintf("failed to start container: %v", err),
					})
					return
				}
				ctx.JSON(200, gin.H{
					"result":  "ok",
					"message": fmt.Sprintf("container %s started successfully", id),
				})
			case "restart":
				_, err = dockerClient.RestartContainer(id)
				if err != nil {
					ctx.JSON(500, gin.H{
						"result":  "error",
						"message": fmt.Sprintf("failed to restart container: %v", err),
					})
					return
				}
				ctx.JSON(200, gin.H{
					"result":  "ok",
					"message": fmt.Sprintf("container %s restarted successfully", id),
				})

			}
		})
	}
}
