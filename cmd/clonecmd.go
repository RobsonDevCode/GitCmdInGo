package cmd

import (
	"GoGitCLI.RobsonDevCode.Com/extentions"
	"GoGitCLI.RobsonDevCode.Com/services/goGitEndpoints"
	"errors"
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
	Args:    cobra.ExactArgs(1),
	Example: "gogit clone [repository Url]",
	Run: func(cmd *cobra.Command, args []string) {
		CloneGoGit(args[0])
	},
}

var (
	repoUrl string
)

func init() {
	rootCmd.AddCommand(cloneCommand)
}

// CloneGoGit implements basic git clone functionality
func CloneGoGit(url string) error {
	isWorkingUrl := isURL(url)
	if !isWorkingUrl {
		return errors.New("Invalid Url")
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
		return initErr
	}

	//Fetch refs from gogit api
	refs, fetchErr := goGitEndpoints.FetchRefs(url)

	if fetchErr != nil {
		log.Errorf("Error Fetching Refs: %s", fetchErr)
		return fetchErr
	}
	if refs == nil {
		log.Error("No Refs Found")
	}

	log.Infof("%s was clone into %s ", url, dir)

	return nil
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
