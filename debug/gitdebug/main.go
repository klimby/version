package main

import (
	"fmt"
	"path/filepath"

	git2 "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/klimby/version/internal/git"
)

func main() {

	p := filepath.Join("/home", "klim", "Projects", "version-test")

	r, err := git.NewRepository(func(options *git.RepoOptions) {
		options.Path = p
	})
	if err != nil {
		fmt.Println(err)

		return
	}

	repo := r.GitRepo()
	tagRefs, err := repo.Tags()
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(r.RemoteURL())

	return

	tags := make([]plumbing.Hash, 0)

	for {
		tagRef, err := tagRefs.Next()
		if err != nil {
			break
		}

		if tagRef == nil {
			break
		}

		fmt.Println(tagRef.Name())
		fmt.Println(tagRef.Target())

		tags = append(tags, tagRef.Hash())
	}

	fmt.Println(tags)

	lastTag := tags[len(tags)-1]

	obj, err := repo.TagObject(lastTag)
	if err != nil {
		fmt.Println(err)

		return
	}

	_ = obj

	return

	tagObs, err := repo.TagObjects()
	if err != nil {
		fmt.Println(err)

		return
	}

	var lastOb *object.Tag

	for {
		tagOb, err := tagObs.Next()
		if err != nil {
			break
		}

		if tagOb == nil {
			break
		}

		if tagOb.Hash == lastTag {
			lastOb = tagOb

			break
		}

		fmt.Println(tagOb.Name)
	}

	fmt.Println(lastOb)

	fmt.Println("=========================================")

	if lastOb == nil {
		return
	}

	logs, err := repo.Log(&git2.LogOptions{
		//From: lastOb.Target,
	})
	if err != nil {
		fmt.Println(err)

		return
	}

	for {
		commit, err := logs.Next()
		if err != nil {
			break
		}

		if commit == nil {
			break
		}

		if commit.Hash == lastOb.Target {
			break
		}

		//fmt.Println()
		fmt.Println(commit.Message)
		fmt.Println(commit.Hash)

	}

	/*nextPatch, exists, err := r.NextPatch()
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(nextPatch, exists)

	if err := r.CommitTag(nextPatch); err != nil {
		fmt.Println(err)

		return
	}*/
}
