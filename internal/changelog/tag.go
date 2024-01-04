package changelog

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/git"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

type TagTpl struct {
	tag             version.V
	prev            version.V
	Date            string
	BreakingChanges []commitTpl
	Blocks          []TagTplBlock
}

type TagTplBlock struct {
	// CommitType is a commit type (feat, fix, etc.).
	CommitType string
	// Name is a commit name.
	Name string
	// Commits is a list of commits.
	Commits []commitTpl
}

func NewTagTpl(tag version.V, date time.Time) TagTpl {
	nms := config.CommitNames()

	blocks := make([]TagTplBlock, len(nms))

	for i, nm := range nms {
		blocks[i] = TagTplBlock{
			CommitType: nm.Type,
			Name:       nm.Name,
			Commits:    []commitTpl{},
		}
	}

	return TagTpl{
		tag:             tag,
		Date:            date.Format("2006-01-02"),
		BreakingChanges: []commitTpl{},
		Blocks:          blocks,
	}
}

// setPrev sets the previous tag.
func (t *TagTpl) setPrev(prev version.V) {
	t.prev = prev
}

// addCommit adds a commit to the tag.
func (t *TagTpl) addCommit(c git.Commit) {
	if c.IsTag() {
		return
	}

	tpl := newCommitTpl(c)

	if tpl.isBreakingChange {
		t.BreakingChanges = append(t.BreakingChanges, tpl)
		return
	}

	for i, b := range t.Blocks {
		if b.CommitType == tpl.CommitType {
			t.Blocks[i].Commits = append(t.Blocks[i].Commits, tpl)
			return
		}
	}
}

// applyTemplate applies the template to the commit message.
func (t *TagTpl) applyTemplate(wr io.Writer) error {
	funcMap := template.FuncMap{
		"versionName": versionName(),
		"commitName":  commitName(),
		"addIssueURL": addIssueURL(),
	}

	tmpl, err := template.New("tag").Funcs(funcMap).Parse(_tagMarkdownTpl)
	if err != nil {
		return fmt.Errorf("parse tag template error: %w", err)
	}

	if err := tmpl.Execute(wr, t); err != nil {
		return fmt.Errorf("execute tag template error: %w", err)
	}

	return nil
}

// versionName returns a version name string in template.
func versionName() func(t TagTpl) string {
	remoteURL := viper.GetString(config.RemoteURL)

	return func(t TagTpl) string {
		if remoteURL == "" || t.tag.Invalid() || t.prev.Invalid() {
			return t.tag.FormatString()
		}

		u, err := url.JoinPath(remoteURL, "compare", fmt.Sprintf("%s...%s", t.prev.GitVersion(), t.tag.GitVersion()))
		if err != nil {
			return t.tag.FormatString()
		}

		return fmt.Sprintf("[%s](%s)", t.tag.FormatString(), u)
	}
}

// commitName returns a commit name string in template.
func commitName() func(c commitTpl) string {
	remoteURL := viper.GetString(config.RemoteURL)
	showAuthor := viper.GetBool(config.ChangelogShowAuthor)

	return func(c commitTpl) string {

		var b strings.Builder

		if c.Scope != "" {
			b.WriteString("**")
			b.WriteString(c.Scope)
			b.WriteString(":** ")
		}

		b.WriteString(c.Message)

		u, err := url.JoinPath(remoteURL, "commit", c.Hash)
		if err != nil || remoteURL == "" {
			u = ""
		}

		if u != "" {
			b.WriteString(" ([" + c.shortHash() + "](" + u + "))")
		} else {
			b.WriteString(" (" + c.shortHash() + ")")
		}

		if showAuthor && c.Author != "" {
			b.WriteString(" - ")
			if c.AuthorHref != "" {
				b.WriteString(fmt.Sprintf("[%s](%s)", c.Author, c.AuthorHref))
			} else {
				b.WriteString(c.Author)
			}
		}

		return b.String()
	}
}

// addIssueURL returns a commit message with issue URL in template.
func addIssueURL() func(s string) string {
	issueURL := viper.GetString(config.ChangelogIssueURL)
	re := regexp.MustCompile(`#\w+`)

	return func(s string) string {
		if issueURL == "" {
			return s
		}

		return re.ReplaceAllStringFunc(s, func(match string) string {
			i := match[1:] // Удаление символа '#'

			u, err := url.JoinPath(issueURL, i)
			if err != nil {
				return s
			}

			return fmt.Sprintf("[%s](%s)", i, u)
		})

	}
}

type TagsTpl struct {
	Tags []TagTpl
}

func NewTagsTpl(commits []git.Commit) (TagsTpl, error) {
	var tags []TagTpl

	for _, c := range commits {
		if c.IsTag() {
			if len(tags) > 0 {
				tags[len(tags)-1].setPrev(c.Version)
			}

			t := NewTagTpl(c.Version, c.Date)

			tags = append(tags, t)

			continue
		}

		if len(tags) == 0 {
			console.Info(fmt.Sprintf("commit %s is not included in version, skip", c.Hash))

			continue
		}

		tags[len(tags)-1].addCommit(c)
	}

	return TagsTpl{
		Tags: tags,
	}, nil
}

// applyTemplate applies the template to the commit message.
func (t TagsTpl) applyTemplate(wr io.Writer) error {
	for _, t := range t.Tags {
		if err := t.applyTemplate(wr); err != nil {
			return err
		}
	}

	return nil
}
