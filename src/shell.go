package src

import (
	"bytes"
	"os/exec"
)

const (
	apiEndpoint = "localhost:9500"
)

func Exec_cmd(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = []string{"LD_LIBRARY_PATH=./", "DYLD_FALLBACK_LIBRARY_PATH=./"}
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	if errb.String() >= outb.String() {
		return errb.String(), nil
	} else {
		return outb.String(), nil
	}
}
