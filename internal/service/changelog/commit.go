package changelog

import (
	"regexp"
	"strings"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/git"
	"github.com/klimby/version/pkg/convert"
	"github.com/spf13/viper"
)

var (
	_titleRegexp          = regexp.MustCompile(`^\s*(?P<tpe>[A-Za-z-_]*)(?:\((?P<scp>.+)\))?(?P<bre>!)?:\s*(?P<msg>.+)\s*$`)
	_breakingChangeRegexp = regexp.MustCompile(`^BREAKING[ -]CHANGE:.+$`)
)

// commitTpl is a commit message.
type commitTpl struct {
	// source is an original commit message.
	source string
	// CommitType is a commit type (feat, fix, etc.).
	CommitType string
	// Scope is a commit Scope.
	Scope string
	// Message is a commit message.
	Message string
	// Commit body.
	Body []string
	// isBreakingChange is a breaking change flag (existed "!" in the title or "BREAKING CHANGE:" in the body).
	isBreakingChange bool
	// Hash is a commit hash.
	Hash string
	// Author is a commit author.
	Author string
	// AuthorHref is a commit author href.
	AuthorHref string
}

// shortHash returns the short commit hash.
func (m commitTpl) shortHash() string {
	return m.Hash[:7]
}

// newCommitTpl returns a new commitTpl.
func newCommitTpl(gc git.Commit) commitTpl {
	m := commitTpl{
		source:     gc.Message,
		Hash:       gc.Hash,
		Author:     gc.Author,
		AuthorHref: gc.AuthorHref(),
	}

	spl := strings.Split(m.source, "\n")

	showBody := viper.GetBool(key.ChangelogShowBody)

	matches := _titleRegexp.FindStringSubmatch(spl[0])

	if len(matches) == 0 {
		m.Message = convert.S2Clear(spl[0])
		m.CommitType = config.CommitChore
	} else {
		m.CommitType = matches[_titleRegexp.SubexpIndex("tpe")]
		m.Scope = matches[_titleRegexp.SubexpIndex("scp")]
		m.Message = matches[_titleRegexp.SubexpIndex("msg")]

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

			if showBody {
				m.Body = append(m.Body, ll)
			}
		}
	}

	return m
}
