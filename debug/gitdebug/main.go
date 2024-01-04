package main

import (
	"fmt"
	"path/filepath"

	"github.com/klimby/version/internal/changelog"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/internal/git"
	"github.com/spf13/viper"
)

func main() {
	p := filepath.Join("/home", "klim", "Projects", "version-test")

	config.Init(func(options *config.ConfigOptions) {
		options.WorkDir = p
		options.ChangelogFileName = "CHANGELOG_2.md"
		//	options.DryRun = true
	})

	repo, err := git.NewRepository(func(options *git.RepoOptions) {
		options.Path = viper.GetString(config.WorkDir)
	})
	if err != nil {
		fmt.Println(err)

		return
	}

	f := file.NewFS()

	remote, _ := repo.RemoteURL()
	config.SetUrlFromGit(remote)

	gen := changelog.NewGenerator(f, repo)

	/*if err := gen.Add(version.V("0.0.3")); err != nil {
		fmt.Println(err)

		return
	}*/

	if err := gen.Generate(); err != nil {
		fmt.Println(err)

		return
	}

	//	repo := r.GitRepo()

}

type CommitMessage struct {
}
