package redir

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"
	"github.com/kshvmdn/redir/template"
)

// CreateOpts holds the set of options for a single create operation.
type CreateOpts struct {
	URL     string
	Name    string
	CNAME   string
	Private bool
	Verbose bool
}

var files = map[string]map[string]interface{}{
	"index.html": {"incl": true, "tmpl": template.Index},
	"CNAME":      {"incl": false, "tmpl": template.CNAME},
	"README.md":  {"incl": true, "tmpl": template.README},
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Create creates a new URL entry with it's associated repository.
func Create(opts *CreateOpts, client *github.Client, ctx context.Context) error {
	if opts.Name == "" {
		opts.Name = randomString(4)
	}

	if opts.CNAME != "" {
		files["CNAME"]["incl"] = true
	}

	repo, _, err := client.Repositories.Create(ctx, "", &github.Repository{
		AutoInit:     github.Bool(true),
		MasterBranch: github.String("refs/heads/gh-pages"),
		Name:         github.String(opts.Name),
		Private:      github.Bool(opts.Private),
	})
	if err != nil {
		return err
	}

	dir, err := createFiles(map[string]string{
		"url":   opts.URL,
		"title": opts.Name,
		"cname": opts.CNAME,
	})
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if err = runGitCmds(dir, repo.GetCloneURL(), opts.Verbose); err != nil {
		return err
	}

	// Use the API to delete the master branch, `git push origin --delete`
	// doesn't work for some reason :(.
	if _, err = client.Git.DeleteRef(
		ctx,
		repo.Owner.GetLogin(),
		repo.GetName(),
		"refs/heads/master",
	); err != nil {
		return err
	}

	// Append the new entry to the log and save the file.
	state.Log.Entries = append(state.Log.Entries, entry{
		Owner:       repo.Owner.GetLogin(),
		Repo:        repo.GetName(),
		RepoURL:     repo.GetHTMLURL(),
		RedirectURL: opts.URL,
		Private:     opts.Private,
	})
	return saveLog()
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

// createFiles creates all necessary files. These will always include an
// index.html, and a README.md. A CNAME file is created iff the --cname flag
// is provided.
func createFiles(params map[string]string) (string, error) {
	dir, err := ioutil.TempDir("", "redir-")
	if err != nil {
		return "", err
	}

	for file, opts := range files {
		if !opts["incl"].(bool) {
			continue
		}

		fi, err := os.Create(filepath.Clean(fmt.Sprintf("%s/%s", dir, file)))
		if err != nil {
			os.RemoveAll(dir)
			return "", err
		}

		tmpl := opts["tmpl"].(func(string, string, string) string)(
			params["url"],
			params["title"],
			params["cname"],
		)
		if _, err = fi.Write([]byte(tmpl)); err != nil {
			os.RemoveAll(dir)
			return "", err
		}
	}

	return dir, nil
}

// runGitCmds executes a series of git commands on the host machine
// to add, commit, and push the associated files to the gh-pages branch.
func runGitCmds(dir, remoteURL string, showOutput bool) error {
	for _, cmd := range [][]string{
		{"-C", dir, "init"},
		{"-C", dir, "add", "."},
		{"-C", dir, "commit", "-m", "Initial commit"},
		{"-C", dir, "branch", "-m", "gh-pages"},
		{"-C", dir, "remote", "add", "origin", remoteURL},
		{"-C", dir, "push", "-u", "origin", "gh-pages"},
	} {
		gitCmd := exec.Command("git", cmd...)
		if showOutput {
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
		}
		if err := gitCmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
