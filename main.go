package main

import (
	"os"

	"github.com/klimby/version/cmd"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/service/console"
)

// -ldflags "-X main.version=$VERSION"
var version = "0.0.0"

func main() {
	config.Init(func(options *config.Options) {
		options.Version = version
	})

	// Di init in cmd.rootCmd in cmd.initConfig()
	if err := cmd.Execute(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
