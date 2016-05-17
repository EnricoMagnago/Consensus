package process

import (
	"consensus/util"
)

// interface for the function executed by threads, must periodically check if AtomicBool is setted, if it's the case the thread must terminate.
type WorkerFunction func(*ProcessConfiguration, *util.AtomicBool, *util.RetVal)  //function type, the last int is the retVal.


// parameters to the function, should not be modified after the start (concurrency).
type ProcessConfiguration struct {
	name         string
	channel      chan string
	processesIds []int
}


// ----------------PROCESS--------------

type Process struct {
	configuration *ProcessConfiguration
	state         *util.AtomicState
	terminator    *util.AtomicBool
	function      WorkerFunction
	retVal        *util.RetVal
	endChannel    chan bool
}

func NewProcess(configuration *ProcessConfiguration, function WorkerFunction) Process {
	return Process{configuration, util.NewAtomicState(), util.NewAtomicBool(), function, util.NewRetVal(), make(chan bool)}
}

func (process *Process)Start() bool {
	switch process.state.Get() {
	case util.ERROR:
		return false
	case util.RUNNING:
		return true
	case util.STOP:
		go process.startFunction() // starts thread.
		process.state.Set(util.RUNNING)
		return true
	}
	return false
}

func (process *Process)startFunction() {
	process.retVal.Lock()
	process.function(process.configuration, process.terminator, process.retVal)
	process.retVal.Unlock()
	process.state.Set(util.STOP)
	process.endChannel <- true// scrivi su channel che hai finito.

}

func (process *Process)IsAlive() bool {
	switch process.state.Get() {
	case util.ERROR:
		return false
	case util.STOP:
		return false
	case util.RUNNING:
		return true
	}
	return false
}

func (process *Process)Stop() bool {
	switch process.state.Get() {
	case util.ERROR:
		return false
	case util.STOP:
		return true
	case util.RUNNING:
		process.terminator.Set()
		return true
	}
	return false
}

func (process *Process)WaitTermination() {
	switch process.state.Get(){
	case util.RUNNING: <-process.endChannel // leggi dal channel per terminazione del thread.
	case util.STOP:
	case util.ERROR:
	}
}