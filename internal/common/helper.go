package common

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Boom(err error, msg ...string) {
	if err != nil {
		fmt.Println(strings.Join(msg, " "))
		panic(err)
	}
}

func Assert(ok bool, msg ...string) {
	if !ok {
		panic("assertion failed: " + strings.Join(msg, " "))
	}
}

func RunCommand(fullCmd string) ([]byte, int, error) {
	cmd := exec.Command("sh", "-c", fullCmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, 0, fmt.Errorf("command failed:\n%s:\n%s%s", fullCmd, stderr.String(), stdout.String())
	}

	pid := cmd.Process.Pid

	return stdout.Bytes(), pid, nil
}
