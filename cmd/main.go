package main

import (
	"context"
	"fmt"
	"os"

	flags "github.com/spf13/pflag"
)

// -ldflags "-X main.version=$VERSION"
var version = "0.0.0"

func main() {
	ctx := context.Background()
	code := 0
	defer os.Exit(code)

	fmt.Println("version:", version)

	var (
		major   bool
		minor   bool
		patch   bool
		version string
		dryRun  bool
		silent  bool

		generateConfig bool
	)

	flagSet := flags.NewFlagSet("version", flags.ContinueOnError)
	flagSet.SortFlags = false

	flagSet.BoolVar(&major, "major", false, "Next major version")
	flagSet.BoolVar(&minor, "minor", false, "Next minor version")
	flagSet.BoolVar(&patch, "patch", false, "Next patch version")
	flagSet.StringVar(&version, "version", "", "Version in format: 1.0.0")

	flagSet.BoolVarP(&dryRun, "dry", "d", false, "Dry run")
	flagSet.BoolVarP(&silent, "silent", "s", false, "Silent mode")

	flagSet.BoolVar(&generateConfig, "generate-config", false, "Generate or update config file")

	flags.Parse()

	switch {
	case generateConfig:
		fmt.Println("generate config")

	case major:
		fmt.Println("major")

	case minor:
		fmt.Println("minor")

	case patch:
		fmt.Println("patch")

	case version != "":
		fmt.Println("version")

	default:
		flagSet.PrintDefaults()
	}

	_ = ctx

	return

}
