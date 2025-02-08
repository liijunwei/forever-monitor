package main

import (
	"fmt"
	"forever-monitor/internal/common"
	"os"
	"time"
)

func main() {
	common.Assert(len(os.Args) == 2, "missing program name")

	foreverProgram := foreverProgram{}
	foreverProgram.Start()
}

type foreverProgram struct{}

func (m foreverProgram) Start() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			work(os.Getenv("LOG_ENABLED") == "1")
		}
	}
}

func work(logEnabled bool) {
	if logEnabled {
		fmt.Println(time.Now().Format(time.RFC3339), "tick")
	}
}
