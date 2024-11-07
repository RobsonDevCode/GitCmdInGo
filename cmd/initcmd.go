package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/labstack/gommon/log"
)

var initCommand = &cobra.Command{
	Use:     "init",
	Short:   "The git init command is responsible for initializing a new Git repository in the specified directory. It does this by creating a .git folder that will contain the necessary subdirectories and files for managing the repository's history, objects, and configuration.",
	Example: "gogit init",
	Run: func(cmd *cobra.Command, args []string) {
		if err := InitGoGit(); err != nil {
			log.Fatalf("Error initializing: %d", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initCommand)
}

func InitGoGit() error {
	//get working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir = `C:\\Repos\\test\\dobetesting\\`

	initPath := filepath.Join(dir, ".gogit")

	//check if path has already been initialized
	if _, err := os.Stat(initPath); err == nil {
		log.Fatalf("Path:%s is already initialized or Error: %s", initPath, err)
		return err
	}

	if err := os.Mkdir(initPath, 0755); err != nil {
		log.Fatalf("failed when initializing path: %s, error: %s ", initPath, err)
		return err
	}

	//build files
	if err := CreateInitDirectories(initPath); err != nil {
		if err := os.Remove(".gogit"); err != nil {
			log.Fatalf("Error backtracking file .gogit: %s", err)
			return err
		}
		log.Fatalf("Error occured: %s. Removing file .gogit", err)
		return err
	}

	fmt.Printf("Initialized path: %s\n", initPath)

	//hide .git
	if err := SetHidden(&initPath); err != nil {
		log.Fatal(err)
	}
	return nil
}

// SetHidden set .gogit to hidden
func SetHidden(path *string) error {

	filenameW, err := syscall.UTF16PtrFromString(*path)
	if err != nil {
		log.Errorf("Error  UTF-16 encoding .gogit")
		return err
	}

	if err := syscall.SetFileAttributes(filenameW, syscall.FILE_ATTRIBUTE_HIDDEN); err != nil {
		log.Errorf("Error setting .gogit to hidden")
		return err
	}

	return nil
}

func CreateInitDirectories(initPath string) error {

	dirs := []string{
		"objects",
		"refs/heads",
		"refs/tags",
	}

	//create sub dir
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(initPath, dir), 0755); err != nil {
			return err
		}
	}

	//Create Head file
	if err := CreateHeadFile(initPath); err != nil {
		return err
	}

	//Create Config file
	if err := CreateConfigFile(initPath); err != nil {
		return err
	}

	//Create Description
	if err := CreateDescriptionFile(initPath); err != nil {
		return err
	}

	fmt.Println("Initialized empty Git repository in", initPath)
	return nil
}

// CreateHeadFile generate head file for init command
func CreateHeadFile(goGitPath string) error {
	headPath := filepath.Join(goGitPath, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: ref/head/master\n"), 0644); err != nil {
		log.Errorf("failed when creating %s: %v", headPath, err)
		return err
	}

	return nil
}

// CreateConfigFile generate config file for init command
func CreateConfigFile(goGitPath string) error {
	configPath := filepath.Join(goGitPath, "config")
	fmt.Printf("Initialized path: %s\n", configPath)

	configContent := `[core]
	repositoryformatversion = 0
	filemode = false
	bare = false
	logallrefupdates = true
	symlinks = false
	ignorecase = true
	[lfs]
	repositoryformatversion = 0`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		log.Errorf("falied when creating config: %s", err)
		return err
	}

	return nil
}

// CreateDescriptionFile generate description file
func CreateDescriptionFile(goGitPath string) error {
	descPath := filepath.Join(goGitPath, "description")

	descContent := "Unnamed repository; edit this file 'description' to name the repository."

	if err := os.WriteFile(descPath, []byte(descContent), 0664); err != nil {
		log.Errorf("failed when creating discription file: %s", err)
		return err
	}

	return nil
}
