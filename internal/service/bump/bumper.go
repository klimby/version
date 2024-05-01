package bump

import (
	"bufio"
	"fmt"
	"io"
	"regexp"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/pkg/version"
)

// process is a content processor.
type process struct{}

// PredefinedJSON process predefined JSON file.
func (process) PredefinedJSON(r io.Reader, bmp config.BumpFile, v version.V) (_ []string, changed bool, err error) {
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
		return nil, changed, fmt.Errorf("scan file %s error: %w", bmp.File.Path(), err)
	}

	return content, changed, nil
}

// CustomFile process custom file.
func (process) CustomFile(r io.Reader, bmp config.BumpFile, v version.V) (_ []string, changed bool, err error) {
	scanner := bufio.NewScanner(r)
	versionRegex := regexp.MustCompile(`\d+\.\d+\.\d+`)
	var content []string

	start, end, regArr := handleBumpFile(bmp)

	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()

		if lineNum >= start && lineNum <= end {
			if versionRegex.MatchString(line) {
				if len(regArr) > 0 {
					for i := range regArr {
						if regArr[i].MatchString(line) {
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

// handleBumpFile handles BumpFile.
// Returns start, end lines and slice of regexp.Regexp.
func handleBumpFile(bmp config.BumpFile) (start, end int, regs []regexp.Regexp) {
	start = 0
	end = int(^uint(0) >> 1)

	if bmp.HasPositions() {
		start = bmp.Start
		end = bmp.End
	}

	regs = make([]regexp.Regexp, 0, len(bmp.RegExp))

	for _, r := range bmp.RegExp {
		rgx, err := regexp.Compile(r)
		if err != nil {
			console.Warn(fmt.Sprintf("file %s regexp %s error: %s", bmp.File.String(), r, err))

			continue
		}

		regs = append(regs, *rgx)
	}

	return start, end, regs
}
