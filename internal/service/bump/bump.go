package bump

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

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
	rw   fsys.ReadWriter
	repo gitRepo
}

type gitRepo interface {
	Add(files ...fsys.File) error
}

// Args is a Bump arguments.
type Args struct {
	RW   fsys.ReadWriter
	Repo gitRepo
}

// New creates new Bump.
func New(args ...func(arg *Args)) *B {
	a := &Args{
		RW: fsys.NewFS(),
	}

	for _, arg := range args {
		arg(a)
	}

	if a.Repo == nil {
		panic("invalid backup argument: repo is nil")
	}

	return &B{
		rw:   a.RW,
		repo: a.Repo,
	}
}

// Apply bumps files.
func (b B) Apply(bumps []config.BumpFile, v version.V) {
	for _, bmp := range bumps {
		if err := backup.Create(b.rw, bmp.File.Path()); err != nil {
			console.Error(err.Error())
		}

		changed, err := b.applyToFile(bmp, v)
		if err != nil {
			console.Warn(err.Error())
		}

		if changed {
			if err := b.repo.Add(bmp.File); err != nil {
				console.Warn(err.Error())
			}
		}
	}
}

// applyToFile bumps file.
func (b B) applyToFile(bmp config.BumpFile, v version.V) (bool, error) {
	var content []string
	changed := false

	if bmp.IsPredefinedJSON() {
		c, ch, err := b.bumpPredefinedJSON(bmp.File.Path(), v)
		if err != nil {
			return false, fmt.Errorf("bump predefined json file %s error: %w", bmp.File.String(), err)
		}

		content = c
		changed = ch
	} else {
		c, ch, err := b.bumpCustomFile(bmp, v)
		if err != nil {
			return false, fmt.Errorf("bump custom file %s error: %w", bmp.File.String(), err)
		}

		content = c
		changed = ch
	}

	if len(content) == 0 {
		return false, fmt.Errorf("file %s is not supported", bmp.File.String())
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

// bumpPredefinedJSON bumps predefined json file (package.json, composer.json).
func (b B) bumpPredefinedJSON(patch string, v version.V) (_ []string, changed bool, err error) {
	r, err := b.rw.Read(patch)
	if err != nil {
		return nil, changed, fmt.Errorf("open file %s error: %w", patch, err)
	}

	defer func() {
		if e := r.Close(); e != nil {
			if err == nil {
				err = fmt.Errorf("close file %s error: %w", patch, e)
			}
		}
	}()

	scanner := bufio.NewScanner(r)
	versionRegex := regexp.MustCompile(`"version"\s*:\s*".*?"`)
	var content []string

	replacer := `"version": "` + v.FormatString() + `"`

	for scanner.Scan() {
		line := scanner.Text()

		if !changed && versionRegex.MatchString(line) {
			line = versionRegex.ReplaceAllString(line, replacer)
			changed = true
		}

		content = append(content, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, changed, fmt.Errorf("scan file %s error: %w", patch, err)
	}

	return content, changed, nil
}

// bumpCustomFile bumps custom file.
func (b B) bumpCustomFile(bmp config.BumpFile, v version.V) (_ []string, changed bool, err error) {
	r, err := b.rw.Read(bmp.File.Path())
	if err != nil {
		return nil, changed, fmt.Errorf("open file %s error: %w", bmp.File.Path(), err)
	}

	defer func() {
		if e := r.Close(); e != nil {
			if err == nil {
				err = fmt.Errorf("close file %s error: %w", bmp.File.Path(), e)
			}
		}
	}()

	scanner := bufio.NewScanner(r)
	versionRegex := regexp.MustCompile(`\d+\.\d+\.\d+`)
	var content []string

	regArr := make([]*regexp.Regexp, 0, len(bmp.RegExp))

	for _, r := range bmp.RegExp {
		rgx, err := regexp.Compile(r)
		if err != nil {
			console.Warn(fmt.Sprintf("file %s regexp %s error: %s", bmp.File.String(), r, err))

			continue
		}

		regArr = append(regArr, rgx)
	}

	start := 0
	end := int(^uint(0) >> 1)

	if bmp.HasPositions() {
		start = bmp.Start
		end = bmp.End
	}

	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()

		if lineNum >= start && lineNum <= end {
			if versionRegex.MatchString(line) {
				if len(regArr) > 0 {
					for _, r := range regArr {
						if r.MatchString(line) {
							line = versionRegex.ReplaceAllString(line, v.FormatString())
							changed = true

							break
						}
					}
				} else {
					line = versionRegex.ReplaceAllString(line, v.FormatString())
					changed = true
				}
			}
		}

		content = append(content, line)

		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, changed, fmt.Errorf("scan file %s error: %w", bmp.File.String(), err)
	}

	return content, changed, nil
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
