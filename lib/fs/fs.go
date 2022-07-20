package fs

type Fs interface {
	ChangeWorkingDirectory(string) error
}
