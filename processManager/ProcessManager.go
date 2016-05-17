package processManager

import(
	"consensus/process"
)

//---------------MANAGER--------------

type Manager struct {
	processes []*process.Process // change data structure.
}

func NewManager() Manager {
	return Manager{make([]*process.Process, 0)}
}

func (manager *Manager) AddProcess(conf *process.ProcessConfiguration, worker process.WorkerFunction) int {
	process := process.NewProcess(conf, worker)
	manager.processes = append(manager.processes, &process)
	return len(manager.processes) - 1 // index in the slice.
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

func (manager *Manager) WaitProcessTermination(processId int) {
	manager.processes[processId].WaitTermination()
}

func (manager *Manager) WaitProcessesTermination() {
	for i := 0; i < len(manager.processes); i++ {
		manager.WaitProcessTermination(i)
	}
}

func (manager *Manager) GetProcessesNumber() int {
	return len(manager.processes)
}
