package extentions

type IDirectoryExtensions interface {
	CheckIfFolderExists(folderPath string) (bool, error)
}
