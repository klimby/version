package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

// File is a file wrapper.
type File string

// String returns the string representation of the file.
func (f File) String() string {
	return string(f)
}

// Path returns the file path.
func (f File) Path() string {
	return filepath.Join(viper.GetString(WorkDir), f.String())
}

// Rel returns the relative path to the file.
func (f File) Rel() string {
	r, err := filepath.Rel(viper.GetString(WorkDir), f.Path())
	if err != nil {
		return f.String()
	}

	return r
}
