package gogitApiModels

type GoGitReference struct {
	Head    Heads
	Remotes Remotes
	Tags    Tags
	Type    ReferenceType
	Hash    string
}

type Heads struct {
	Name string
	Id   string
}

type Remotes struct {
	Name string
	Id   string
}

type Tags struct {
	Name string
	Id   string
}
type ReferenceType struct {
	Name string
	ID   int
}
