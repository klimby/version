package changelog

import (
	"testing"
	"time"

	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/git"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_commitTpl_shortHash(t *testing.T) {
	type fields struct {
		Hash string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "short hash",
			fields: fields{
				Hash: "0123456789",
			},
			want: "0123456",
		},
		{
			name: "short hash empty",
			fields: fields{
				Hash: "11",
			},
			want: "11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := commitTpl{
				Hash: tt.fields.Hash,
			}
			assert.Equalf(t, tt.want, m.shortHash(), "shortHash()")
		})
	}
}

func Test_newCommitTpl(t *testing.T) {
	tests := []struct {
		name     string
		gc       git.Commit
		showBody bool
		want     commitTpl
	}{
		{
			name: "new commit tpl",
			gc: git.Commit{
				Message: `feat(scope): message

body 1
body 2

close #123
`,
				Hash:   "0123456789",
				Author: "author",
				Date:   time.Now(),
				Email:  "foo@bar.com",
			},
			showBody: true,
			want: commitTpl{
				CommitType:       "feat",
				Scope:            "scope",
				Message:          "message",
				Body:             []string{"body 1", "body 2", "close #123"},
				isBreakingChange: false,
				Hash:             "0123456789",
				Author:           "author",
				AuthorHref:       "mailto:" + "foo@bar.com",
			},
		},
		{
			name: "new commit tpl no match",
			gc: git.Commit{
				Message: "message",
			},
			want: commitTpl{
				CommitType: "chore",
				Message:    "message",
			},
		},
		{
			name: "new commit tpl breaking change 1",
			gc: git.Commit{
				Message: `feat(scope): message

BREAKING CHANGE: breaking change
`,
			},
			showBody: true,
			want: commitTpl{
				CommitType:       "feat",
				Scope:            "scope",
				Message:          "message",
				Body:             []string{"BREAKING CHANGE: breaking change"},
				isBreakingChange: true,
			},
		},
		{
			name: "new commit tpl breaking change 2",
			gc: git.Commit{
				Message: "feat(scope)!: message",
			},
			want: commitTpl{
				CommitType:       "feat",
				Scope:            "scope",
				Message:          "message",
				isBreakingChange: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.ChangelogShowBody, tt.showBody)

			got := newCommitTpl(tt.gc)

			if tt.want.CommitType != "" {
				assert.Equal(t, tt.want.CommitType, got.CommitType, "CommitType")
			}

			if tt.want.Scope != "" {
				assert.Equal(t, tt.want.Scope, got.Scope, "Scope")
			}

			if tt.want.Message != "" {
				assert.Equal(t, tt.want.Message, got.Message, "Message")
			}

			if len(tt.want.Body) > 0 {
				assert.Equal(t, tt.want.Body, got.Body, "Body")
			}

			assert.Equal(t, tt.want.isBreakingChange, got.isBreakingChange, "isBreakingChange")

			if tt.want.Hash != "" {
				assert.Equal(t, tt.want.Hash, got.Hash, "Hash")
			}

			if tt.want.Author != "" {
				assert.Equal(t, tt.want.Author, got.Author, "Author")
			}

			if tt.want.AuthorHref != "" {
				assert.Equal(t, tt.want.AuthorHref, got.AuthorHref, "AuthorHref")
			}

		})
	}
}
