package changelog

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/klimby/version/internal/backup"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/file"
	"github.com/spf13/viper"
)

var (
	errGenerate = fmt.Errorf("changelog generate error")
)

type Generator struct {
	repo gitRepo
	f    file.ReadWriter
	path string
}

type gitRepo interface {
	LastCommit() (*object.Commit, error)
}

func NewGenerator(f file.ReadWriter, g gitRepo) *Generator {
	n := config.File(viper.GetString(config.ChangelogFileName))
	return &Generator{
		repo: g,
		f:    f,
		path: n.Path(),
	}
}

func (g Generator) Generate() error {
	if err := g.check(); err != nil {
		return err
	}

	if err := backup.Create(g.f, g.path); err != nil {
		return err
	}

	return nil
}

func (g Generator) check() error {
	last, err := g.repo.LastCommit()
	if err != nil {
		return err
	}

	if last.MergeTag == "" {
		return fmt.Errorf("%w: last commit is not tagged (last commit must be version)", errGenerate)
	}

	return nil

}
