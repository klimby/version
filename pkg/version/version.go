package version

import (
	"cmp"
	"regexp"
	"strings"

	"github.com/klimby/version/pkg/convert"
)

var (
	// Version regexp.
	re = regexp.MustCompile(`^(?:[a-zA-Z-.]*?)?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)(?:\.(?P<patch>0|[1-9]\d*))?(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	// Number regexp.
	reNum = regexp.MustCompile(`\d+`)
)

const (
	_maxInt = int(^uint(0) >> 1) // max int value.
)

// V version.
// Example: 1.0.0-asd+dd
//
// [SEMANTIC VERSIONING] https://semver.org/
type V string

// String returns the string value of the version.
func (v V) String() string {
	return string(v)
}

// Empty returns true if the version is empty.
func (v V) Empty() bool {
	return v.String() == "" || v.Invalid()
}

// Equal returns true if the version is equal to the argument.
func (v V) Equal(o V) bool {
	return v.Compare(o) == 0
}

// FormatString returns the string value of the version in the specified format.
func (v V) FormatString() string {
	if v.Invalid() {
		return ""
	}

	major, minor, patch, prerelease, buildmetadata := v.semver()

	var b strings.Builder

	b.WriteString(convert.I2S(major))
	b.WriteString(".")
	b.WriteString(convert.I2S(minor))
	b.WriteString(".")
	b.WriteString(convert.I2S(patch))

	if prerelease != "" {
		b.WriteString("-")
		b.WriteString(prerelease)

		if buildmetadata != "" {
			b.WriteString("+")
			b.WriteString(buildmetadata)
		}
	}

	return b.String()
}

// Invalid returns true if the version is invalid.
func (v V) Invalid() bool {
	if v == "" {
		return true
	}

	return len(re.FindStringSubmatch(string(v))) == 0
}

// NextMajor returns the next major version.
func (v V) NextMajor() V {
	major, _, _, _, _ := v.semver()

	return V(convert.I2S(major+1) + ".0.0")
}

// NextMinor returns the next minor version.
func (v V) NextMinor() V {
	major, minor, _, _, _ := v.semver()

	return V(convert.I2S(major) + "." + convert.I2S(minor+1) + ".0")
}

// NextPatch returns the next patch version.
func (v V) NextPatch() V {
	major, minor, patch, _, _ := v.semver()

	return V(convert.I2S(major) + "." + convert.I2S(minor) + "." + convert.I2S(patch+1))
}

// Start returns the start version.
func (v V) Start() V {
	return V("0.0.0")
}

// GitVersion returns the git version (version, started from "v").
func (v V) GitVersion() string {
	return "v" + v.FormatString()
}

// semver returns the semver value of the version.
//
// Returns 0.0.0 if the version is invalid.
// Return values are: major, minor, patch, prerelease, buildmetadata.
//
// Example usage:
//
//	major, minor, patch, prerelease, buildmetadata := version.Semver()
//
//nolint:revive
func (v V) semver() (major, minor, patch int, prerelease, buildmetadata string) {
	matches := re.FindStringSubmatch(v.String())

	if len(matches) != 0 {
		major = convert.S2Int(matches[re.SubexpIndex("major")])
		minor = convert.S2Int(matches[re.SubexpIndex("minor")])
		patch = convert.S2Int(matches[re.SubexpIndex("patch")])
		prerelease = matches[re.SubexpIndex("prerelease")]
		buildmetadata = matches[re.SubexpIndex("buildmetadata")]

		return major, minor, patch, prerelease, buildmetadata
	}

	return 0, 0, 0, "", ""
}

// Compare with other version:
// @see https://semver.org/#spec-item-11
//
// Precedence refers to how versions are compared to each other when ordered.
//
// Precedence MUST be calculated by separating the version into major, minor, patch and pre-release identifiers
// in that order (Build metadata does not figure into precedence).
//
// Precedence is determined by the first difference when comparing each of these identifiers from left to right
// as follows: Major, minor, and patch versions are always compared numerically.
//
// Example: 1.0.0 < 2.0.0 < 2.1.0 < 2.1.1.
//
// When major, minor, and patch are equal, a pre-release version has lower precedence than a normal version:
//
// Example: 1.0.0-alpha < 1.0.0.
//
// Precedence for two pre-release versions with the same major, minor, and patch version
// MUST be determined by comparing each dot separated identifier from left to right until a difference is found
// as follows:
//
//   - Identifiers consisting of only digits are compared numerically.
//   - Identifiers with letters or hyphens are compared lexically in ASCII sort order.
//     -Numeric identifiers always have lower precedence than non-numeric identifiers.
//   - A larger set of pre-release fields has a higher precedence than a smaller set, if all the preceding identifiers are equal.
//
// Example: 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0.
//
// Return values are:
//   - 0: equal;
//   - 1: greater, then argument;
//   - -1: les, then argument
func (v V) Compare(o V) int {
	if v.String() == o.String() {
		return 0
	}

	if v.Invalid() && o.Invalid() {
		return cmp.Compare[string](v.String(), o.String())
	}

	vMajor, vMinor, vPatch, vPrerelease, vBuildmetadata := v.semver()
	oMajor, oMinor, oPatch, oPrerelease, oBuildmetadata := o.semver()

	// compare major
	c := cmp.Compare[int](vMajor, oMajor)
	if c != 0 {
		return c
	}

	// compare minor
	c = cmp.Compare[int](vMinor, oMinor)
	if c != 0 {
		return c
	}

	// compare patch
	c = cmp.Compare[int](vPatch, oPatch)
	if c != 0 {
		return c
	}

	// compare prerelease
	c = comparePart(vPrerelease, oPrerelease)
	if c != 0 {
		return c
	}

	// compare buildmetadata
	return comparePart(vBuildmetadata, oBuildmetadata)
}

// Version returns the version for interface compatibility.
func (v V) Version() V {
	return v
}

// HasVersion interface.
type HasVersion interface {
	Version() V
}

// CompareASC compares two versions (ASC) (less to greater).
// Use in [slices.SortFunc].
//
// Example usage:
//
//	slices.SortFunc(tags, version.CompareDESC[version.V])
func CompareASC[T HasVersion](v, o T) int {
	return v.Version().Compare(o.Version())
}

// CompareDESC compares two versions (DESC) (greater to less).
// Use in [slices.SortFunc].
//
// Example usage:
//
//	slices.SortFunc(tags, version.CompareDESC[version.V])
func CompareDESC[T HasVersion](v, o T) int {
	c := v.Version().Compare(o.Version())

	switch c {
	case 1:
		return -1
	case -1:
		return 1
	default:
		return 0
	}
}

// CompareStrings compares two strings (build or prerelease).
//
// Returns:
//   - 0: equal;
//   - 1: v > o;
//   - -1: v < o
func comparePart(v, o string) int {
	sepV := strings.Split(v, ".")
	sepO := strings.Split(o, ".")
	maxLen := len(sepV)

	if len(sepO) > maxLen {
		maxLen = len(sepO)
	}

	switch {
	case v == o:
		return 0
	case v == "" && o != "":
		return 1
	case v != "" && o == "":
		return -1
	default:
		for i := 0; i < maxLen; i++ {
			first := ""
			second := ""

			if i < len(sepV) {
				first = sepV[i]
			}

			if i < len(sepO) {
				second = sepO[i]
			}

			c := comparePartElement(first, second)

			if c != 0 {
				return c
			}
		}

		return 0
	}
}

// comparePartElement compares two strings (part of build or prerelease).
func comparePartElement(v, o string) int {
	vIsNum := reNum.MatchString(v)
	oIsNum := reNum.MatchString(o)

	switch {
	case v == o:
		return 0
	case v == "" && o != "":
		return -1
	case v != "" && o == "":
		return 1
	case vIsNum && !oIsNum:
		return -1
	case !vIsNum && oIsNum:
		return 1
	case vIsNum && oIsNum:
		return cmp.Compare[int](convert.S2Int(v), convert.S2Int(o))
	default:
		return cmp.Compare[string](v, o)
	}
}
