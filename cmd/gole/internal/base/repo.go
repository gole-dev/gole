package base

import (
	"context"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Repo is git repository manager.
type Repo struct {
	url  string
	home string
}

// NewRepo new a repository manager.
func NewRepo(url string) *Repo {
	var start int
	start = strings.Index(url, "//")
	if start == -1 {
		start = strings.Index(url, ":") + 1
	} else {
		start += 2
	}
	end := strings.LastIndex(url, "/")
	return &Repo{
		url:  url,
		home: goleHomeWithDir("repo/" + url[start:end]),
	}
}

func (r *Repo) Path() string {
	start := strings.LastIndex(r.url, "/")
	end := strings.LastIndex(r.url, ".git")
	if end == -1 {
		end = len(r.url)
	}
	return path.Join(r.home, r.url[start+1:end])
}

// Pull fetch the repository from remote url.
func (r *Repo) Pull(ctx context.Context) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = r.Path()
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

// Clone clones the repository to cache path.
func (r *Repo) Clone(ctx context.Context) error {
	if _, err := os.Stat(r.Path()); !os.IsNotExist(err) {
		return r.Pull(ctx)
	}
	cmd := exec.Command("git", "clone", r.url, r.Path())
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (r *Repo) CopyTo(ctx context.Context, to string, modPath string, ignores []string) error {
	if err := r.Clone(ctx); err != nil {
		return err
	}
	mod, err := ModulePath(path.Join(r.Path(), "go.mod"))
	if err != nil {
		return err
	}
	return copyDir(r.Path(), to, []string{mod, modPath}, ignores)
}
