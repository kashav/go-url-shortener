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
the URL with it's associated repository name (may contains more than 1 repo separated by space as parameters).

Example:
  $ shorten remove go
  $ shorten remove x29jzI8m
  $ shorten remove a243234a jasdH234 524aAdsd`,
	Run: func(cmd *cobra.Command, args []string) {

		// do input check before running command
		for len(args) == 0 {
			var repoName string
			fmt.Print("Please input the name(s) of the repos you'd like to remove:  : ")
			fmt.Scanln(&repoName)

			if repoName != "" {
				args = append(args, repoName)
			}
		}

		// then do the repo removal
		for _, arg := range args {
			err := removeRepo(arg)
			checkError(err)
		}

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

	return entry{}, errors.New("repository not found")
}

func removeRepo(repoName string) error {
	repo, err := findRepo(repoName)
	if err != nil {
		return fmt.Errorf("repository '%s' not found", repoName)
	}

	_, err = client.Repositories.Delete(ctx, repo.Owner, repo.Repo)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully removed %s/%s.\n", repo.Owner, repo.Repo)

	return nil
}
