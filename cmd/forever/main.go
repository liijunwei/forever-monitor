package main

func main() {
	foreverProgram := foreverProgram{}
	foreverProgram.Start()
}

type foreverProgram struct{}

func (m foreverProgram) Start() {}
