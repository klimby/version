package git

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

// NextType is a next version type.
type NextType int

// NextType values.
const (
	NextNone   NextType = iota // NextNone is a next version type none (invalid).
	NextMajor  NextType = iota // NextMajor is a next version type major.
	NextMinor                  // NextMinor is a next version type minor.
	NextPatch                  // NextPatch is a next version type patch.
	NextCustom                 // NextCustom is a next version type custom (need to set custom version).
)

var (
	errTagsNotFound = errors.New("tags not found")
)

// Repository is a git repository wrapper.
type Repository struct {
	repo *git.Repository
	path string
}

// RepoOptions is a Repository options.
type RepoOptions struct {
	Path string
	Repo *git.Repository
}

// NewRepository returns a new Repository.
func NewRepository(opts ...func(options *RepoOptions)) (*Repository, error) {
	ro := &RepoOptions{
		Path: viper.GetString(config.WorkDir),
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
		path: ro.Path,
	}, nil
}

// IsClean returns true if all the files are in Unmodified status.
func (r Repository) IsClean() (bool, error) {
	w, err := r.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("get worktree error: %w", err)
	}

	st, err := w.Status()
	if err != nil {
		return false, fmt.Errorf("get status error: %w", err)
	}

	if len(st) == 0 {
		return true, nil
	}

	for _, s := range st {
		if !(s.Staging == git.Untracked && s.Worktree == git.Untracked) {
			return false, nil
		}
	}

	return true, nil
}

// AddModified adds modified files to the index.
func (r Repository) AddModified() error {
	if viper.GetBool(config.DryRun) {
		return nil
	}

	w, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("get worktree error: %w", err)
	}

	st, err := w.Status()
	if err != nil {
		return fmt.Errorf("get status error: %w", err)
	}

	for path, s := range st {
		if s.Worktree == git.Modified || s.Staging == git.Added {
			if err := w.AddWithOptions(&git.AddOptions{
				Path: path,
			}); err != nil {
				return fmt.Errorf("add file %s error: %w", path, err)
			}
		}
	}

	return nil
}

// RemoteURL returns a repository name.
func (r Repository) RemoteURL() (string, error) {
	rem, err := r.repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("get remotes error: %w", err)
	}

	reg := regexp.MustCompile(`^.+(github\.com.+).git$`)

	for _, rm := range rem {
		matches := reg.FindStringSubmatch(rm.Config().URLs[0])

		if len(matches) == 2 {
			return "https://" + matches[1], nil
		}
	}

	return "", nil
}

// NextVersion returns a next version.
func (r Repository) NextVersion(nt NextType, custom version.V) (_ version.V, exists bool, _ error) {
	if nt == NextNone {
		return "", false, nil
	}

	var lastV version.V

	lastTag, err := r.lastTag()
	if err != nil {
		if !errors.Is(err, errTagsNotFound) {
			return "", false, err
		}

		lastV = lastV.Start()
	} else {
		lastV = lastTag.Version()
	}

	if lastV.Empty() {
		lastV = lastV.Start()
	}

	var next version.V

	switch nt {
	case NextMajor:
		next = lastV.NextMajor()
	case NextMinor:
		next = lastV.NextMinor()
	case NextPatch:
		next = lastV.NextPatch()
	case NextCustom:
		next = custom
	case NextNone:
		return "", false, fmt.Errorf("unknown next type")
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

		next = next.NextPatch()
	}

	return next, exists, nil
}

// CheckDowngrade checks if the version is not downgraded.
func (r Repository) CheckDowngrade(v version.V) error {
	lastTag, err := r.lastTag()
	if err != nil {
		if errors.Is(err, errTagsNotFound) {
			return nil
		}

		return err
	}

	last := lastTag.Version()

	if last.Empty() {
		last = last.Start()
	}

	if v.LessThen(last) {
		return fmt.Errorf("version downgrade: %s -> %s", last.Version().FormatString(), v.FormatString())
	}

	return nil
}

// Add files to the index.
// files is list from path to files FROM WORKDIR.
func (r Repository) Add(files ...config.File) error {
	if viper.GetBool(config.DryRun) {
		return nil
	}

	w, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("get worktree error: %w", err)
	}

	for _, f := range files {
		if err := w.AddWithOptions(&git.AddOptions{
			Path: f.Rel(),
		}); err != nil {
			return fmt.Errorf("add file %s error: %w", f.Rel(), err)
		}
	}

	return nil
}

