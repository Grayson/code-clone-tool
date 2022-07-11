package lib

import (
	"bytes"
	"log"
	"os/exec"
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
	cmd := exec.Command("git", "clone", gitUrl, path)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return stdout.String(), nil
}

func (gc *GitClient) Pull(destinationDir string) (string, error) {
	gc.log.Println("Executing `git pull` in", destinationDir)
	cmd := exec.Command("git", "pull")
	cmd.Dir = destinationDir
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return stdout.String(), nil
}
