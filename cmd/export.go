package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Print the current config. file.",
	Long: `Prints the Shorten configuration TOML file to stdout. Can be used with
the "import" command to restore a pre-existing configuration.

Assumes you've initialized the program with a valid access token, if not,
run this command with no arguments.

Example:
  $ shorten export > shorten.toml`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile(confFile)
		checkError(err)
		fmt.Print(string(data))
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
