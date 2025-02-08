package main

func main() {
	monitor := monitor{}
	monitor.Start()
}

type monitor struct{}

func (m monitor) Start() {}
