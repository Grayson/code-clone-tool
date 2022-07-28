package shell

import (
	"bytes"
	"os/exec"
)

func Do(first string, args ...string) (string, error) {
	cmd := exec.Command(first, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return stdout.String(), nil
}
