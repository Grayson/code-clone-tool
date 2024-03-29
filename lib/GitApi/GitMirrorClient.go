package git

import (
	"log"

	"github.com/grayson/code-clone-tool/lib/shell"
)

type GitMirrorClient struct {
	log *log.Logger
}

func CreateMirrorClient(log *log.Logger) GitApi {
	return &GitMirrorClient{
		log,
	}
}

func (g *GitMirrorClient) Clone(gitUrl string, path string) (string, error) {
	g.log.Printf("Executing `git clone --mirror %v %v`", gitUrl, path)
	g.log.Println()
	return shell.Do("git", "clone", "--mirror", gitUrl, path)
}

func (g *GitMirrorClient) Pull(destinationDir string) (string, error) {
	g.log.Printf("Executing `git remote update` in %v", destinationDir)
	g.log.Println()
	return shell.In(destinationDir).Do("git", "remote", "update")
}
