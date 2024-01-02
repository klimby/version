package bump

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/klimby/version/internal/backup"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/pkg/convert"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

func Apply(f file.ReadWriter, bumps []config.BumpFile, v version.V) {
	for _, bmp := range bumps {
		if err := backup.Create(f, bmp.File.Path()); err != nil {
			console.Error(err.Error())
		}

		if err := applyToFile(f, bmp, v); err != nil {
			console.Error(err.Error())
		}
	}
}

func applyToFile(f file.ReadWriter, bmp config.BumpFile, v version.V) error {

	var content [][]byte
	changed := false

	if bmp.IsPredefinedJSON() {
		c, ch, err := bumpPredefinedJSON(f, bmp.File.Path(), v)
		if err != nil {
			return fmt.Errorf("bump predefined json file %s error: %w", bmp.File.String(), err)
		}

		content = c
		changed = ch
	} else {
		c, ch, err := bumpCustomFile(f, bmp, v)
		if err != nil {
			return fmt.Errorf("bump custom file %s error: %w", bmp.File.String(), err)
		}

		content = c
		changed = ch
	}

	if len(content) == 0 {
		return fmt.Errorf("file %s is not supported", bmp.File.String())
	}

	if !changed {
		console.Warn(fmt.Sprintf("File %s is not changed", bmp.File.String()))

		return nil
	}

	if !viper.GetBool(config.DryRun) {
		if err := write(f, bmp.File.Path(), content); err != nil {
			return fmt.Errorf("write file %s error: %w", bmp.File.String(), err)
		}
	}

	console.Info(fmt.Sprintf("Bump file %s", bmp.File.String()))

	return nil
}

func bumpPredefinedJSON(f file.Reader, patch string, v version.V) (_ [][]byte, changed bool, err error) {
	r, err := f.Read(patch)
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
	var content [][]byte

	replacer := convert.S2B(`"version": "` + v.FormatString() + `"`)

	for scanner.Scan() {
		line := scanner.Bytes()

		if !changed && versionRegex.Match(line) {
			line = versionRegex.ReplaceAll(line, replacer)
			changed = true
		}

		content = append(content, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, changed, fmt.Errorf("scan file %s error: %w", patch, err)
	}

	return content, changed, nil
}

func bumpCustomFile(f file.Reader, bmp config.BumpFile, v version.V) (_ [][]byte, changed bool, err error) {
	r, err := f.Read(bmp.File.Path())
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
	var content [][]byte

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
		line := scanner.Bytes()

		if lineNum >= start && lineNum <= end {
			if versionRegex.Match(line) {
				if len(regArr) > 0 {
					for _, r := range regArr {
						if r.Match(line) {
							line = versionRegex.ReplaceAll(line, convert.S2B(v.FormatString()))
							changed = true

							break
						}
					}
				} else {
					line = versionRegex.ReplaceAll(line, convert.S2B(v.FormatString()))
					changed = true
				}
			}
		}

		content = append(content, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, changed, fmt.Errorf("scan file %s error: %w", bmp.File.String(), err)
	}

	return content, changed, nil
}

func write(f file.Writer, patch string, content [][]byte) (err error) {
	w, err := f.Write(patch, os.O_WRONLY|os.O_TRUNC)
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

	n := convert.S2B("\n")

	for _, line := range content {
		_, err := w.Write(append(line, n...))
		if err != nil {
			return fmt.Errorf("write file %s line %s error: %w", patch, convert.B2S(line), err)
		}
	}

	return nil
}
