package main

import (
	"fmt"
	"forever-monitor/internal/common"
	"os"
	"strconv"
	"strings"
	"time"
)

// Note: when there is no target process running, program startup will panic
//
// TODO: need to figure out how to properly `fork` a new process
// instead of new goroutine, because we don't expect the target process to
// shutdown when the monitor program stops

// TODO haven't test on linux yet
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
	monitoredProcesses map[int]*parsedProcess
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
			// spew.Dump(m.monitoredProcesses)
		case process, ok := <-m.addProcess:
			if ok {
				fmt.Println("adding process", process.pid)
				m.monitoredProcesses[process.pid] = process
			}
		case pid, ok := <-m.removePid:
			if ok {
				delete(m.monitoredProcesses, pid)
				fmt.Println("removing process", pid)
			}
		}
	}
}

func (m *monitor) parseProcess() map[int]*parsedProcess {
	fullCmd := fmt.Sprintf("ps | grep -v grep | grep -v main | grep -v %s | grep %s", m.currentProg, m.targetProg)

	result, _, err := common.RunCommand(fullCmd)
	common.Boom(err, "failed to run command", fullCmd)

	fmt.Println(string(result))
	lines := strings.Split(string(result), "\n")

	processes := make(map[int]*parsedProcess)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		detail, err := parseProcess(line)
		common.Boom(err)

		processes[detail.pid] = detail
		// spew.Dump(detail)
	}

	return processes
}

// expected running: 1,2,3
// actual   running: 1,2
func (m *monitor) compareProcess(running map[int]*parsedProcess) {
	for pid, process := range m.monitoredProcesses {
		if _, ok := running[pid]; !ok {
			m.restart(pid, process.cmd)
		}
	}
}

func (m *monitor) restart(pid int, cmd string) {
	m.removePid <- pid

	// TODO ideally fork a new process instead of fire new goroutine
	go func() {
		_, pid, err := common.RunCommand(cmd)
		common.Boom(err, "failed to restart", cmd)

		process := &parsedProcess{
			pid: pid,
			cmd: cmd,
		}

		m.addProcess <- process
	}()
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
