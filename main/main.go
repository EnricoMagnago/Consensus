package main

import (
	"consensus/processManager"
	"consensus/process"
	"fmt"
	"consensus/util"
)

type Command int

const (
	QUIT Command = iota
	START_ALL
	START
	STOP_ALL
	STOP
	WAIT_ALL
	WAIT
	ERROR
	LOG
)

func readCommand() Command {
	var tmp string = ""
	var done = false
	for !done {
		done = true
		fmt.Scanf("%s", &tmp)
		switch(tmp){
		case "quit": return QUIT
		case "startA": return START_ALL
		case "start": return START
		case "stopA": return STOP_ALL
		case "stop": return STOP
		case "waitA": return WAIT_ALL
		case "wait": return WAIT
		case "log" : return LOG
		default:
			fmt.Print("invalid command\n> ")
			done = false
		}
	}
	return ERROR
}

func readInt(n *int) {
	var tmp int = 0
	tmp, _ = fmt.Scanf("%d", n)
	for tmp == 0 {
		fmt.Print("\n\tplease insert an integer: ")
		tmp, _ = fmt.Scanf("%d", n)
	}
}

func main() {
	var processNumber int = 3
	var delayMean int = 2000
	var variance int = 1000
	var F int = 0
	var MaxVal int = 2

	fmt.Println("\tBen-Or protocol simulator\n\t\tby: Roberto Fellin, Enrico Magnago\n\n")

	fmt.Print("GLOBAL SETTINGS:\n\t- processes number: ")
	readInt(&processNumber)

	fmt.Print("\n\t- number of crashable processes: ")
	readInt(&F)
	fmt.Print("\n\t- number of decideable values: ")
	readInt(&MaxVal)

	fmt.Print("\n\t- channel mean delay: ")
	readInt(&delayMean)
	fmt.Print("\n\t- channel delay variance: ")
	readInt(&variance)



	var manager processManager.Manager = processManager.NewManager(processNumber, delayMean, variance, F, MaxVal)
	var workers []process.WorkerFunction = make([]process.WorkerFunction, 0, processNumber)
	for i := 0; i < processNumber; i++ {
		workers = append(workers, BenOr)
	}
	manager.AddProcesses(workers)

	var command Command = ERROR
	for command != QUIT {
		fmt.Print("> ")
		command = readCommand()
		fmt.Println()
		switch(command){
		case QUIT:
			manager.StopProcesses()
		case START_ALL:
			manager.StartProcesses()
			fmt.Printf("All processes started.\n");
		case START:
			var id int = -1
			for id >= processNumber && id < 0 {
				fmt.Printf("id of the process to start [0;%d]: ", processNumber - 1)
				readInt(&id)
			}
			manager.StartProcess(id)
		case STOP_ALL:
			manager.StopProcesses()
			fmt.Printf("All processes stopped.\n");
		case STOP:
			var id int = -1
			for id >= processNumber || id < 0 {
				fmt.Printf("id of the process to stop [0;%d]: ", processNumber - 1)
				readInt(&id)
			}
			manager.StopProcess(id)
		case WAIT_ALL:
			var retVals []*util.RetVal = manager.WaitProcessesTermination()
			for i := 0; i < len(retVals); i++ {
				fmt.Printf("process %d returned: %d\n", i, retVals[i].Get());
			}

		case WAIT:
			var id int = -1
			for id >= processNumber || id < 0 {
				fmt.Printf("id of the process to wait [0;%d]: ", processNumber - 1)
				readInt(&id)
			}
			manager.GetRetval(id)
		case LOG:
			manager.LogState()
		}
	}
}
