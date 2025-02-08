package main

import (
	"fmt"
	"forever-monitor/internal/common"
	"os"
	"time"
)

func main() {
	common.Assert(len(os.Args) == 2, "missing program name")

	foreverProgram := foreverProgram{
		args: os.Args[1],
	}
	foreverProgram.Start()
}

type foreverProgram struct {
	args any
}

func (m foreverProgram) Start() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			work(m.args, os.Getenv("LOG_ENABLED") == "1")
		}
	}
}

func work(args any, logEnabled bool) {
	if logEnabled {
		fmt.Println(time.Now().Format(time.RFC3339), args, "tick")
	}
}
