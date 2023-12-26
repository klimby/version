package git

import (
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/klimby/version/pkg/version"
)

type nextType int

const (
	nextMajor nextType = iota + 1
	nextMinor
	nextPatch
)

type Repository struct {
	repo *git.Repository
}

type RepoOptions struct {
	Path string
}

// NewRepository returns a new Repository.
func NewRepository(opts ...func(options *RepoOptions)) (*Repository, error) {
	ro := &RepoOptions{
		Path: "",
	}

	for _, opt := range opts {
		opt(ro)
	}

	r, err := git.PlainOpen(ro.Path)
	if err != nil {
		return nil, err
	}

	return &Repository{
		repo: r,
	}, nil
}

// Repo returns a repository.
func (r Repository) GitRepo() *git.Repository {
	return r.repo
}

func (r Repository) LastVersion() (version.V, error) {
	lastTag, err := r.lastTag()
	if err != nil {
		return "", err
	}

	v := tagVersion(lastTag)

	return v, nil
}

// IsClean returns true if all the files are in Unmodified status.
func (r Repository) IsClean() (bool, error) {
	w, err := r.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("get worktree error: %s", err)
	}

	st, err := w.Status()
	if err != nil {
		return false, fmt.Errorf("get status error: %s", err)
	}

	return st.IsClean(), nil
}

// RemoteURL returns a repository name.
func (r Repository) RemoteURL() (string, error) {
	rem, err := r.repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("get remotes error: %s", err)
	}

	reg := regexp.MustCompile(`^.+(github.com.+).git$`)

	for _, rm := range rem {
		matches := reg.FindStringSubmatch(rm.Config().URLs[0])

		if len(matches) == 2 {
			return "https://" + matches[1], nil
		}
		// if contains
	}

	return "", nil
}

func (r Repository) NextMajor() (version.V, bool, error) {
	return r.nextVersion(nextMajor)
}

func (r Repository) NextMinor() (version.V, bool, error) {
	return r.nextVersion(nextMinor)
}

func (r Repository) NextPatch() (version.V, bool, error) {
	return r.nextVersion(nextPatch)
}

func (r Repository) nextVersion(nt nextType) (_ version.V, exists bool, _ error) {
	last, err := r.LastVersion()
	if err != nil {
		return "", exists, err
	}

	if last == "" {
		last = last.Start()
	}

	switch nt {
	case nextMajor:
		last = last.NextMajor()
	case nextMinor:
		last = last.NextMinor()
	case nextPatch:
		last = last.NextPatch()
	}

	for {
		tag := last.GitVersion()

		exists, err := r.tagExists(tag)
		if err != nil {
			return "", exists, err
		}

		if !exists {
			break
		}

		exists = true
		last = last.NextPatch()
	}

	return last, exists, nil
}

// CommitTag stores a tag.
func (r Repository) CommitTag(v version.V) (*object.Commit, error) {
	w, err := r.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("get worktree error: %s", err)
	}

	commit, err := w.Commit(fmt.Sprintf("chore(release): %s", v.FormatString()), &git.CommitOptions{})

	if _, err = r.repo.CreateTag(v.GitVersion(), commit, &git.CreateTagOptions{
		Message: fmt.Sprintf("chore(release): %s", v.FormatString()),
	}); err != nil {
		return nil, fmt.Errorf("create tag error: %s", err)
	}

	obj, err := r.repo.CommitObject(commit)
	if err != nil {
		return nil, fmt.Errorf("get commit object error: %s", err)
	}

	return obj, nil
}

// LastCommits returns a commits from last tag to HEAD.
func (r Repository) LastCommits() ([]object.Commit, error) {
	lastTag, err := r.lastTag()
	if err != nil {
		return nil, err
	}

	commits, err := r.repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	var cs []object.Commit

	for {
		c, err := commits.Next()
		if err != nil {
			break
		}

		if lastTag != nil && c.Hash == lastTag.Target {
			break
		}

		cs = append(cs, *c)
	}

	return cs, nil
}

func (r Repository) tagExists(tag string) (bool, error) {
	tags, err := r.repo.TagObjects()
	if err != nil {
		return false, err
	}

	for {
		t, err := tags.Next()
		if err != nil {
			break
		}
		if t.Name == tag {
			return true, nil
		}
	}

	return false, nil
}

// lastTag returns a last tag.
func (r Repository) lastTag() (*object.Tag, error) {
	var lastTag *object.Tag

	tagRefs, err := r.repo.Tags()
	if err != nil {
		return nil, err
	}

	h := plumbing.ZeroHash

	for {
		tagRef, err := tagRefs.Next()
		if err != nil {
			break
		}
		if tagRef == nil {
			break
		}

		h = tagRef.Hash()
	}

	if h == plumbing.ZeroHash {
		return lastTag, nil
	}

	lastTag, err = r.repo.TagObject(h)
	if err != nil {
		return nil, err
	}

	return lastTag, nil
}

// tagVersion returns a tag version.
func tagVersion(tag *object.Tag) version.V {
	if tag == nil {
		return ""
	}

	v := version.V(tag.Name)

	if v.Invalid() {
		return ""
	}

	return v
}
