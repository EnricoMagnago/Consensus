package main

import (
	"fmt"
	"strconv"
	"sync"
)

type State int

const (
	ERROR State = -1 + iota // ERROR = -1, STOP = 0, RUNNING = 1
	STOP
	RUNNING
)
// interface for the function executed by threads, must periodically check if AtomicBool is setted, if it's the case the thread must terminate.
type WorkerFunction func(*ProcessConfiguration, *AtomicBool, *RetVal)  //function type, the last int is the retVal.

// parameters to the function, should not be modified after the start (concurrency).
type ProcessConfiguration struct {
	name         string
	channel      chan string
	processesIds []int
}

func newProcessConfiguration(name string, channel chan string, processesIds []int) *ProcessConfiguration {
	return &ProcessConfiguration{name, channel, processesIds}
}

//--------------ATOMIC STATE-------------

// atomic access to a bool variable.
type AtomicState struct {
	value State
	mutex sync.RWMutex // read-write mutex.
}

func newAtomicState() *AtomicState {
	return &AtomicState{STOP, sync.RWMutex{}}
}

func (astate *AtomicState) set(new State) bool {
	var res = true
	astate.mutex.Lock()
	switch astate.value{
	case ERROR: res = false
	case STOP: astate.value = new
	case RUNNING: astate.value = new
	default:
		fmt.Errorf("ERROR: AtomicState.set: Unknow state")
	}
	astate.mutex.Unlock()
	return res
}

func (astate *AtomicState) get() State {
	var res State
	astate.mutex.RLock()
	res = astate.value
	astate.mutex.RUnlock()
	return res
}


//--------------ATOMIC BOOL-------------

// atomic access to a bool variable.
type AtomicBool struct {
	boolean bool
	mutex   sync.RWMutex // read-write mutex.
}

func newAtomicBool() *AtomicBool {
	return &AtomicBool{false, sync.RWMutex{}}
}

func (abool *AtomicBool) set() {
	abool.mutex.Lock()
	abool.boolean = true
	abool.mutex.Unlock()
}

func (abool *AtomicBool) get() bool {
	var res bool = false
	abool.mutex.RLock()
	res = abool.boolean
	abool.mutex.RUnlock()
	return res
}

// --------------RETVAL---------------

type RetVal struct {
	integer int
	mutex   sync.Mutex
}

func newRetVal() *RetVal {
	return &RetVal{0, sync.Mutex{}}
}

func (retVal *RetVal)lock() {
	retVal.mutex.Lock()
}

func (retVal *RetVal)unlock() {
	retVal.mutex.Unlock()
}

func (retVal *RetVal)set(newValue int) {
	retVal.integer = newValue
}


// ----------------PROCESS--------------

type Process struct {
	configuration *ProcessConfiguration
	state         *AtomicState
	terminator    *AtomicBool
	function      WorkerFunction
	retVal        *RetVal
	endChannel    chan bool
}

func newProcess(configuration *ProcessConfiguration, function WorkerFunction) Process {
	return Process{configuration, newAtomicState(), newAtomicBool(), function, newRetVal(), make(chan bool)}
}

func (process *Process)start() bool {
	switch process.state.get() {
	case ERROR:
		return false
	case RUNNING:
		return true
	case STOP:
		go process._startFunction() // starts thread.
		process.state.set(RUNNING)
		return true
	}
	return false
}

func (process *Process)_startFunction() {
	process.retVal.lock()
	process.function(process.configuration, process.terminator, process.retVal)
	process.retVal.unlock()
	process.state.set(STOP)
	process.endChannel <- true// scrivi su channel che hai finito.

}

func (process *Process)isAlive() bool {
	switch process.state.get() {
	case ERROR:
		return false
	case STOP:
		return false
	case RUNNING:
		return true
	}
	return false
}

func (process *Process)stop() bool {
	switch process.state.get() {
	case ERROR:
		return false
	case STOP:
		return true
	case RUNNING:
		process.terminator.set()
		return true
	}
	return false
}

func (process *Process)waitTermination() {
	switch process.state.get(){
	case RUNNING: <-process.endChannel // leggi dal channel per terminazione del thread.
	case STOP:
	case ERROR:
	}
}

//---------------MANAGER--------------

type Manager struct {
	processes []*Process // change data structure.
}

func newManager() Manager {
	return Manager{make([]*Process, 0)}
}

func (manager *Manager) addProcess(conf *ProcessConfiguration, worker WorkerFunction) int {
	process := newProcess(conf, worker)
	manager.processes = append(manager.processes, &process)
	return len(manager.processes) - 1 // index in the slice.
}

func (manager *Manager) startProcess(processId int) bool {
	return manager.processes[processId].start()
}

func (manager *Manager) startProcesses() bool {
	for i := 0; i < len(manager.processes); i++ {
		if (!manager.startProcess(i)) {
			return false
		}
	}
	return true
}

func (manager *Manager) stopProcess(processId int) bool {
	return manager.processes[processId].stop()
}

func (manager *Manager) stopProcesses() bool {
	for i := 0; i < len(manager.processes); i++ {
		if (!manager.stopProcess(i)) {
			return false
		}
	}
	return true
}

func (manager *Manager) waitProcessTermination(processId int) {
	manager.processes[processId].waitTermination()
}

func (manager *Manager) waitProcessesTermination() {
	for i := 0; i < len(manager.processes); i++ {
		manager.waitProcessTermination(i)
	}
}

func (manager *Manager) getProcessesNumber() int {
	return len(manager.processes)
}

//--------------FUNCTIONS IMPLEMENTING WORKER INTERFACE---------------

func write(conf *ProcessConfiguration, terminator *AtomicBool, retVal *RetVal) {
	fmt.Printf("\tave sono processo %s, scrivo nel channel\n", conf.name)
	conf.channel <- conf.name
	//for !terminator.get() {
	//retVal.set(1)
	//}
}

func read(conf *ProcessConfiguration, terminator *AtomicBool, retVal *RetVal) {
	fmt.Printf("\tave sono processo %s: ho letto: %s\n", conf.name, <-conf.channel)
	//for !terminator.get() {
	retVal.set(1)
	//}
}

func main() {
	var channel chan string = make(chan string, 2)
	var processManager Manager = newManager()
	var conf *ProcessConfiguration = newProcessConfiguration("pro" + strconv.Itoa(0), channel, make([]int, 0))
	processManager.addProcess(conf, write)
	conf = newProcessConfiguration("pro" + strconv.Itoa(1), channel, make([]int, 0))
	processManager.addProcess(conf, read)
	/*
	for i := 0; i < 2; i++ {
		var conf *ProcessConfiguration = newProcessConfiguration("pro" + strconv.Itoa(i), channel, make([]int, 0))
		processManager.addProcess(conf, BenOr)
	}
	*/
	//fmt.Printf("processes created, tot:%d\n", processManager.getProcessesNumber())
	processManager.startProcesses()
	fmt.Print("processes started");
	processManager.waitProcessesTermination()
	fmt.Print("arrivederci!\n")
}