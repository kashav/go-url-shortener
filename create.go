package redir

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/kshvmdn/redir/template"
)

// Creator implements Runner for the create operation.
type Creator struct {
	URL     string
	Name    string
	CNAME   string
	Private bool
	Repo    string
	Subdir  bool
	Verbose bool
}

var files = map[string]map[string]interface{}{
	"index.html": {"incl": true, "tmpl": template.Index},
	"CNAME":      {"incl": false, "tmpl": template.CNAME},
	"README.md":  {"incl": true, "tmpl": template.README},
}

func (c *Creator) run(ctx context.Context, client *github.Client) error {
	if c.Name == "" {
		c.Name = randomString(4)
	}

	dir, err := ioutil.TempDir("", "redir-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	var createFunc = c.createRepo
	if c.Subdir {
		createFunc = c.createSubdir
	}
	ent, err := createFunc(ctx, client, dir)
	if err != nil {
		return err
	}
	if ent == nil {
		return errors.New("failed to create entry")
	}
	state.Log.Entries = append(state.Log.Entries, *ent)
	return saveLog()
}

func (c *Creator) createSubdir(ctx context.Context, client *github.Client, dir string) (*entry, error) {
	if c.Repo == "" {
		return nil, errors.New("expected repo with --subdir flag (use --repo)")
	}

	split := strings.Split(c.Repo, "/")
	owner, repoName := split[0], split[1]

	repo, _, err := client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, fmt.Errorf("repository %s doesn't exist", c.Repo)
	}

	if err := c.runGitCmds([]string{"clone", repo.GetCloneURL(), dir}); err != nil {
		return nil, err
	}

	subdir := filepath.Clean(fmt.Sprintf("%s/%s", dir, c.Name))
	if err := os.MkdirAll(subdir, os.ModePerm); err != nil {
		return nil, err
	}
	if err := c.createFiles(subdir); err != nil {
		return nil, err
	}

	cmds := [][]string{
		{"-C", dir, "add", "."},
		{"-C", dir, "commit", "-m", "redir: new entry", "-m", fmt.Sprintf("%s -> %s", c.Name, c.URL)},
		{"-C", dir, "push"},
	}
	if err := c.runGitCmds(cmds...); err != nil {
		return nil, err
	}

	return &entry{
		Name:        fmt.Sprintf("%s/%s", repo.GetFullName(), c.Name),
		Owner:       repo.Owner.GetLogin(),
		Repo:        repo.GetName(),
		RepoURL:     repo.GetHTMLURL(),
		RedirectURL: c.URL,
		IsSubdir:    true,
		IsPrivate:   c.Private,
	}, nil
}

func (c *Creator) createRepo(ctx context.Context, client *github.Client, dir string) (*entry, error) {
	if c.CNAME != "" {
		files["CNAME"]["incl"] = true
	}

	repo, _, err := client.Repositories.Create(ctx, "", &github.Repository{
		AutoInit:     github.Bool(true),
		MasterBranch: github.String("refs/heads/gh-pages"),
		Name:         github.String(c.Name),
		Private:      github.Bool(c.Private),
	})
	if err != nil {
		return nil, err
	}

	if err := c.createFiles(dir); err != nil {
		return nil, err
	}

	cmds := [][]string{
		{"-C", dir, "init"},
		{"-C", dir, "add", "."},
		{"-C", dir, "commit", "-m", "Initial commit"},
		{"-C", dir, "branch", "-m", "gh-pages"},
		{"-C", dir, "remote", "add", "origin", repo.GetCloneURL()},
		{"-C", dir, "push", "-u", "origin", "gh-pages"},
	}
	if err := c.runGitCmds(cmds...); err != nil {
		return nil, err
	}

	// Use the API to delete the master branch, `git push origin --delete`
	// doesn't work for some reason :(.
	if _, err = client.Git.DeleteRef(
		ctx,
		repo.Owner.GetLogin(),
		repo.GetName(),
		"refs/heads/master",
	); err != nil {
		return nil, err
	}

	return &entry{
		Name:        repo.GetFullName(),
		Owner:       repo.Owner.GetLogin(),
		Repo:        repo.GetName(),
		RepoURL:     repo.GetHTMLURL(),
		RedirectURL: c.URL,
		IsSubdir:    false,
		IsPrivate:   c.Private,
	}, nil
}

// createFiles creates all necessary files. These will always include an
// index.html, and a README.md. A CNAME file is created iff the --cname flag
// is provided.
func (c *Creator) createFiles(dir string) error {
	for file, opts := range files {
		if !opts["incl"].(bool) {
			continue
		}

		fi, err := os.Create(filepath.Clean(fmt.Sprintf("%s/%s", dir, file)))
		if err != nil {
			return err
		}

		tmpl := opts["tmpl"].(func(string, string, string) string)(
			c.URL,
			c.Name,
			c.CNAME,
		)
		if _, err = fi.Write([]byte(tmpl)); err != nil {
			return err
		}
	}

	return nil
}

// runGitCmds executes a series of git commands on the host machine
// to add, commit, and push the associated files to the gh-pages branch.
func (c *Creator) runGitCmds(cmds ...[]string) error {
	for _, cmd := range cmds {
		gitCmd := exec.Command("git", cmd...)
		if c.Verbose {
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
		}
		if err := gitCmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// randomString returns a random hash of n characters.
func randomString(n int) string {
	var chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
