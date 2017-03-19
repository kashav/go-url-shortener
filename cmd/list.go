package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print a list of currently active URLs.",
	Long: `Print a list of currently active URLs.

Example:
  $ shorten list`,
	Run: func(cmd *cobra.Command, args []string) {
		for i, entry := range conf.Entries {
			fmt.Printf("%d) %s -> %s\n", i+1, entry.Repo, entry.RedirectUrl)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
