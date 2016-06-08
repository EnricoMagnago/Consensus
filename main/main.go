package main

import (
	"consensus/processManager"
	"consensus/process"
)

func main() {
	var processNumber int = 5
	var delayMean int = 000
	var variance int = 000
	var manager processManager.Manager = processManager.NewManager(processNumber, delayMean, variance)
	var workers []process.WorkerFunction = make([]process.WorkerFunction, 0, processNumber)
	for i := 0; i < processNumber; i++ {
		workers = append(workers, BenOr)
	}
	manager.AddProcesses(workers)
	manager.StartProcesses()
	//time.Sleep(time.Duration(5) * time.Second)
	manager.WaitProcessesTermination()

}
