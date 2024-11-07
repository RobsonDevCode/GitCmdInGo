package cmd

import (
	. "GoGitCLI.RobsonDevCode.Com/extentions"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"net/url"
)

type DirectoryExtension struct {
	_dirExtentions IDirectoryExtensions
}

var cloneCommand = &cobra.Command{
	Use:     "clone",
	Short:   "Clone a repository into a new directory",
	Example: "gogit clone <repository Url>",
	Run:     CloneGoGit,
}

func init() {
	rootCmd.AddCommand(cloneCommand)
}

func CloneGoGit(cmd *cobra.Command, args []string) {
	input, err := cmd.Flags().GetString("repoUrl")
	if err != nil {
		log.Fatalf("Error reading Url: %d", err)
		return
	}

	if input == "" {
		log.Fatal("Repository Url is required")
		return
	}

	isWorkingUrl := isURL(input)
	if !isWorkingUrl {
		log.Fatal("Invalid Url")
	}

	//create destination for repo
}

func isURL(input string) bool {
	//check if input is a URL
	parsedUrl, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}
	return parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
}

func (c *DirectoryExtension) CreateDirectory(dir string) error {
	isExist, err := c._dirExtentions.CheckIfFolderExists(dir)
	if err != nil {
		log.Errorf("Error checking if dir exists: %s", err)
		return err
	}

	if isExist {

	}
	return nil
}
