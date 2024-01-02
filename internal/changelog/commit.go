package changelog

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/pkg/convert"
	"github.com/spf13/viper"
)

var (
	_titleRegexp          = regexp.MustCompile(`^\s*(?P<tpe>[A-Za-z-_]*)(?:\((?P<scp>.+)\))?(?P<bre>!)?:\s*(?P<msg>.+)\s*$`)
	_breakingChangeRegexp = regexp.MustCompile(`^BREAKING[ -]CHANGE:.+$`)
)

// CommitTpl is a commit message.
type CommitTpl struct {
	// source is an original commit message.
	source string
	// CommitType is a commit type (feat, fix, etc.).
	CommitType string
	// Scope is a commit Scope.
	Scope string
	// Short is a Short commit message.
	Short string
	// Commit body.
	Body []string
	// isBreakingChange is a breaking change flag (existed "!" in the title or "BREAKING CHANGE:" in the body).
	isBreakingChange bool
	// ShortHash is a short commit hash.
	ShortHash string
	Hash      string
	// IssueURL is an issue URL.
	IssueURL string
	// URL is a commit URL.
	URL string
}

// NewCommitTpl returns a new CommitTpl.
func NewCommitTpl(s, hash string) CommitTpl {
	m := CommitTpl{
		source:    s,
		ShortHash: shortHash(hash),
		Hash:      hash,
		URL:       hashURL(hash),
	}

	spl := strings.Split(s, "\n")

	matches := _titleRegexp.FindStringSubmatch(spl[0])

	if len(matches) == 0 {
		m.Short = convert.S2Clear(spl[0])
		m.CommitType = config.CommitChore
	} else {
		m.CommitType = matches[_titleRegexp.SubexpIndex("tpe")]
		m.Scope = matches[_titleRegexp.SubexpIndex("scp")]
		m.Short = matches[_titleRegexp.SubexpIndex("msg")]

		if matches[_titleRegexp.SubexpIndex("bre")] == "!" {
			m.isBreakingChange = true
		}
	}

	if len(spl) > 1 {
		for _, l := range spl[1:] {
			ll := convert.S2Clear(l)
			if ll == "" {
				continue
			}

			if _breakingChangeRegexp.MatchString(ll) {
				m.isBreakingChange = true
			}

			m.Body = append(m.Body, ll)
		}
	}

	return m
}

// HashURL returns the commit hash URL.
func hashURL(hash string) string {
	ru := viper.GetString(config.RemoteURL)

	if ru == "" {
		return ""
	}

	u, err := url.JoinPath(ru, "commit", hash)
	if err != nil {
		return ""
	}

	return u
}

// IssueURL returns the issue URL.
func issueURL(issue string) string {
	ru := viper.GetString(config.ChangelogIssueURL)

	if ru == "" {
		return ""
	}

	u, err := url.JoinPath(ru, issue)
	if err != nil {
		return ""
	}

	return u
}

// ShortHash returns the short commit hash.
func shortHash(hash string) string {
	return hash[:7]
}
