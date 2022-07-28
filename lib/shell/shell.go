package shell

import (
	"bytes"
	"os/exec"
)

func Do(first string, args ...string) (string, error) {
	return In("").Do(first, args...)
}

type impl struct {
	destination string
}

func In(destination string) *impl {
	return &impl{
		destination: destination,
	}
}

func (i *impl) Do(first string, args ...string) (string, error) {
	cmd := exec.Command(first, args...)
	if i.destination != "" {
		cmd.Dir = i.destination
	}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return stdout.String(), nil
}
