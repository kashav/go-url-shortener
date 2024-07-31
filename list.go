package shortener

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// Lister implements Runner for the list operation.
type Lister struct{}

func (l *Lister) run(ctx context.Context, client *github.Client) (err error) {
	for i, en := range state.Log.Entries {
		fmt.Printf("%d) %s -> %s\n", i+1, en.Name, en.RedirectURL)
	}
	return nil
}
