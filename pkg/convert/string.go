package convert

import (
	"regexp"
	"strconv"
	"strings"
)

// S2Int convert string to int.
func S2Int(s string, d ...int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		}

		return 0
	}

	return i
}

// I2S convert int to string.
func I2S(i int) string {
	return strconv.Itoa(i)
}

// S2Clear convert string to clear string.
func S2Clear(s string) string {
	replaced := strings.ReplaceAll(s, "\t", "")
	replaced = strings.ReplaceAll(replaced, "\u0001", " ") // This symbol is used in the git log output.
	replaced = strings.ReplaceAll(replaced, "\n", " ")
	reg := regexp.MustCompile(`\s+`)
	replaced = reg.ReplaceAllString(replaced, " ")

	return strings.TrimSpace(replaced)
}
