package backup

import (
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/spf13/viper"
)

const (
	_suffix = ".bak"
)

// Create backup of file.
func Create(f fsys.ReadWriter, path string) (err error) {
	if !viper.GetBool(key.Backup) || viper.GetBool(key.DryRun) {
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

	console.Success("Backup file " + backPath + " created.")

	return nil
}

// Remove backup of file.
func Remove(f fsys.Remover, path ...string) {
	for _, p := range path {
		backPath := p + _suffix

		if err := f.Remove(backPath); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}

			console.Error("Remove backup file " + backPath + " error: " + err.Error())

			continue
		}

		console.Success("Backup file " + backPath + " removed.")
	}
}
