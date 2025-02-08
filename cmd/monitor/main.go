package main

import (
	"fmt"
	"forever-monitor/internal/common"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	common.Assert(len(os.Args) == 2, "missing program name")

	currentProg := os.Args[0]
	targetProg := os.Args[1]

	monitor := newMonitor(currentProg, targetProg)
	monitor.Start()
}

type monitor struct {
	currentProg        string
	targetProg         string
	monitoredProcesses []*parsedProcess
	addProcess         chan *parsedProcess
	removePid          chan int
}

func newMonitor(currentProg, targetProg string) *monitor {
	m := &monitor{
		currentProg: currentProg,
		targetProg:  targetProg,
		addProcess:  make(chan *parsedProcess, 1),
		removePid:   make(chan int, 1),
	}

	m.monitoredProcesses = m.parseProcess()

	return m
}

// parse
// storge
// compare
// restart if needed
func (m *monitor) Start() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			runningProcesses := m.parseProcess()
			m.compareProcess(runningProcesses)
		case process, ok := <-m.addProcess:
			if ok {
				m.monitoredProcesses = append(m.monitoredProcesses, process)
			}
		case pid, ok := <-m.removePid:
			if ok {
				// m.removeProcess(pid)
				fmt.Println("removing pid", pid)
			}
		}
	}
}

func (m *monitor) parseProcess() []*parsedProcess {
	fullCmd := fmt.Sprintf("ps | grep -v grep | grep -v main | grep -v %s | grep %s", m.currentProg, m.targetProg)

	result, _, err := common.RunCommand(fullCmd)
	common.Boom(err, "failed to run command", fullCmd)

	fmt.Println(string(result))
	lines := strings.Split(string(result), "\n")

	processes := make([]*parsedProcess, 0, 100) // guess a number

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		detail, err := parseProcess(line)
		common.Boom(err)

		processes = append(processes, detail)
		// spew.Dump(detail)
	}

	return processes
}

// expected running: 1,2,3
// actuall  running: 1,2
func (m *monitor) compareProcess(running []*parsedProcess) {
	runningPids := make(map[int]bool)
	for _, p := range running {
		runningPids[p.pid] = true
	}

	// restart if needed
	for _, p := range m.monitoredProcesses {
		if !runningPids[p.pid] {
			fmt.Println("process", p.pid, "is not running, restarting")
			go m.restart(p.pid, p.cmd)
		}
	}
}

func (m *monitor) restart(pid int, cmd string) {
	m.removePid <- pid

	_, pid, err := common.RunCommand(cmd)
	common.Boom(err, "failed to restart", cmd)

	process := &parsedProcess{
		pid: pid,
		cmd: cmd,
	}

	m.addProcess <- process
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
