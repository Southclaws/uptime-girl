package service

import (
	"fmt"
	"strings"

	"github.com/Southclaws/uptime-girl/dockerwatch"
	"github.com/docker/docker/client"
)

// Labels that the service looks out for
const (
	LabelMonitorEndpoint = "com.southclaws.uptime-girl.endpoint"
	LabelMonitorTraefik  = "com.southclaws.uptime-girl.traefik"
)

// Config stores configuration for loading from the environment
type Config struct {
	APIKey string
}

type App struct {
	docker *client.Client
}

// Initialise prepares the service
func Initialise(config Config) (app *App, err error) {
	docker, err := client.NewEnvClient()
	if err != nil {
		return
	}

	app = &App{docker: docker}

	return
}

// Start runs the service and blocks until failure
func (app *App) Start() (err error) {
	tasks := make(chan string)

	dockerwatch.New(docker, LabelMonitorEndpoint, func(value string) {
		if value == "" {
			return
		}

		tasks <- value
	})
	dockerwatch.New(docker, LabelMonitorTraefik, func(value string) {
		if value == "" {
			return
		}

		if traefik, ok := container.Labels[LabelMonitorTraefik]; ok && traefik == "true" {
			hostRules, ok := container.Labels["traefik.frontend.rule"]
			if !ok {
				return
			}
			endpoints := strings.Split(strings.TrimPrefix(hostRules, "Host:"), ",")
			if len(endpoints) == 0 {
				return
			}

			tasks <- value
		}
	})

	for task := range tasks {
		fmt.Println(task)
	}

	return
}
