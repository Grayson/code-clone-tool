package git

type GitApi interface {
	Clone(gitUrl string, path string) (string, error)
	Pull(destinationDir string) (string, error)
}
