package main

import (
	"fmt"
	"forever-monitor/internal/common"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	common.Assert(len(os.Args) == 2, "missing program name")

	currentProg := os.Args[0]
	targetProg := os.Args[1]

	monitor := newMonitor(currentProg, targetProg)
	monitor.Start()
}

type monitor struct {
	currentProg string
	targetProg  string
}

func newMonitor(currentProg, targetProg string) *monitor {
	return &monitor{
		currentProg: currentProg,
		targetProg:  targetProg,
	}
}

func (m *monitor) Start() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			m.checkProcess()
		}
	}
}

// parse
// storge
// compare
// restart if needed
func (m *monitor) checkProcess() *parsedProcess {
	fullCmd := fmt.Sprintf("ps | grep -v grep | grep -v main | grep -v %s | grep %s", m.currentProg, m.targetProg)

	result, err := common.RunCommand(fullCmd)
	common.Boom(err, "failed to run command", fullCmd)

	fmt.Println(string(result))
	lines := strings.Split(string(result), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		detail, err := parseProcess(line)
		common.Boom(err)
		spew.Dump(detail)
	}

	// panic("TODO")
	return nil
}

// line example
// 91013 ttys009    0:00.17 ./tmp/forever 100
func parseProcess(line string) (*parsedProcess, error) {
	fields := strings.Fields(line)
	common.Assert(len(fields) >= 4, "unexpected output, not enough fields")

	pid, err := strconv.Atoi(fields[0])
	common.Boom(err)

	cmd := strings.Join(fields[3:], " ")

	return &parsedProcess{
		pid: pid,
		cmd: cmd,
	}, nil
}

type parsedProcess struct {
	pid int
	cmd string
}
