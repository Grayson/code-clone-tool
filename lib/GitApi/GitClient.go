package git

import (
	"log"

	"github.com/grayson/code-clone-tool/lib/shell"
)

type GitClient struct {
	log *log.Logger
}

func CreateGitClient(log *log.Logger) *GitClient {
	return &GitClient{
		log,
	}
}

func (gc *GitClient) Clone(gitUrl string, path string) (string, error) {
	gc.log.Println("Executing `git clone", gitUrl, path, "`")
	return shell.Do("git", "clone", gitUrl, path)
}

func (gc *GitClient) Pull(destinationDir string) (string, error) {
	gc.log.Println("Executing `git pull` in", destinationDir)
	return shell.In(destinationDir).Do("git", "pull")
}
