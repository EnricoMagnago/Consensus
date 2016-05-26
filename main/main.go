package main

import (
	"consensus/processManager"
	"consensus/process"
)

func main() {
	var processNumber int = 3
	var delayMean int = 2000
	var variance int = 1000
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
