package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the version number of gogit",
	Long:    `All software has versions. This is gogit's version.`,
	Example: "gogit version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Running GoGit Version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
