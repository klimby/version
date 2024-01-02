package changelog

import (
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

func TestDebug(t *testing.T) {

	//viper.Set(config.RemoteURL, "https://github.om/company/project")
	//viper.Set(config.ChangelogIssueURL, "https://github.om/company/project/issues/")

	config.Init()
	viper.Set(config.RemoteURL, "https://github.om/company/project")
	viper.Set(config.ChangelogIssueURL, "https://github.om/company/project/issues/")

	tagTpl := NewTagTpl(version.V("1.0.0"), time.Now())
	tagTpl.prev = version.V("0.0.1")

	c1 := __newFakeCommit("feat(SCOPES): short message", version.V("1.0.0"))
	tagTpl.AddCommit(c1)

	c2 := __newFakeCommit(`fix(SCOPES): short message

long message 01 #123

BREAKING CHANGE: test
`, version.V(""))

	tagTpl.AddCommit(c2)

	c3 := __newFakeCommit("feat(SCOPES): short message", version.V(""))
	tagTpl.AddCommit(c3)

	var b strings.Builder

	if err := tagTpl.ApplyTemplate(MarkdownTpl, &b); err != nil {
		t.Error(err)
	}

	res := b.String()
	console.Info(res)

}

type __fakeCommit struct {
	message string
	hash    string
	version version.V
	date    time.Time
}

func __newFakeCommit(message string, version version.V) __fakeCommit {
	return __fakeCommit{
		message: message,
		hash:    plumbing.ZeroHash.String(),
		version: version,
		date:    time.Now(),
	}
}

func (c __fakeCommit) Message() string {
	return c.message
}

func (c __fakeCommit) Hash() string {
	return c.hash
}

func (c __fakeCommit) IsTag() bool {
	return !c.version.Invalid()
}

func (c __fakeCommit) Version() version.V {
	return c.version
}

func (c __fakeCommit) Date() time.Time {
	return c.date
}
