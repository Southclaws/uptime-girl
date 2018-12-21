package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/Southclaws/uptime-girl/service"
	_ "github.com/joho/godotenv/autoload" // load vars from .env
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

func main() {
	config := service.Config{}
	envconfig.MustProcess("", &config)

	app, err := service.Initialise(config)
	if err != nil {
		zap.L().Fatal("failed to initialise", zap.Error(err))
	}

	zap.L().Info("service initialised")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := make(chan error, 1)
	go func() { e <- app.Start(ctx) }()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)

	select {
	case sig := <-s:
		err = errors.New(sig.String())
	case err = <-e:
	}

	zap.L().Fatal("exit", zap.Error(err))
}
