package config

import "github.com/klimby/version/pkg/version"

type C struct {
	// Version is a version of the application.
	Version version.V `yaml:"version"`
}
