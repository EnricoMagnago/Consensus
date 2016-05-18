package main

import (
	"consensus/processManager"
	"consensus/process"
)

func main() {
	var processNumber int = 3
	var delayMean int = 100
	var variance int = 10
	var manager processManager.Manager = processManager.NewManager(processNumber, delayMean, variance)
	var workers [processNumber]process.WorkerFunction
	// TODO set workers functions
	manager.AddProcesses(workers)
}
