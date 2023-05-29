package app

type TestGit struct {
	DidClone bool
	DidPull  bool
}

func (g *TestGit) Clone(gitUrl string, path string) (string, error) {
	g.DidClone = true
	return "stdout", nil
}

func (g *TestGit) Pull(destinationDir string) (string, error) {
	g.DidPull = true
	return "stdout", nil
}
