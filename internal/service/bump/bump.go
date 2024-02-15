package bump

import (
	"fmt"
	"io"
	"os"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/backup"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/klimby/version/pkg/convert"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

// B bump files by version.
type B struct {
	rw      readWriter
	repo    gitRepo
	bcp     backupSrv
	process contentProcessor
}

type gitRepo interface {
	Add(files ...fsys.File) error
}

type backupSrv interface {
	Create(path string) error
}

type readWriter interface {
	Read(string) (io.ReadCloser, error)
	Write(patch string, flag int) (io.WriteCloser, error)
}

type contentProcessor interface {
	CustomFile(r io.Reader, bmp config.BumpFile, v version.V) ([]string, bool, error)
	PredefinedJSON(r io.Reader, bmp config.BumpFile, v version.V) ([]string, bool, error)
}

// Args is a Bump arguments.
type Args struct {
	RW     readWriter
	Repo   gitRepo
	Backup backupSrv
	Proc   contentProcessor
}

// New creates new Bump.
func New(args ...func(arg *Args)) *B {
	a := &Args{
		RW:     fsys.New(),
		Backup: backup.New(),
		Proc:   &process{},
	}

	for _, arg := range args {
		arg(a)
	}

	return &B{
		rw:      a.RW,
		repo:    a.Repo,
		bcp:     a.Backup,
		process: a.Proc,
	}
}

// Apply bumps files.
func (b B) Apply(bumps []config.BumpFile, v version.V) {
	for _, bmp := range bumps {
		if err := b.bcp.Create(bmp.File.Path()); err != nil {
			console.Error(fmt.Sprintf("create backup file %s error: %s", bmp.File.String(), err.Error()))
		}

		changed, err := b.applyToFile(bmp, v)
		if err != nil {
			console.Warn(fmt.Sprintf("bump file %s error: %s", bmp.File.String(), err.Error()))

			continue
		}

		if changed {
			if err := b.repo.Add(bmp.File); err != nil {
				console.Warn(fmt.Sprintf("add file %s to git error: %s", bmp.File.String(), err.Error()))
			}
		}
	}
}

// applyToFile bumps file.
func (b B) applyToFile(bmp config.BumpFile, v version.V) (bool, error) {
	content, changed, err := b.read(bmp, v)
	if err != nil {
		return false, fmt.Errorf("bump file %s error: %w", bmp.File.String(), err)
	}

	if len(content) == 0 {
		return false, fmt.Errorf("file %s is empty", bmp.File.String())
	}

	if !changed {
		console.Warn(fmt.Sprintf("File %s is not changed", bmp.File.String()))

		return false, nil
	}

	if !viper.GetBool(key.DryRun) {
		if err := b.write(bmp.File.Path(), content); err != nil {
			return false, fmt.Errorf("write file %s error: %w", bmp.File.String(), err)
		}
	}

	console.Success(fmt.Sprintf("Bump file %s", bmp.File.String()))

	return changed, nil
}

// read reads file.
func (b B) read(bmp config.BumpFile, v version.V) (_ []string, changed bool, err error) {
	r, err := b.rw.Read(bmp.File.Path())
	if err != nil {
		return nil, false, fmt.Errorf("open file %s error: %w", bmp.File.Path(), err)
	}

	defer func() {
		if e := r.Close(); e != nil {
			if err == nil {
				err = fmt.Errorf("close file %s error: %w", bmp.File.Path(), e)
			}
		}
	}()

	if bmp.IsPredefinedJSON() {
		return b.process.PredefinedJSON(r, bmp, v)
	}

	return b.process.CustomFile(r, bmp, v)
}

// write writes content to file.
func (b B) write(patch string, content []string) (err error) {
	w, err := b.rw.Write(patch, os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return fmt.Errorf("open file %s error: %w", patch, err)
	}

	defer func() {
		if e := w.Close(); e != nil {
			if err == nil {
				err = fmt.Errorf("close file %s error: %w", patch, e)
			}
		}
	}()

	for _, line := range content {
		_, err := w.Write(convert.S2B(line + "\n"))
		if err != nil {
			return fmt.Errorf("write file %s line %s error: %w", patch, line, err)
		}
	}

	return nil
}
