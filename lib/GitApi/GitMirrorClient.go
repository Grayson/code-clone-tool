package git

type GitMirrorClient struct {
}

func (g *GitMirrorClient) Clone(gitUrl string, path string) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (g *GitMirrorClient) Pull(destinationDir string) (string, error) {
	panic("not implemented") // TODO: Implement
}
