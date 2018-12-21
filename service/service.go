package service

import (
	"strings"

	"github.com/Southclaws/dockwatch"
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

// App stores app state
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
	w := dockwatch.New(app.docker)

	f := func() (e error) {
		select {
		case evt := <-w.Events:
			endpoint := app.react(evt)
			if evt.Type == dockwatch.EventTypeCreate {
				// create monitor
			} else if evt.Type == dockwatch.EventTypeDelete {
				// delete monitor
			}
		case e = <-w.Errors:
		}
		return
	}

	for {
		err = f()
		if err != nil {
			break
		}
	}

	return
}

func (app *App) react(e dockwatch.Event) (endpoint string) {
	if endpoint, ok := e.Container.Labels[LabelMonitorEndpoint]; ok {
		return endpoint
	}
	if traefik, ok := e.Container.Labels[LabelMonitorTraefik]; ok && traefik == "true" {
		hostRules, ok := e.Container.Labels["traefik.frontend.rule"]
		if !ok {
			return
		}
		endpoints := strings.Split(strings.TrimPrefix(hostRules, "Host:"), ",")
		if len(endpoints) == 0 {
			return
		}
		return endpoints[0]
	}
	return
}
