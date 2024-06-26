// Package changelog provides changelog generator.
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

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/backup"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/klimby/version/internal/service/git"
	"github.com/klimby/version/pkg/convert"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

var (
	// ErrWarning is a warning error.
	ErrWarning = fmt.Errorf("changelog warning")
)

// Generator generates changelog.
type Generator struct {
	repo gitRepo
	rw   readWriter
	bcp  backupSrv
	f    fsys.File
	nms  []config.CommitName
}

// gitRepo is git repository.
type gitRepo interface {
	// Commits returns commits.
	// If nextV is set, then the tag with this version is not created yet and nextV - new created version.
	// In this case will ber returned commits from last tag to HEAD and last commit will be with nextV.
	// If nextV is not set, then will be returned all commits.
	Commits(...func(options *git.CommitsArgs)) ([]git.Commit, error)

	// Add adds files to git.
	Add(files ...fsys.File) error
}

// backupSrv is a backup service.
type backupSrv interface {
	// Create backup of file.
	Create(path string) error
}

type readWriter interface {
	Read(string) (io.ReadCloser, error)
	Write(patch string, flag int) (io.WriteCloser, error)
}

// Args is a Generator arguments.
type Args struct {
	Repo        gitRepo
	RW          readWriter
	Backup      backupSrv
	ConfigFile  fsys.File
	CommitNames []config.CommitName
}

// New creates new Generator.
func New(args ...func(arg *Args)) *Generator {
	a := &Args{
		RW:         fsys.New(),
		Backup:     backup.New(),
		ConfigFile: fsys.File(viper.GetString(key.ChangelogFileName)),
	}

	for _, arg := range args {
		arg(a)
	}

	if a.Repo == nil {
		panic("invalid changelog generator argument: repo is nil")
	}

	return &Generator{
		repo: a.Repo,
		rw:   a.RW,
		bcp:  a.Backup,
		f:    a.ConfigFile,
		nms:  a.CommitNames,
	}
}

// Add adds new version to changelog.
func (g Generator) Add(nextV version.V) (err error) {
	if err := g.bcp.Create(g.f.Path()); err != nil {
		return err
	}

	var b strings.Builder

	if err := g.load(nextV, &b); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return g.generateAll(func(args *git.CommitsArgs) {
				args.NextV = nextV
			})
		}

		return err
	}

	if viper.GetBool(key.DryRun) {
		return nil
	}

	w, err := g.rw.Write(g.f.Path(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
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

	if err := g.repo.Add(g.f); err != nil {
		return fmt.Errorf("add changelog file error: %w", err)
	}

	console.Success(fmt.Sprintf("Changelog %s updated to %s", g.f.String(), nextV.FormatString()))

	return nil
}

// load changes file.
func (g Generator) load(nextV version.V, wr io.Writer) (err error) {
	src, err := g.rw.Read(g.f.Path())
	if err != nil {
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

	if err := g.applyTemplate(&b, func(args *git.CommitsArgs) {
		args.NextV = nextV
		args.LastOnly = true
	}); err != nil {
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

	if viper.GetBool(key.Verbose) {
		console.Info("Changelog changed:")
		console.Info(b.String())
	}

	return nil
}

// Generate generates changelog.
func (g Generator) Generate() (err error) {
	return g.generateAll()
}

// generateAll generates changelog.
func (g Generator) generateAll(opt ...func(*git.CommitsArgs)) (err error) {
	if err := g.rewrite(opt...); err != nil {
		return err
	}

	if err := g.repo.Add(g.f); err != nil {
		return fmt.Errorf("add changelog file error: %w", err)
	}

	console.Success(fmt.Sprintf("Changelog %s created", g.f.String()))

	return nil
}

// rewrite changelog.
func (g Generator) rewrite(opt ...func(*git.CommitsArgs)) (err error) {
	if err := g.bcp.Create(g.f.Path()); err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString("# " + viper.GetString(key.ChangelogTitle) + "\n\n")

	b.WriteString("All notable changes to this project will be documented in this file. ")
	b.WriteString("See [Conventional CommitsFromLast](https://www.conventionalcommits.org/en/v1.0.0/) for commit guidelines.\n")

	if err := g.applyTemplate(&b, opt...); err != nil {
		return err
	}

	if viper.GetBool(key.Verbose) {
		console.Info("Changelog created:")
		console.Info(b.String())
	}

	if viper.GetBool(key.DryRun) {
		return nil
	}

	w, err := g.rw.Write(g.f.Path(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
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
func (g Generator) applyTemplate(wr io.Writer, opt ...func(*git.CommitsArgs)) error {
	c, err := g.repo.Commits(opt...)
	if err != nil {
		return err
	}

	tagsTpl := newTagsTpl(g.nms, c)

	if len(tagsTpl.Tags) == 0 {
		return fmt.Errorf("%w: no new commits", ErrWarning)
	}

	return tagsTpl.applyTemplate(wr)
}

// Normalize changelog builder and convert to []byte.
func builder2B(b strings.Builder) []byte {
	re := regexp.MustCompile(`(\n\s*){2,}`)
	r := re.ReplaceAllString(b.String(), "\n\n")

	return convert.S2B(r)
}
