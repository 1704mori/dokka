package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"

	env "github.com/1704mori/dokka/backend"
	"github.com/1704mori/dokka/backend/api"
	"github.com/1704mori/dokka/backend/docker"
	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version = "dev"

//go:embed dist
var frontend embed.FS

func main() {
	arg.MustParse(&env.Args)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msgf("Starting dokka version %s", Version)

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		log.Error().Msg("Failed to create new docker client (cli)")
		return
	}

	router := gin.Default()
	router.Use(corsMiddleware())

	sub, err := fs.Sub(frontend, "dist")
	if err != nil {
		panic(err)
	}

	router.GET("/assets/*filepath", func(c *gin.Context) {
		http.FileServer(http.FS(sub)).ServeHTTP(c.Writer, c.Request)
	})
	router.NoRoute(func(c *gin.Context) {
		executeTemplate(c)
	})
	api.APIRoutes(router, dockerClient)

	router.Run(fmt.Sprintf(":%s", env.Args.Port))
	log.Info().Msgf("Listening on port %v", env.Args.Port)
}

func executeTemplate(c *gin.Context) {
	file, err := os.Open("dist/index.html")
	if err != nil {
		log.Panic().Msgf("Failed to open index.html: %v", err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Panic().Msgf("Failed to read index.html: %v", err)
	}

	file.Close()
	log.Info().Msgf("content: %s", string(bytes))

	mashaled, err := json.Marshal(env.Args)
	if err != nil {
		log.Panic().Msgf("Failed to marshal config: %v", err)
	}
	log.Info().Msgf("config: %s", string(mashaled))

	tmpl, err := template.New("dist/index.html").Funcs(template.FuncMap{
		"env": func() template.JS {
			return template.JS(fmt.Sprintf(`{"env": %s}`, string(mashaled)))
		},
	}).Parse(string(bytes))
	if err != nil {
		log.Panic().Msgf("Failed to parse index.html: %v", err)
	}

	err = tmpl.Execute(os.Stdout, nil)
	if err != nil {
		log.Panic().Msgf("Failed to execute index.html: %v", err)
	}

	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		log.Panic().Msgf("Failed to execute index.html: %v", err)
	}
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
