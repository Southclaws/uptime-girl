package service

import (
	"context"
	"strings"

	"github.com/Southclaws/dockwatch"
	"github.com/docker/docker/client"
	"go.uber.org/zap"

	"github.com/Southclaws/uptime-girl/uptimerobot"
)

// Labels that the service looks out for
const (
	LabelMonitorEndpoint = "com.southclaws.uptime-girl.endpoint"
	LabelMonitorTraefik  = "com.southclaws.uptime-girl.traefik"
)

// Config stores configuration for loading from the environment
type Config struct {
	UptimeRobotKey string
}

// App stores app state
type App struct {
	docker *client.Client
	uptime *uptimerobot.Client
}

// Initialise prepares the service
func Initialise(config Config) (app *App, err error) {
	docker, err := client.NewEnvClient()
	if err != nil {
		return
	}

	app = &App{
		docker: docker,
		uptime: uptimerobot.New(config.UptimeRobotKey),
	}

	return
}

// Start runs the service and blocks until failure
func (app *App) Start(ctx context.Context) (err error) {
	w := dockwatch.New(app.docker)

	f := func() error {
		var e error
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e = <-w.Errors:

		case evt := <-w.Events:
			e = app.react(evt)
		}
		if e != nil {
			zap.L().Error("error occurred", zap.Error(e))
		}
		return nil
	}

	for {
		f()
		if err != nil {
			break
		}
	}

	return
}

func (app *App) react(e dockwatch.Event) (err error) {
	endpoint := app.getEndpoint(e)
	if endpoint == "" {
		return
	}

	if e.Type == dockwatch.EventTypeCreate {
		app.uptime.NewMonitor(uptimerobot.Monitor{
			URL:          endpoint,
			FriendlyName: endpoint,
		})
	}

	if e.Type == dockwatch.EventTypeDelete {
		var monitors []uptimerobot.Monitor
		monitors, err = app.uptime.GetMonitors()
		if err != nil {
			return
		}
		for _, m := range monitors {
			if m.URL == endpoint {
				app.uptime.DeleteMonitor(m.ID)
				break
			}
		}
	}

	return
}

func (app *App) getEndpoint(e dockwatch.Event) (endpoint string) {
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
