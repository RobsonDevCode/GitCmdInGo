package extentions

import "os"

type IDirectoryExtensions interface {
	CheckIfDirectoryExists(folderPath string) (bool, error)
}

type osDirectoryExtensions struct{}

func NewOSDirectoryExtensions() IDirectoryExtensions {
	return &osDirectoryExtensions{}
}

func (e *osDirectoryExtensions) CheckIfDirectoryExists(folderPath string) (bool, error) {
	if _, err := os.Stat(folderPath); err != nil {
		if os.IsExist(err) {
			return true, err
		} else {
			//directory doesn't exist but cannot be created for other reasons
			return false, err
		}
	}

	return false, nil
}

type ExtensionHandler struct {
	IDirectoryExtensions IDirectoryExtensions
}
