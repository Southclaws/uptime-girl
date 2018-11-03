package main

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/Southclaws/uptime-girl/server"
)

func main() {
	config := server.Config{}
	envconfig.MustProcess("", &config)

	app, err := server.Initialise(config)
	if err != nil {
		fmt.Println("failed to initialise:", err)
		os.Exit(2)
	}

	os.Exit(app.Start())
}
