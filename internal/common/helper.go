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

func RunCommand(fullCmd string) ([]byte, error) {
	fmt.Println(fullCmd)

	cmd := exec.Command("sh", "-c", fullCmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("command failed:\n%s:\n%s%s", fullCmd, stderr.String(), stdout.String())
	}

	return stdout.Bytes(), nil
}
