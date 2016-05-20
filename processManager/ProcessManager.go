package processManager

import (
	"consensus/process"
	"consensus/channel"
	"consensus/util"
)

//---------------MANAGER--------------

type Manager struct {
	processes     []*process.Process // change data structure.
	channel       *channel.Channel
	processNumber int
}

func NewManager(processNumber int, mean int, variance int) Manager {
	return Manager{make([]*process.Process, 0, processNumber), channel.NewChannel(processNumber, mean, variance), processNumber}
}

func (manager *Manager) addProcess(worker process.WorkerFunction, conf *process.ProcessConfiguration) int {
	process := process.NewProcess(conf, worker)
	manager.processes = append(manager.processes, &process)
	return len(manager.processes) - 1 // index in the slice.
}

func (manager *Manager) AddProcesses(workers []process.WorkerFunction) {
	for i := 0; i < manager.processNumber; i++ {
		var conf *process.ProcessConfiguration = process.NewProcessConfiguration(manager.channel, i, manager.processNumber)
		manager.addProcess(workers[i], conf)
	}
}

func (manager *Manager) StartProcess(processId int) bool {
	return manager.processes[processId].Start()
}

func (manager *Manager) StartProcesses() bool {
	for i := 0; i < len(manager.processes); i++ {
		if (!manager.StartProcess(i)) {
			return false
		}
	}
	return true
}

func (manager *Manager) StopProcess(processId int) bool {
	return manager.processes[processId].Stop()
}

func (manager *Manager) StopProcesses() bool {
	for i := 0; i < len(manager.processes); i++ {
		if (!manager.StopProcess(i)) {
			return false
		}
	}
	return true
}

func (manager *Manager) WaitProcessesTermination() {
	for i := 0; i < len(manager.processes); i++ {
		manager.GetRetval(i)
	}
}

func (manager *Manager) GetRetval(processId int) *util.RetVal {
	return manager.processes[processId].WaitTermination()
}

func (manager *Manager) GetProcessesNumber() int {
	return len(manager.processes)
}
