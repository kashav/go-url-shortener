package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a shortened URL and delete the associated repository.",
	Long: `Remove a shortened URL and delete the associated repository. Specify
the URL with it's associated repository name.

Example:
  $ shorten remove go
  $ shorten remove x29jzI8m`,
	Run: func(cmd *cobra.Command, args []string) {
		repoName := args[0]

		repo, err := findRepo(repoName)
		if err != nil {
			fmt.Printf("Repository `%s` not found.\n", repoName)
			return
		}

		_, err = client.Repositories.Delete(ctx, repo.Owner, repo.Repo)
		checkError(err)

		fmt.Printf("Successfully removed %s/%s.\n", repo.Owner, repo.Repo)
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)
}

func findRepo(repoName string) (entry, error) {
	for i, entry := range conf.Entries {
		if entry.Repo == repoName {
			conf.Entries = append(conf.Entries[:i], conf.Entries[i+1:]...)
			saveConfig()
			return entry, nil
		}
	}

	return entry{}, errors.New("Repository not found.")
}
