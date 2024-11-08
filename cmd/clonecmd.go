package cmd

import (
	"GoGitCLI.RobsonDevCode.Com/extentions"
	. "GoGitCLI.RobsonDevCode.Com/models/gogitApiModels"
	"GoGitCLI.RobsonDevCode.Com/services/goGitEndpoints"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

// GoGitApiConnection Dependency Injection
type GoGitApiConnection struct {
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

// CloneGoGit implements basic git clone functionality
func CloneGoGit(cmd *cobra.Command, args []string) {
	repoUrl, err := cmd.Flags().GetString("repoUrl")
	if err != nil {
		log.Fatalf("Error reading Url: %d", err)
		return
	}

	if repoUrl == "" {
		log.Fatal("Repository Url is required")
		return
	}

	isWorkingUrl := isURL(repoUrl)
	if !isWorkingUrl {
		log.Fatal("Invalid Url")
	}

	//get working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//create destination for repo
	dirErr := CreateDirectory(dir)
	if dirErr != nil {
		log.Fatal(dirErr)
	}

	//initialize .gogit
	initErr := InitGoGit()
	if initErr != nil {
		log.Fatal(initErr)
	}

	//Fetch refs from gogit api
	refs, err := goGitEndpoints.FetchRefs(repoUrl)
	if err != nil {
		log.Errorf("")
	}
	if refs == nil {

	}
}

func isURL(input string) bool {
	//check if input is a URL
	parsedUrl, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}
	return parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
}

func CreateDirectory(dir string) error {

	dirExtensions := extentions.NewOSDirectoryExtensions()
	//force OOP on go so we can use these extensions, open to feed back from someone smarter than me
	isExist, err := dirExtensions.CheckIfDirectoryExists(dir)

	if err != nil {
		log.Errorf("Error checking if dir exists: %s", err)
		return err
	}

	if isExist {
		log.Errorf("Directory %s alrady exists", dir)
	}

	//create directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Errorf("failed to create destination directory: %v", err)
		return err
	}

	return nil
}
