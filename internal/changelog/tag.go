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
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

type TagTpl struct {
	tag             version.V
	prev            version.V
	d               time.Time
	BreakingChanges []CommitTpl
	Blocks          []TagTplBlock
}

type TagTplBlock struct {
	// CommitType is a commit type (feat, fix, etc.).
	CommitType string
	// Name is a commit name.
	Name string
	// Commits is a list of commits.
	Commits []CommitTpl
}

type gitCommit interface {
	Message() string
	Hash() string
	IsTag() bool
	Version() version.V
	Date() time.Time
}

func NewTagTpl(tag version.V, date time.Time) TagTpl {
	return TagTpl{
		tag:             tag,
		d:               date,
		BreakingChanges: []CommitTpl{},
		Blocks:          bloks(),
	}
}

// Date returns a tag date as a string YYYY-MM-DD.
func (t *TagTpl) Date() string {
	return t.d.Format("2006-01-02")
}

// URL returns a tag URL.
func (t *TagTpl) URL() string {
	remoteURL := viper.GetString(config.RemoteURL)
	if remoteURL == "" || t.tag.Invalid() || t.prev.Invalid() {
		return ""
	}

	return fmt.Sprintf("%s/compare/%s...%s", remoteURL, t.prev.GitVersion(), t.tag.GitVersion())
}

// Name returns a tag name.
func (t *TagTpl) Name() string {
	return t.tag.FormatString()
}

// SetPrev sets the previous tag.
func (t *TagTpl) SetPrev(prev version.V) {
	t.prev = prev
}

// AddCommit adds a commit to the tag.
func (t *TagTpl) AddCommit(c gitCommit) {
	if c.IsTag() {
		return
	}

	tpl := NewCommitTpl(c.Message(), c.Hash())

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

func bloks() []TagTplBlock {
	nms := config.CommitNames()

	blocks := make([]TagTplBlock, len(nms))

	for i, nm := range nms {
		blocks[i] = TagTplBlock{
			CommitType: nm.Type,
			Name:       nm.Name,
			Commits:    []CommitTpl{},
		}
	}

	return blocks
}

// ApplyTemplate applies the template to the commit message.
func (t *TagTpl) ApplyTemplate(tplType TemplateType, wr io.Writer) error {
	switch tplType {
	case MarkdownTpl:
		return t.apply(wr, tplType, _tagMarkdownTpl)
	case ConsoleTpl:
		return t.apply(wr, tplType, _tagConsoleTpl)
	default:
		return fmt.Errorf("unknown template type: %d", tplType)
	}
}

func (t *TagTpl) apply(wr io.Writer, tplType TemplateType, tpl string) error {
	funcMap := template.FuncMap{
		"versionName": versionName(tplType),
		"commitName":  commitName(tplType),
		"addIssueURL": addIssueURL(tplType),
	}

	tmpl, err := template.New("tag").Funcs(funcMap).Parse(tpl)
	if err != nil {
		return fmt.Errorf("parse tag template error: %w", err)
	}

	if err := tmpl.Execute(wr, t); err != nil {
		return fmt.Errorf("execute tag template error: %w", err)
	}

	return nil
}

func versionName(tplType TemplateType) func(t TagTpl) string {
	remoteURL := viper.GetString(config.RemoteURL)

	return func(t TagTpl) string {
		if tplType == ConsoleTpl {
			return t.tag.FormatString()
		}

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

func commitName(tplType TemplateType) func(c CommitTpl) string {
	remoteURL := viper.GetString(config.RemoteURL)

	return func(c CommitTpl) string {

		var b strings.Builder

		if c.Scope != "" {
			b.WriteString("**")
			b.WriteString(c.Scope)
			b.WriteString(":**  ")
		}

		b.WriteString(c.Short)

		u, err := url.JoinPath(remoteURL, "commit", c.Hash)
		if err != nil {
			u = ""
		}

		if tplType == MarkdownTpl && u != "" {
			b.WriteString(" ([")
			b.WriteString(c.ShortHash)
			b.WriteString("](")
			b.WriteString(u)
			b.WriteString("))")

			return b.String()
		}

		b.WriteString(" (")
		b.WriteString(c.ShortHash)
		b.WriteString(")")

		return b.String()
	}
}

func addIssueURL(tplType TemplateType) func(s string) string {
	issueURL := viper.GetString(config.ChangelogIssueURL)
	re := regexp.MustCompile(`#\w+`)

	return func(s string) string {
		if issueURL == "" || tplType == ConsoleTpl {
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

func NewTagsTpl(commits []gitCommit) (TagsTpl, error) {
	var tags []TagTpl

	for _, c := range commits {
		if c.IsTag() {
			if len(tags) > 0 {
				tags[len(tags)-1].SetPrev(c.Version())
			}

			t := NewTagTpl(c.Version(), c.Date())

			tags = append(tags, t)

			continue
		}

		if len(tags) == 0 {
			return TagsTpl{}, fmt.Errorf("%w: last commit is not tagged (last commit must be version)", errGenerate)
		}

		tags[len(tags)-1].AddCommit(c)
	}

	return TagsTpl{
		Tags: tags,
	}, nil
}
