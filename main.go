package main

import (
	"os"

	"github.com/klimby/version/cmd"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/di"
)

// -ldflags "-X main.version=$VERSION"
var version = "0.0.0"

func main() {
	code := 0
	defer os.Exit(code)

	config.Init(func(options *config.ConfigOptions) {
		options.Version = version
	})

	if err := di.C.Init(""); err != nil {
		code = errorHandler(err)

		return
	}

	if err := cmd.Execute(); err != nil {
		code = errorHandler(err)
	}
}

func errorHandler(err error) int {
	if err != nil {
		console.Error(err.Error())

		return 1
	}

	return 0
}
