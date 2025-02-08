package main

import (
	"bytes"
	"fmt"
	"forever-monitor/internal/common"
	"os"
	"os/exec"
)

func main() {
	common.Assert(len(os.Args) == 2, "missing program name")

	currentProg := os.Args[0]
	targetProg := os.Args[1]
	fullCmd := fmt.Sprintf("ps aux | grep -v grep | grep -v main | grep -v %s | grep %s", currentProg, targetProg)

	result, err := runCommand(fullCmd)
	if err != nil {
		fmt.Println("failed to run ps aux", err.Error())
		os.Exit(1)
	}

	common.Boom(err, "failed to run ps aux")

	fmt.Println(string(result))
	monitor := monitor{}
	monitor.Start()
}

type monitor struct{}

func (m monitor) Start() {}

func runCommand(fullCmd string) ([]byte, error) {
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