// CommitTag stores a tag and commit changes.
func (r Repository) CommitTag(v version.V) (*Commit, error) {
	if viper.GetBool(config.DryRun) {
		return nil, nil
	}

	w, err := r.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("get worktree error: %w", err)
	}

	commit, err := w.Commit(fmt.Sprintf("chore(release): %s", v.FormatString()), &git.CommitOptions{})
	if err != nil {
		return nil, fmt.Errorf("commit error: %w", err)
	}

	if _, err = r.repo.CreateTag(v.GitVersion(), commit, &git.CreateTagOptions{
		Message: fmt.Sprintf("chore(release): %s", v.FormatString()),
	}); err != nil {
		return nil, fmt.Errorf("create tag error: %w", err)
	}

	obj, err := r.repo.CommitObject(commit)
	if err != nil {
		return nil, fmt.Errorf("get commit object error: %w", err)
	}

	cmt := newCommitFromGit(*obj)

	cmt.Version = v

	return &cmt, nil
}

// CommitsArgs is a Commits options.
type CommitsArgs struct {
	NextV    version.V
	LastOnly bool
}

// Commits returns commits.
// If nextV is set, then the tag with this version is not created yet and nextV - new created version.
// In this case will ber returned commits from last tag to HEAD and last commit will be with nextV.
// If nextV is not set, then will be returned all commits.
func (r Repository) Commits(opt ...func(options *CommitsArgs)) ([]Commit, error) {
	a := &CommitsArgs{
		NextV:    version.V(""),
		LastOnly: false,
	}

	for _, o := range opt {
		o(a)
	}

	tags, err := r.tags()
	if err != nil {
		return nil, err
	}

	commits, err := r.repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	defer commits.Close()

	var cs []Commit

	if !a.NextV.Empty() {
		lastCommit := Commit{
			Hash:    plumbing.ZeroHash.String(),
			Message: fmt.Sprintf("chore(release): %s", a.NextV.GitVersion()),
			Version: a.NextV,
			Date:    time.Now(),
		}

		cs = make([]Commit, 0, 1)
		cs = append(cs, lastCommit)
	}

	for {
		c, err := commits.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		cmt := newCommitFromGit(*c)

		setTagToCommit(&cmt, tags)

		if a.LastOnly && cmt.IsTag() {
			break
		}

		cs = append(cs, cmt)
	}

	return cs, nil
}

// setTagToCommit sets a tag to commit.
func setTagToCommit(c *Commit, tags []tag) {
	for i := range tags {
		if c.Hash == tags[i].t.Target.String() {
			tagVer := tags[i].Version()

			if !tagVer.Invalid() {
				c.Version = tagVer
			}

			break
		}
	}
}

// tagExists returns true if the tag exists.
func (r Repository) tagExists(tag string) (bool, error) {
	tags, err := r.repo.TagObjects()
	if err != nil {
		return false, err
	}

	defer tags.Close()

	for {
		t, err := tags.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return false, err
		}

		if t.Name == tag {
			return true, nil
		}
	}

	return false, nil
}

// lastTag returns a last tag.
func (r Repository) lastTag() (*tag, error) {
	tags, err := r.tags()
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return nil, errTagsNotFound
	}

	return &tags[len(tags)-1], nil
}

// tags returns a list of tags.
func (r Repository) tags() ([]tag, error) {
	tagRefs, err := r.repo.Tags()
	if err != nil {
		return nil, err
	}

	defer tagRefs.Close()

	var tags []tag

	for {
		tagRef, err := tagRefs.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		if tagRef == nil {
			break
		}

		t, err := r.repo.TagObject(tagRef.Hash())
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag{t: *t})
	}

	return tags, nil
}

// Commit is a commit wrapper for return to external services.
type Commit struct {
	// Hash is a commit hash string.
	Hash string
	// Message is a commit message.
	Message string
	// Author is a commit author.
	Author string
	// Version is a commit version (for tag only).
	Version version.V
	// Date is a commit date.
	Date time.Time
	// Email is an user email.
	Email string
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

// tag is a tag wrapper for return to external services.
type tag struct {
	t object.Tag
}

// Version returns a tag version.
func (t tag) Version() version.V {
	v := version.V(t.t.Name)

	if v.Invalid() {
		return ""
	}

	return v
}
