package changelog

import (
	"testing"
	"time"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_addIssueURL(t *testing.T) {
	tests := []struct {
		name     string
		issueURL string
		s        string
		want     string
	}{
		{
			name:     "add issue url",
			issueURL: "https://example.com/issues/",
			s:        "fix: issue #123",
			want:     "fix: issue [123](https://example.com/issues/123)",
		},
		{
			name:     "empty issue url",
			issueURL: "",
			s:        "fix: issue #123",
			want:     "fix: issue #123",
		},
		{
			name:     "invalid issue url",
			issueURL: "https://example.com:issues/",
			s:        "fix: issue #123",
			want:     "fix: issue #123",
		},
		{
			name:     "not issue",
			issueURL: "https://example.com/issues/",
			s:        "fix: issue 123",
			want:     "fix: issue 123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.ChangelogIssueURL, tt.issueURL)

			assert.Equal(t, tt.want, addIssueURL()(tt.s))
		})
	}
}

func Test_commitName(t *testing.T) {
	tests := []struct {
		name                string
		c                   commitTpl
		remoteUrl           string
		changelogShowAuthor bool
		want                string
	}{
		{
			name:                "commit",
			remoteUrl:           "https://example.com",
			changelogShowAuthor: true,
			c: commitTpl{
				Scope:      "scope",
				Message:    "message",
				Hash:       "0123456789",
				Author:     "author",
				AuthorHref: "https://example.com/author",
			},
			want: "**scope:** message ([0123456](https://example.com/commit/0123456789)) - [author](https://example.com/author)",
		},
		{
			name:                "commit no remote url",
			changelogShowAuthor: true,
			c: commitTpl{
				Scope:      "scope",
				Message:    "message",
				Hash:       "0123456789",
				Author:     "author",
				AuthorHref: "https://example.com/author",
			},
			want: "**scope:** message (0123456) - [author](https://example.com/author)",
		},
		{
			name:                "commit no author href",
			remoteUrl:           "https://example.com",
			changelogShowAuthor: true,
			c: commitTpl{
				Scope:   "scope",
				Message: "message",
				Hash:    "0123456789",
				Author:  "author",
			},
			want: "**scope:** message ([0123456](https://example.com/commit/0123456789)) - author",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.ChangelogShowAuthor, tt.changelogShowAuthor)
			viper.Set(key.RemoteURL, tt.remoteUrl)

			assert.Equal(t, tt.want, commitName()(tt.c))
		})
	}
}

func Test_versionName(t *testing.T) {
	tests := []struct {
		name      string
		remoteUrl string
		t         tagTpl
		want      string
	}{
		{
			name:      "version",
			remoteUrl: "https://example.com",
			t: tagTpl{
				tag:  "v1.0.0",
				prev: "v0.1.0",
			},
			want: "[1.0.0](https://example.com/compare/v0.1.0...v1.0.0)",
		},
		{
			name: "empty remote url",
			t: tagTpl{
				tag:  "v1.0.0",
				prev: "v0.1.0",
			},
			want: "1.0.0",
		},
		{
			name:      "invalid remote url",
			remoteUrl: "https://example.com:compare",
			t: tagTpl{
				tag:  "v1.0.0",
				prev: "v0.1.0",
			},
			want: "1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.RemoteURL, tt.remoteUrl)

			assert.Equal(t, tt.want, versionName()(tt.t))
		})
	}
}

func Test_tagTpl_setPrev(t1 *testing.T) {
	type args struct {
		prev version.V
	}
	tests := []struct {
		name   string
		tagTpl tagTpl
		args   args
	}{
		{
			name: "set prev",
			tagTpl: newTagTpl([]config.CommitName{
				{Type: "feat", Name: "Features"},
			}, "v1.0.0", time.Now()),
			args: args{
				prev: "v0.1.0",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.tagTpl.setPrev(tt.args.prev)

			assert.Equal(t1, tt.args.prev, tt.tagTpl.prev)
		})
	}
}
