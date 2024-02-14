package changelog

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/git"
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

func Test_tagTpl_addCommit(t1 *testing.T) {
	type fields struct {
		BreakingChanges []commitTpl
		Blocks          []tagTplBlock
	}
	type args struct {
		c git.Commit
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		wantBreakingChangesLen int
		wantBlockCommitsLen    int
	}{
		{
			name: "add BreakingChanges commit",
			fields: fields{
				BreakingChanges: []commitTpl{},
				Blocks: []tagTplBlock{
					{
						CommitType: "feat",
						Name:       "Features",
						Commits:    []commitTpl{},
					},
				},
			},
			args: args{
				c: git.Commit{
					Message: "feat(cmd): message\n\nBREAKING CHANGE: foo",
				},
			},
			wantBreakingChangesLen: 1,
			wantBlockCommitsLen:    0,
		},
		{
			name: "add feat commit",
			fields: fields{
				BreakingChanges: []commitTpl{},
				Blocks: []tagTplBlock{
					{
						CommitType: "feat",
						Name:       "Features",
						Commits:    []commitTpl{},
					},
				},
			},
			args: args{
				c: git.Commit{
					Message: "feat(cmd): message",
				},
			},
			wantBreakingChangesLen: 0,
			wantBlockCommitsLen:    1,
		},
		{
			name: "add tag commit",
			fields: fields{
				BreakingChanges: []commitTpl{},
				Blocks: []tagTplBlock{
					{
						CommitType: "feat",
						Name:       "Features",
						Commits:    []commitTpl{},
					},
				},
			},
			args: args{
				c: git.Commit{
					Message: "feat(cmd): message",
					Version: "v1.0.0",
				},
			},
			wantBreakingChangesLen: 0,
			wantBlockCommitsLen:    0,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tagTpl{
				tag:             version.V("1.0.0"),
				prev:            version.V("0.1.0"),
				Date:            "2021-01-01",
				BreakingChanges: tt.fields.BreakingChanges,
				Blocks:          tt.fields.Blocks,
			}
			t.addCommit(tt.args.c)

			assert.Len(t1, t.BreakingChanges, tt.wantBreakingChangesLen, "BreakingChanges length")

			if len(t.Blocks) > 0 {
				assert.Len(t1, t.Blocks[0].Commits, tt.wantBlockCommitsLen, "Blocks commits length")
			}
		})
	}
}

func Test_tagTpl_applyTemplate(t1 *testing.T) {
	type fields struct {
		tag             version.V
		prev            version.V
		Date            string
		BreakingChanges []commitTpl
		Blocks          []tagTplBlock
	}
	tests := []struct {
		name    string
		fields  fields
		wantWr  string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "apply template",
			fields: fields{
				tag:  "v1.0.0",
				prev: "v0.1.0",
				Date: "2021-01-01",
				BreakingChanges: []commitTpl{
					{
						Scope:   "scope",
						Message: "message",
						Hash:    "0123456789",
					},
				},
				Blocks: []tagTplBlock{
					{
						CommitType: "feat",
						Name:       "Features",
						Commits: []commitTpl{
							{
								Scope:   "scope",
								Message: "message",
								Hash:    "0123456789",
							},
						},
					},
				},
			},
			wantWr: `
## [1.0.0](https://example.com/compare/v0.1.0...v1.0.0) (2021-01-01)

### Breaking changes

* **scope:** message ([0123456](https://example.com/commit/0123456789))

### Features

* **scope:** message ([0123456](https://example.com/commit/0123456789))
`,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tagTpl{
				tag:             tt.fields.tag,
				prev:            tt.fields.prev,
				Date:            tt.fields.Date,
				BreakingChanges: tt.fields.BreakingChanges,
				Blocks:          tt.fields.Blocks,
			}

			viper.Set(key.RemoteURL, "https://example.com")
			viper.Set(key.ChangelogIssueURL, "https://example.com/issues/")

			wr := &bytes.Buffer{}

			err := t.applyTemplate(wr)

			if !tt.wantErr(t1, err, fmt.Sprintf("applyTemplate(%v)", wr)) {
				return
			}

			assert.Equalf(t1, tt.wantWr, wr.String(), "applyTemplate(%v)", wr)
		})
	}
}

func Test_newTagsTpl(t *testing.T) {
	nms := []config.CommitName{
		{Type: "feat", Name: "Features"},
		//{Type: "chore", Name: "Other changes"},
	}

	commits := []git.Commit{
		{
			Message: "feat: message 0",
			Date:    time.Now(),
		},
		{
			Message: "chore: message",
			Version: version.V("v1.0.0"),
			Date:    time.Now(),
		},
		{
			Message: "feat: message 1",
			Date:    time.Now(),
		},
		{
			Message: "chore: message",
			Version: version.V("v0.1.0"),
			Date:    time.Now(),
		},
		{
			Message: "feat: message 2",
			Date:    time.Now(),
		},
	}

	tTpl := newTagsTpl(nms, commits)

	assert.Len(t, tTpl.Tags, 2, "tags length")

	first := tTpl.Tags[0]
	assert.Equal(t, version.V("v1.0.0"), first.tag, "tag first")
	assert.Equal(t, version.V("v0.1.0"), first.prev, "tag first prev")
	assert.Len(t, first.Blocks, 1, "tag first Blocks length")
	assert.Len(t, first.Blocks[0].Commits, 1, "tag first Blocks commits length")
	assert.Equal(t, "message 1", first.Blocks[0].Commits[0].Message, "tag first Blocks commits message")

	second := tTpl.Tags[1]
	assert.Equal(t, version.V("v0.1.0"), second.tag, "tag second")
	assert.Equal(t, version.V(""), second.prev, "tag second prev")
	assert.Len(t, second.Blocks, 1, "tag second Blocks length")
	assert.Len(t, second.Blocks[0].Commits, 1, "tag second Blocks commits length")
	assert.Equal(t, "message 2", second.Blocks[0].Commits[0].Message, "tag second Blocks commits message")

}

func Test_tagsTpl_applyTemplate(t *testing.T) {
	nms := []config.CommitName{
		{Type: "feat", Name: "Features"},
	}

	commits := []git.Commit{
		{
			Message: "chore: message",
			Version: version.V("v1.0.0"),
			Date:    time.Now(),
		},
		{
			Message: "feat: message 1",
			Date:    time.Now(),
		},
	}

	tTpl := newTagsTpl(nms, commits)

	viper.Set(key.RemoteURL, "https://example.com")
	viper.Set(key.ChangelogIssueURL, "https://example.com/issues/")

	err := tTpl.applyTemplate(&bytes.Buffer{})

	assert.NoError(t, err, "applyTemplate")
}
