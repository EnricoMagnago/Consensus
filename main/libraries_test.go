package main

import (
	"testing"
	"consensus/processManager"
	"consensus/process"
	"consensus/util"
	"fmt"
	"consensus/channel"
)

func benOr(conf *process.ProcessConfiguration, terminator *util.AtomicBool, retVal *util.RetVal) {
	var ID int = conf.ProcessId

	var message *channel.Message = channel.NewMessage(ID, -1, channel.REPORT, 0, 1)

	conf.Channel.BroadcastSend(message)
	for i := 0; i < conf.ProcessesNumber - 1; i++ {
		var recMessage *channel.Message = nil
		for recMessage == nil {
			recMessage = conf.Channel.Deliver(ID)
			if recMessage != nil {
				fmt.Printf("\t%d) r: %d ; s: %d\n", ID, recMessage.GetReceiver(), recMessage.GetSender())
			}
		}
	}
	retVal.Set(ID)
}

func Test(t *testing.T) {
	var processNumber int = 20
	var delayMean int = 0
	var variance int = 0
	var manager processManager.Manager = processManager.NewManager(processNumber, delayMean, variance)
	var workers []process.WorkerFunction = make([]process.WorkerFunction, 0, processNumber)
	for i := 0; i < processNumber; i++ {
		workers = append(workers, benOr)
	}
	manager.AddProcesses(workers)
	manager.StartProcesses()
	//time.Sleep(time.Duration(5) * time.Second)
	manager.StopProcesses()
	manager.WaitProcessesTermination()

	for i := 0; i < processNumber; i++ {
		if manager.GetRetval(i).Get() != i {
			t.Errorf("found : %d expected: %d", manager.GetRetval(i).Get(), i)
		}
	}
}
