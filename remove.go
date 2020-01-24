package point

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
)

// Remover implements Runner for the remove operation.
type Remover struct {
	Repos   []string
	Verbose bool
}

func (r *Remover) run(ctx context.Context, client *github.Client) error {
	for _, repo := range r.Repos {
		i, ent := r.findRepo(repo)
		if i < 0 {
			return fmt.Errorf("could not find entry %s", repo)
		}

		var err error
		if ent.IsSubdir {
			err = r.removeSubdir(ctx, client, ent)
		} else {
			_, err = client.Repositories.Delete(ctx, ent.Owner, ent.Repo)
		}
		if err != nil {
			return err
		}

		state.Log.Entries = append(state.Log.Entries[:i], state.Log.Entries[i+1:]...)
		if err := saveLog(); err != nil {
			return err
		}
		fmt.Printf("Removed entry %s.\n", ent.Name)
	}
	return nil
}

func (r *Remover) removeSubdir(ctx context.Context, client *github.Client, ent *entry) error {
	repo, _, err := client.Repositories.Get(ctx, ent.Owner, ent.Repo)
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "point-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	for _, cmd := range [][]string{
		{"clone", repo.GetCloneURL(), dir},
		{"-C", dir, "rm", "-r", ent.Name[strings.LastIndex(ent.Name, "/")+1:]},
		{"-C", dir, "commit", "-m", fmt.Sprintf("point: remove entry (%s)", ent.Name)},
		{"-C", dir, "push"},
	} {
		gitCmd := exec.Command("git", cmd...)
		if r.Verbose {
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
		}
		if err := gitCmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Remover) findRepo(repo string) (int, *entry) {
	// try to parse as int
	entryNum, err := strconv.Atoi(repo)
	if err == nil {
		// interpert input as entry number
		entryNum-- // `point list` index starts from 1
		if 0 <= entryNum && entryNum < len(state.Log.Entries) {
			return entryNum, &state.Log.Entries[entryNum]
		}
	}
	for i, ent := range state.Log.Entries {
		if ent.Name == repo {
			return i, &ent
		}
	}
	return -1, nil
}
