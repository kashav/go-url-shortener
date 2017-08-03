package redir

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func Remove(repos []string, client *github.Client, ctx context.Context) error {
	for _, repo := range repos {
		i, en := findRepo(repo)
		if i < 0 {
			return fmt.Errorf("repository %s not found", repo)
		}

		if _, err := client.Repositories.Delete(ctx, en.Owner, en.Repo); err != nil {
			return err
		}

		state.Log.Entries = append(state.Log.Entries[:i], state.Log.Entries[i+1:]...)
		saveLog()
	}
	return nil
}

func findRepo(repo string) (int, *entry) {
	for i, en := range state.Log.Entries {
		if en.Repo == repo {
			return i, &en
		}
	}
	return -1, nil
}
