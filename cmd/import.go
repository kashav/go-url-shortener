package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a pre-existing config. file.",
	Long: `Replace the current configuration file. Can be used with the "export"
to restore a pre-existing configuration.

Example:
  $ shorten import shorten.toml`,
	Run: func(cmd *cobra.Command, args []string) {
		newConfFile := args[0]

		if _, err := os.Stat(newConfFile); os.IsNotExist(err) {
			log.Fatalf("File `%s` does not exist. Please provide a valid file.",
				newConfFile)
		}

		data, err := ioutil.ReadFile(newConfFile)
		checkError(err)

		err = ioutil.WriteFile(confFile, data, 0644)
		checkError(err)
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
}
