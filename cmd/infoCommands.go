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

var helpCmd = &cobra.Command{
	Use:     "help",
	Short:   "Print the help list for all gogit Commands",
	Example: "gogit help",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`usage: git [--version] [--help] [-C <path>] [-c <name>=<value>]
           <command> [<args>]

These are common GoGit commands used in various situations:

start a working area (see also: git help tutorial)
   init        Create an empty Git repository or reinitialize an existing one
   clone       Clone a repository into a new directory

work on the current change (see also: git help everyday)
   add         Add file contents to the index
   mv          Move or rename a file, a directory, or a symlink
   restore     Restore working tree files
   rm          Remove files from the working tree and from the index

examine the history and state (see also: git help revisions)
   log         Show commit logs
   status      Show the working tree status
   diff        Show changes between commits, commit and working tree, etc
   show        Show various types of objects

grow, mark and tweak your common history
   branch      List, create, or delete branches
   commit      Record changes to the repository
   merge       Join two or more development histories together
   rebase      Reapply commits on top of another base tip
   tag         Create, list, delete or verify a tag object signed with GPG

collaborate (see also: git help workflows)
   fetch       Download objects and refs from another repository
   pull        Fetch from and integrate with another repository or a local branch
   push        Update remote refs along with associated objects`)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(helpCmd)
}
