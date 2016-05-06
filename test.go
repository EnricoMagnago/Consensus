package main

import (
	"fmt"
	"strconv"
	"sync"
	"log"
)

type State int

const (
	ERROR State = -1 + iota // ERROR = -1, STOP = 0, RUNNING = 1
	STOP
	RUNNING
)
// interface for the function executed by threads.
type WorkerFunction func(*ProcessConfiguration, *AtomicBool, *int)  //function type, the last int is the retVal.

// parameters to the function, should not be modified after the start (concurrency).
type ProcessConfiguration struct {
	name string
}

//--------------ATOMIC BOOL-------------

// atomic access to a bool variable.
type AtomicBool struct {
	boolean bool
	mutex   sync.Mutex
}

func (abool *AtomicBool) set() {
	abool.mutex.Lock()
	abool.boolean = true
	abool.mutex.Unlock()
}

func (abool *AtomicBool) get() bool {
	var res bool = false
	abool.mutex.Lock()
	res = abool.boolean
	abool.mutex.Unlock()
	return res
}


// ----------------PROCESS--------------

type Process struct {
	configuration *ProcessConfiguration
	state         State
	terminator    AtomicBool
	function      WorkerFunction
	retVal        int
}

func newProcess(configuration *ProcessConfiguration, function WorkerFunction) Process {
	return Process{configuration, STOP, AtomicBool{false, sync.Mutex{}}, function, 0}
}

func (process *Process)start() bool {
	switch process.state {
	case ERROR:
		return false
	case RUNNING:
		return true
	case STOP:
		process.function(process.configuration, &process.terminator, &process.retVal) // start thread.
		process.state = RUNNING
		return true
	default:
		log.Fatal("ERROR: unknown state")
		return false
	}
}

func (process *Process)isAlive() bool {
	switch process.state {
	case ERROR:
		return false
	case STOP:
		return false
	case RUNNING:
		return true
	default:
		log.Fatal("ERROR: unknown state")
		return false
	}
}

func (process *Process)stop() bool {
	switch process.state {
	case ERROR:
		return false
	case STOP:
		return true
	case RUNNING:
		process.terminator.set()
		process.state = RUNNING
		return true
	default:
		log.Fatal("ERROR: unknown state")
		return false
	}
}

//---------------MANAGER--------------

type Manager struct {
	processes []*Process
}

func newManager() Manager {
	return Manager{make([]*Process, 0)}
}

func (manager *Manager) addProcess(conf *ProcessConfiguration, worker WorkerFunction) int {
	process := newProcess(conf, worker)
	manager.processes = append(manager.processes, &process)
	return len(manager.processes) - 1 // index in the slice.
}

func (manager *Manager) getProcessesNumber() int {
	return len(manager.processes)
}

//--------------FUNCTIONS IMPLEMENTING WORKER INTERFACE---------------

func BenOr(conf *ProcessConfiguration, terminator *AtomicBool, retVal *int) {
	for !*terminator.get() {
		*retVal = 1
	}
}

func main() {
	var processManager Manager = newManager()
	fmt.Printf("initSize: %d\n", processManager.getProcessesNumber())
	for i := 0; i < 10; i++ {
		var conf ProcessConfiguration = ProcessConfiguration{"pro" + strconv.Itoa(i)}
		var id int = processManager.addProcess(&conf, BenOr)
		fmt.Printf("%d)\n\tadded process, id: %d\n\tsize: %d\n", i, id, processManager.getProcessesNumber())
	}
}
