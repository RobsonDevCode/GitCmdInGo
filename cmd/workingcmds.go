package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	BlobType = "blob"
)

// Entry represents a file entry in the index
type Entry struct {
	Path     string
	Hash     string
	Mode     int
	Size     int64
	Modified int64
}

// Index represents the git staging area
type Index struct {
	Entries map[string]Entry
}

// IndexHeader represents the binary format of Git's index file header
type IndexHeader struct {
	Signature  [4]byte
	Version    uint32
	EntryCount uint32
}

const (
	IndexSignature = "DIRC"
	IndexVersion   = 2
)

var add = &cobra.Command{
	Use:   "add",
	Short: "Add file contents to index",
	Args:  cobra.ExactArgs(10),
	Run: func(cmd *cobra.Command, args []string) {

		//get working dir
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		if err := Add(dir, args); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(add)
}

func Add(repoPath string, paths []string) error {
	if repoPath == "" {
		return errors.New("file is required when using cmd add")
	}

	if repoPath == "." {
		err := addAllFiles()
		if err != nil {
			log.Fatalf("Error adding all files: %s", err)
			return err
		}
		return nil
	}

	index := &Index{
		Entries: make(map[string]Entry),
	}

	if err := loadIndex(repoPath, index); err != nil {
		log.Fatalf("Error loading index: %s", err)
		return err
	}

	for _, path := range paths {
		err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
			if err != nil {
				log.Error(err)
			}

			//skip .git file
			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			//skipping dir
			if info.IsDir() {
				return nil
			}

			if _, exists := index.Entries[currentPath]; !exists {
				//check if the file changed
				changed, fileErr := hasFileChanged(repoPath, currentPath, index.Entries[currentPath])
				if fileErr != nil {
					log.Errorf("error checking if change in file %s: %s", currentPath, fileErr)
					return fileErr
				}

				// no need to add as the file is unchanged
				if !changed {
					log.Infof(" %s ", currentPath)
					return nil
				}

				//create blob
				hash, err := createBlob(repoPath, currentPath)
				if err != nil {
					log.Errorf("error trying to create blob for file,%s error: %s", err)
				}

				index.Entries[currentPath] = Entry{
					Path:     currentPath,
					Hash:     hash,
					Mode:     int(info.Mode()),
					Size:     info.Size(),
					Modified: info.ModTime().Unix(),
				}

				return nil

			}

			return nil
		})

		if err != nil {
			log.Fatal(err)
			return err
		}

		return saveIndex(repoPath, index)
	}

	return nil
}

func addAllFiles() error {
	return nil
}

func loadIndex(repoPath string, index *Index) error {
	indexPath := filepath.Join(repoPath, ".git", "index")

	//check if index already exists
	if _, err := os.Stat(indexPath); err != nil {
		if os.IsNotExist(err) {
			return createNewIndex(indexPath, index)
		}
		log.Errorf("Error Creating file index: %s", err)
		return err
	}

	file, err := os.Open(indexPath)
	if err != nil {
		log.Errorf("Error when opening index file: %s", err)
		return err
	}
	defer file.Close()

	//check if file being added has changed

	return nil
}

// hasFileChanged checks if the current file has been updated
func hasFileChanged(repoPath string, currentPath string, entry Entry) (bool, error) {
	info, err := os.Stat(currentPath)
	if err != nil {
		return false, err
	}

	//check if file size or the mod time on the file has changed
	if info.Size() != entry.Size || info.ModTime().Unix() != entry.Modified {
		return true, nil
	}

	return false, nil
}

func createNewIndex(indexPath string, index *Index) error {

	file, err := os.Create(indexPath)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := file.Close(); err != nil {
		log.Errorf("Error closing file at path, %s Error: %s ", indexPath, err)
		return err
	}

	index.Entries = make(map[string]Entry)

	return nil
}

func createBlob(repoPath, currentPath string) (string, error) {

	file, err := os.Open(currentPath)
	if err != nil {
		log.Errorf("Error opening file at path, %s Error: %s ", currentPath, err)
		return "", err
	}
	defer file.Close()

	hashing := sha256.New()

	//Add all file to contents
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Error reading file at path, %s Error: %s ", currentPath, err)
		return "", err
	}

	//Write headers
	if _, fmtErr := fmt.Fprintf(hashing, "%s %d\x00", BlobType, len(content)); fmtErr != nil {
		log.Errorf("Error hashing file at path, %s Error: %s ", currentPath, fmtErr)
		return "", fmtErr
	}

	hashing.Write(content)
	hash := fmt.Sprintf("%x", hashing.Sum(nil))

	//create obj directory if it doesnt exist
	objDir := filepath.Join(repoPath, ".git", "objects", hash[:2])
	if writeErr := os.WriteFile(objDir, content, 0644); writeErr != nil {
		log.Errorf("Error writing file at path, %s Error: %s ", objDir, writeErr)
		return "", writeErr
	}

	return hash, nil
}

func saveIndex(repoPath string, index *Index) error {
	var buf bytes.Buffer

	_, err := buf.WriteString(IndexSignature)
	if err != nil {
		log.Errorf("Error writing index signature: %s", err)
		return err
	}

	if writeErr := binary.Write(&buf, binary.LittleEndian, uint32(IndexVersion)); writeErr != nil {
		log.Errorf("Error writing index version: %s", err)
		return writeErr
	}

	if writeErr := binary.Write(&buf, binary.LittleEndian, uint32(len(index.Entries))); writeErr != nil {
		log.Errorf("Error writing index entries: %s", writeErr)
		return writeErr
	}

	return nil

}
