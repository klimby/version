package fsys

import (
	"path/filepath"

	"github.com/klimby/version/internal/config/key"
	"github.com/spf13/viper"
)

// File is a file wrapper.
type File string

// String returns the string representation of the file.
func (f File) String() string {
	return string(f)
}

// IsAbs returns true if the file path is absolute.
func (f File) IsAbs() bool {
	return filepath.IsAbs(f.String())
}

// Empty returns true if the file is empty.
func (f File) Empty() bool {
	return f.String() == ""
}

// Path returns the file path.
func (f File) Path() string {
	p := f.String()

	if f.IsAbs() {
		return p
	}

	return filepath.Join(viper.GetString(key.WorkDir), p)
}

// Rel returns the relative path to the file.
func (f File) Rel() string {
	r, err := filepath.Rel(viper.GetString(key.WorkDir), f.Path())
	if err != nil {
		return f.String()
	}

	return r
}
