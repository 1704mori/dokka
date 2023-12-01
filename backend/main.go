package main

import (
	"fmt"

	"github.com/1704mori/dokka/api"
	"github.com/1704mori/dokka/docker"
	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version = "dev"

var args struct {
	Port string `arg:"env:DOKKA_PORT" default:"7070"`
}

func main() {
	arg.MustParse(&args)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msgf("Starting dokka version %s", Version)

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		log.Error().Msg("Failed to create new docker client (cli)")
		return
	}

	router := gin.Default()
	router.Use(corsMiddleware())
	api.APIRoutes(router, dockerClient)

	router.Run(fmt.Sprintf(":%s", args.Port))
	// err = http.ListenAndServe(fmt.Sprintf(":%s", args.Port), router)
	// log.Fatal().Msg(err.Error())
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
