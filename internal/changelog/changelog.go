package changelog

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strings"

	"github.com/klimby/version/internal/backup"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/internal/git"
	"github.com/klimby/version/pkg/convert"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

var (
	ErrWarning  = fmt.Errorf("changelog warning")
	errGenerate = fmt.Errorf("changelog generate error")
)

type Generator struct {
	repo gitRepo
	f    file.ReadWriter
	path string
}

type gitRepo interface {
	// Commits returns commits.
	// If nextV is set, then the tag with this version is not created yet and nextV - new created version.
	// In this case will ber returned commits from last tag to HEAD and last commit will be with nextV.
	// If nextV is not set, then will be returned all commits.
	Commits(nextV ...version.V) ([]git.Commit, error)
}

func NewGenerator(f file.ReadWriter, g gitRepo) *Generator {
	n := config.File(viper.GetString(config.ChangelogFileName))
	return &Generator{
		repo: g,
		f:    f,
		path: n.Path(),
	}
}

func (g Generator) Add(nextV version.V) (err error) {
	if err := backup.Create(g.f, g.path); err != nil {
		return err
	}

	var b strings.Builder

	if err := g.load(nextV, &b); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return g.Generate()
		}

		return err
	}

	if viper.GetBool(config.DryRun) {
		return nil
	}

	w, err := g.f.Write(g.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return fmt.Errorf("open changelog file error: %w", err)
	}

	defer func() {
		if e := w.Close(); e != nil {
			if err == nil {
				err = fmt.Errorf("close config file error: %w", e)
			}
		}
	}()

	if _, err := w.Write(builder2B(b)); err != nil {
		return fmt.Errorf("write changelog file error: %w", err)
	}

	return nil
}

// load changes file.
func (g Generator) load(nextV version.V, wr io.Writer) (err error) {
	src, err := g.f.Read(g.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return g.Generate()
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

	var b strings.Builder

	if err := g.applyTemplate(&b, nextV); err != nil {
		return err
	}

	// read file by lines
	scanner := bufio.NewScanner(src)

	insert := false

	for scanner.Scan() {
		byteLine := append(scanner.Bytes(), '\n')
		// Insert before first ## line.
		if !insert && strings.HasPrefix(scanner.Text(), "##") {
			insert = true

			if _, err := wr.Write(convert.S2B(b.String() + "\n")); err != nil {
				return err
			}

			if _, err := wr.Write(byteLine); err != nil {
				return err
			}
		} else {
			if _, err := wr.Write(byteLine); err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// insert after last line
	if !insert {
		if _, err := wr.Write(convert.S2B(b.String() + "\n")); err != nil {
			return err
		}
	}

	if !viper.GetBool(config.Silent) {
		console.Notice("\nChangelog changed:\n")
		console.Info(b.String())
	}

	return nil
}

func (g Generator) Generate() (err error) {
	if err := backup.Create(g.f, g.path); err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString("# " + viper.GetString(config.ChangelogTitle) + "\n\n")

	b.WriteString("All notable changes to this project will be documented in this file. ")
	b.WriteString("See [Conventional CommitsFromLast](https://www.conventionalcommits.org/en/v1.0.0/) for commit guidelines.\n")

	if err := g.applyTemplate(&b); err != nil {
		return err
	}

	if !viper.GetBool(config.Silent) {
		console.Info("\n")
		console.Info(b.String())
	}

	if viper.GetBool(config.DryRun) {
		return nil
	}

	w, err := g.f.Write(g.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return fmt.Errorf("open changelog file error: %w", err)
	}

	defer func() {
		if e := w.Close(); e != nil {
			if err == nil {
				err = fmt.Errorf("close config file error: %w", e)
			}
		}
	}()

	if _, err := w.Write(builder2B(b)); err != nil {
		return fmt.Errorf("write changelog file error: %w", err)
	}

	return nil
}

// applyTemplate applies template to writer.
func (g Generator) applyTemplate(wr io.Writer, nextV ...version.V) error {
	c, err := g.repo.Commits(nextV...)
	if err != nil {
		return err
	}

	tagsTpl, err := NewTagsTpl(c)
	if err != nil {
		return err
	}

	if len(tagsTpl.Tags) == 0 {
		return fmt.Errorf("%w: no new commits", ErrWarning)
	}

	if err := tagsTpl.applyTemplate(wr); err != nil {
		return err
	}

	return nil
}

// Normalize changelog builder and convert to []byte.
func builder2B(b strings.Builder) []byte {
	re := regexp.MustCompile(`(\n\s*){2,}`)
	r := re.ReplaceAllString(b.String(), "\n\n")

	return convert.S2B(r)

}
