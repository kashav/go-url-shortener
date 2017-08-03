package redir

import "fmt"

func List() {
	for i, en := range state.Log.Entries {
		fmt.Printf("%d) %s -> %s\n", i+1, en.Repo, en.RedirectURL)
	}
}
