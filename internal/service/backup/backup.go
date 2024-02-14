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

// Service is a backup service.
type Service struct {
	rw rwService
}

// rwService is a read/write service.
type rwService interface {
	Read(patch string) (io.ReadCloser, error)
	Write(patch string, flag int) (io.WriteCloser, error)
	RemoveAll(patch string) error
}

// Args is a Service arguments.
type Args struct {
	RW rwService
}

// New creates new Service.
func New(args ...func(arg *Args)) *Service {
	a := &Args{
		RW: fsys.New(),
	}

	for _, arg := range args {
		arg(a)
	}

	return &Service{
		rw: a.RW,
	}
}

// Create backup of file.
func (s Service) Create(path string) (err error) {
	if !viper.GetBool(key.Backup) || viper.GetBool(key.DryRun) {
		return nil
	}

	backPath := path + _suffix

	src, err := s.rw.Read(path)
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

	back, err := s.rw.Write(backPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
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
func (s Service) Remove(path ...string) {
	for _, p := range path {
		backPath := p + _suffix

		if err := s.rw.RemoveAll(backPath); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}

			console.Error("RemoveAll backup file " + backPath + " error: " + err.Error())

			continue
		}

		console.Success("Backup file " + backPath + " removed.")
	}
}
