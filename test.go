package main

import (
	"fmt"
	"strconv"
)

type Process struct {
	name string
}

type Manager struct {
	processes []Process
}

func (manager *Manager) addProcess(name string) int {
	process := Process{name}
	manager.processes = append(manager.processes, process)
	return len(manager.processes) - 1
}

func (manager *Manager) getProcessesNumber() int {
	return len(manager.processes)
}

func main() {
	var processManager Manager = Manager{make([]Process, 0)}
	fmt.Println("initSize: %d", processManager.getProcessesNumber())
	for i := 0; i < 10; i++ {
		var id int = processManager.addProcess("pro" + strconv.Itoa(i))
		fmt.Printf("%d)\n\tadded process, id: %d\n\tsize: %d\n", i, id, processManager.getProcessesNumber())
	}
}
