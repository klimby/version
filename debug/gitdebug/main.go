package main

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/klimby/version/internal/git"
)

func main() {

	reg := regexp.MustCompile(`#(?P<iss>[A-Za-z0-9-_]+)`)

	str := "Close 123"

	res := reg.ReplaceAllStringFunc(str, func(match string) string {
		word := match[1:] // Удаление символа '#'
		return fmt.Sprintf("[%s](https://foo/%s)", word, word)
	})

	fmt.Println(res)
	return

	/*messages := []string{
			`init

	BREAKING CHANGE: test

	Closes: #123
	`,
			`feat(SCOPES): short message

	long message 01
	long message 02

	BREAKING CHANGE: test

	Closes: #123
	`,
		}

		cm := changelog.NewCommitTpl(messages[1], func(options *changelog.CommitTplOptions) {
			options.RepoHref = "https://github.com/company/project/"
			options.IssueHref = "https://github.com/company/project/issues/"
		})

		var b strings.Builder

		if err := cm.ApplyTemplate(changelog.MarkdownTpl, &b); err != nil {
			fmt.Println(err)

			return
		}

		console.Info(b.String())

		return*/

	p := filepath.Join("/home", "klim", "Projects", "version-test")

	r, err := git.NewRepository(func(options *git.RepoOptions) {
		options.Path = p
	})
	if err != nil {
		fmt.Println(err)

		return
	}

	_ = r

	//	repo := r.GitRepo()

}

type CommitMessage struct {
}
