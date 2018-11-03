package server

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Labels that the service looks out for
const (
	LabelMonitorEndpoint = "com.southclaws.uptime-girl.endpoint"
	LabelMonitorTraefik  = "com.southclaws.uptime-girl.traefik"
)

type Config struct {
}

type App struct {
	config Config
	docker *client.Client
}

func Initialise(config Config) (app *App, err error) {
	client, err := client.NewEnvClient()
	if err != nil {
		return
	}

	app = &App{
		config: config,
		docker: client,
	}
	return
}

func (app *App) Start() (exit int) {
	err := app.start()
	if err != nil {
		log.Print(err)
		return 1
	}
	return 0
}

func (app *App) start() (err error) {
	containers, err := app.docker.ContainerList(
		context.Background(),
		types.ContainerListOptions{})
	if err != nil {
		return
	}
	for _, container := range containers {
		endpoint := app.checkContainer(container)
		if endpoint != "" {
			// create monitor
		}
	}
	return
}

func (app *App) checkContainer(container types.Container) (endpoint string) {
	if endpoint, ok := container.Labels[LabelMonitorEndpoint]; ok {
		return endpoint
	}
	if traefik, ok := container.Labels[LabelMonitorTraefik]; ok && traefik == "true" {
		hostRules, ok := container.Labels["traefik.frontend.rule"]
		if !ok {
			return
		}
		endpoints := strings.Split(strings.TrimPrefix(hostRules, "Host:"), ",")
		if len(endpoints) > 0 {
			return endpoints[0]
		}
	}
	return
}
