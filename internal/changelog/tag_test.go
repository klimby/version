package changelog

import (
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/git"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

func TestDebug(t *testing.T) {

	//viper.Set(config.RemoteURL, "https://github.om/company/project")
	//viper.Set(config.ChangelogIssueURL, "https://github.om/company/project/issues/")

	config.Init()
	viper.Set(config.RemoteURL, "https://github.om/company/project")
	viper.Set(config.ChangelogIssueURL, "https://github.om/company/project/issues/")
	viper.Set(config.ChangelogShowAuthor, true)

	tagTpl := newTagTpl(version.V("1.0.0"), time.Now())
	tagTpl.prev = version.V("0.0.1")

	breackComm := __newFakeCommit(func(c *git.Commit) {
		c.Message = `fix(SCOPES): short message

long message 01 #123

BREAKING CHANGE: test
`
		c.Author = "Author"
		c.Email = "asd@asd.com"
	})

	featComm01 := __newFakeCommit(func(c *git.Commit) {
		c.Message = `feat(SCOPES): short message

long message 02 #123`
	})

	featComm02 := __newFakeCommit(func(c *git.Commit) {
		c.Message = `feat(SCOPES): short message`
	})

	tagTpl.addCommit(breackComm)
	tagTpl.addCommit(featComm01)
	tagTpl.addCommit(featComm02)
	tagTpl.addCommit(featComm02)

	//c3 := __newFakeCommit("feat(SCOPES): short message", version.Version(""))
	//tagTpl.AddCommit(c3)
	//tagTpl.AddCommit(c3)

	var b strings.Builder

	if err := tagTpl.applyTemplate(&b); err != nil {
		t.Error(err)
	}

	tagTpl2 := newTagTpl(version.V("0.9.0"), time.Now())
	tagTpl2.prev = version.V("0.0.1")

	/*if err := tagTpl2.applyTemplate(MarkdownTpl, &b); err != nil {
		t.Error(err)
	}*/

	res := b.String()
	console.Info(res)

}

func __newFakeCommit(opt ...func(c *git.Commit)) git.Commit {
	cm := git.Commit{
		Hash: plumbing.ZeroHash.String(),
		Date: time.Now(),
	}

	for _, o := range opt {
		o(&cm)
	}

	return cm
}
