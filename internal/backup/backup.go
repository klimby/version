package backup

import (
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/file"
	"github.com/spf13/viper"
)

const (
	_suffix = ".bak"
)

// Create backup of file.
func Create(f file.ReadWriter, path string) (err error) {
	if !viper.GetBool(config.Backup) || viper.GetBool(config.DryRun) {
		return nil
	}

	backPath := path + _suffix

	src, err := f.Read(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		return err
	}

	defer func() {
		if e := src.Close(); e != nil {
			if err == nil {
				err = e
			}
		}
	}()

	back, err := f.Write(backPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return err
	}

	defer func() {
		if e := back.Close(); e != nil {
			if err == nil {
				err = e
			}
		}
	}()

	_, err = io.Copy(back, src)
	if err != nil {
		return err
	}

	console.Info("Backup file " + backPath + " created.")

	return nil
}

// Remove backup of file.
func Remove(f file.Remover, path ...string) {

	for _, p := range path {
		backPath := p + _suffix

		if err := f.Remove(backPath); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}

			console.Error("Remove backup file " + backPath + " error: " + err.Error())

			continue
		}

		console.Info("Backup file " + backPath + " removed.")
	}
}
