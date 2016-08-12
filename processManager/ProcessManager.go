package processManager

import (
	"consensus/process"
	"consensus/channel"
	"consensus/util"
	"fmt"
)

//---------------MANAGER--------------

type Manager struct {
	processes     []*process.Process // change data structure.
	channel       *channel.Channel
	processNumber int
	f             int
	maxVal        int
}

func NewManager(processNumber int, mean int, variance int, f int, maxVal int) Manager {
	return Manager{make([]*process.Process, 0, processNumber), channel.NewChannel(processNumber, mean, variance), processNumber, f, maxVal}
}

func (manager *Manager) addProcess(worker process.WorkerFunction, conf *process.ProcessConfiguration) int {
	process := process.NewProcess(conf, worker)
	manager.processes = append(manager.processes, &process)
	return len(manager.processes) - 1 // index in the slice.
}

func (manager *Manager) AddProcesses(workers []process.WorkerFunction) {
	for i := 0; i < manager.processNumber; i++ {
		var conf *process.ProcessConfiguration = process.NewProcessConfiguration(manager.channel, i, manager.processNumber, manager.f, manager.maxVal)
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

func (manager *Manager) WaitProcessesTermination() []*util.RetVal {
	var res []*util.RetVal = make([]*util.RetVal, 0)
	for i := 0; i < len(manager.processes); i++ {
		//fmt.Printf("Waiting for %d...\n", i);
		res = append(res, manager.GetRetval(i))
	}
	return res
}

func (manager *Manager) GetRetval(processId int) *util.RetVal {
	return manager.processes[processId].WaitTermination()
}

func (manager *Manager) GetProcessesNumber() int {
	return len(manager.processes)
}

func (manager *Manager) GetSendedMessageNumber() int {
	return manager.channel.GetSendedMessagesNumber()
}

func (manager *Manager) GetDeliveredMessageNumber() int {
	return manager.channel.GetDeliveredMessagesNumber()
}

func (manager *Manager) LogState() {
	fmt.Print("Channel State:\n")
	manager.channel.PrintState()
}
