package util

import (
	"sync"
	"fmt"
)


// --------------RETVAL---------------

type RetVal struct {
	integer int
	mutex   sync.Mutex
}

func NewRetVal() *RetVal {
	return &RetVal{0, sync.Mutex{}}
}

func (retVal *RetVal)Lock() {
	retVal.mutex.Lock()
}

func (retVal *RetVal)Unlock() {
	retVal.mutex.Unlock()
}

func (retVal *RetVal)Set(newValue int) {
	retVal.integer = newValue
}


//--------------ATOMIC BOOL-------------

// atomic access to a bool variable.
type AtomicBool struct {
	boolean bool
	mutex   sync.RWMutex // read-write mutex.
}

func NewAtomicBool() *AtomicBool {
	return &AtomicBool{false, sync.RWMutex{}}
}

func (abool *AtomicBool) Set() {
	abool.mutex.Lock()
	abool.boolean = true
	abool.mutex.Unlock()
}

func (abool *AtomicBool)Get() bool {
	var res bool = false
	abool.mutex.RLock()
	res = abool.boolean
	abool.mutex.RUnlock()
	return res
}


//--------------ATOMIC STATE-------------

type State int

const (
	ERROR State = -1 + iota // ERROR = -1, STOP = 0, RUNNING = 1
	STOP
	RUNNING
)

// atomic access to a bool variable.
type AtomicState struct {
	value State
	mutex sync.RWMutex // read-write mutex.
}

func NewAtomicState() *AtomicState {
	return &AtomicState{STOP, sync.RWMutex{}}
}

func (astate *AtomicState) Set(new State) bool {
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

func (astate *AtomicState) Get() State {
	var res State
	astate.mutex.RLock()
	res = astate.value
	astate.mutex.RUnlock()
	return res
}