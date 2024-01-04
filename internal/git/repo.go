package git

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/klimby/version/pkg/version"
)

type NextType int

const (
	NextNone  NextType = iota
	NextMajor NextType = iota
	NextMinor
	NextPatch
	NextCustom
)

type Repository struct {
	repo *git.Repository
}

type RepoOptions struct {
	Path string
	Repo *git.Repository
}

// NewRepository returns a new Repository.
func NewRepository(opts ...func(options *RepoOptions)) (*Repository, error) {
	ro := &RepoOptions{
		Path: "",
	}

	for _, opt := range opts {
		opt(ro)
	}

	if ro.Repo != nil {
		return &Repository{
			repo: ro.Repo,
		}, nil
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

/*func (r Repository) LastVersion() (version.Version, error) {
	lastTag, err := r.LastTag()
	if err != nil {
		return "", err
	}

	return lastTag.Version(), nil
}*/

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
	}

	return "", nil
}

func (r Repository) NextVersion(nt NextType, custom version.V) (_ version.V, exists bool, _ error) {
	if nt == NextNone {
		return "", false, nil
	}

	lastTag, err := r.LastTag()
	if err != nil {
		return "", exists, err
	}

	last := lastTag.Version()

	if last.Empty() {
		last = last.Start()
	}

	var next version.V

	switch nt {
	case NextMajor:
		next = last.NextMajor()
	case NextMinor:
		next = last.NextMinor()
	case NextPatch:
		next = last.NextPatch()
	case NextCustom:
		next = custom
	default:
		return "", exists, fmt.Errorf("unknown next type")
	}

	for {
		tag := next.GitVersion()

		exists, err := r.tagExists(tag)
		if err != nil {
			return "", exists, err
		}

		if !exists {
			break
		}

		exists = true
		next = next.NextPatch()
	}

	return next, exists, nil
}

// CheckDowngrade checks if the version is not downgraded.
func (r Repository) CheckDowngrade(v version.V) error {
	lastTag, err := r.LastTag()
	if err != nil {
		return err
	}

	last := lastTag.Version()

	if last.Empty() {
		last = last.Start()
	}

	if last.Version().Compare(v) == -1 {
		return fmt.Errorf("version downgrade: %s -> %s", last.Version().FormatString(), v.FormatString())
	}

	return nil
}

// CommitTag stores a tag.
func (r Repository) CommitTag(v version.V) (*Commit, error) {
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

	cmt := newCommitFromGit(*obj)

	cmt.Version = v

	return &cmt, nil
}

// Commits returns a commits from last tag to HEAD.
func (r Repository) Commits(v version.V) ([]Commit, error) {
	tags, err := r.Tags()
	if err != nil {
		return nil, err
	}

	commits, err := r.repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	var cs []Commit

	for {
		c, err := commits.Next()
		if err != nil {
			break
		}

		cmt := newCommitFromGit(*c)

		setTagToCommit(&cmt, tags)

		if !v.Empty() && v.Equal(cmt.Version) {
			break
		}

		cs = append(cs, cmt)
	}

	return cs, nil
}

// LastCommit returns a last commit.
func (r Repository) LastCommit() (*Commit, error) {
	tags, err := r.Tags()
	if err != nil {
		return nil, err
	}

	commits, err := r.repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	c, err := commits.Next()
	if err != nil {
		return nil, err
	}

	cmt := newCommitFromGit(*c)

	setTagToCommit(&cmt, tags)

	return &cmt, nil
}

func setTagToCommit(c *Commit, tags []Tag) {
	for _, t := range tags {
		if c.Hash == t.t.Target.String() {
			tagVer := t.Version()

			if !tagVer.Invalid() {
				c.Version = tagVer
			}

			break
		}
	}
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

// LastTag returns a last tag.
func (r Repository) LastTag() (*Tag, error) {
	tags, err := r.Tags()
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return nil, nil
	}

	return &tags[len(tags)-1], nil
}

// Tags returns a list of Tags.
func (r Repository) Tags() ([]Tag, error) {
	tagRefs, err := r.repo.Tags()
	if err != nil {
		return nil, err
	}

	var tags []Tag

	for {
		tagRef, err := tagRefs.Next()
		if err != nil {
			break
		}
		if tagRef == nil {
			break
		}

		tag, err := r.repo.TagObject(tagRef.Hash())
		if err != nil {
			return nil, err
		}

		tags = append(tags, Tag{t: *tag})
	}

	return tags, nil
}

type Commit struct {
	Hash    string
	Message string
	Author  string
	Version version.V
	Date    time.Time
	Email   string
}

// newCommitFromGit returns a new Commit.
func newCommitFromGit(c object.Commit) Commit {
	return Commit{
		Hash:    c.Hash.String(),
		Message: c.Message,
		Date:    c.Author.When,
		Author:  c.Author.Name,
		Email:   c.Author.Email,
	}
}

// IsTag returns true if the commit is tagged.
func (c Commit) IsTag() bool {
	return !c.Version.Invalid()
}

// AuthorHref returns a commit Author href.
func (c Commit) AuthorHref() string {
	if c.Author == "" {
		return ""
	}

	// if Author start with @, then it is a GitHub username
	// and return GitHub Author href
	if c.Author[0] == '@' {
		return "https://github.com/" + c.Author[1:]
	}

	if c.Email == "" {
		return ""
	}

	return fmt.Sprintf("mailto:%s", c.Email)
}

type Tag struct {
	t object.Tag
}

func (t Tag) Version() version.V {
	v := version.V(t.t.Name)

	if v.Invalid() {
		return ""
	}

	return v
}
