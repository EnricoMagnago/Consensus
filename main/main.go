package main

import (
	"consensus/processManager"
	"consensus/process"
	"fmt"
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

	fmt.Println("\tBen-Or protocol simulator\n\t\tby: Roberto Fellin, Enrico Magnago\n\n")

	fmt.Print("GLOBAL SETTINGS:\n\t- processes number: ")
	readInt(&processNumber)

	fmt.Print("\n\t- channel mean delay: ")
	readInt(&delayMean)
	fmt.Print("\n\t- channel delay variance: ")
	readInt(&variance)

	var manager processManager.Manager = processManager.NewManager(processNumber, delayMean, variance)
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
		case START:
			var id int = -1
			for id >= processNumber && id < 0 {
				fmt.Printf("id of the process to start [0;%d]: ", processNumber - 1)
				readInt(&id)
			}
			manager.StartProcess(id)
		case STOP_ALL:
			manager.StopProcesses()
		case STOP:
			var id int = -1
			for id >= processNumber && id < 0 {
				fmt.Printf("id of the process to stop [0;%d]: ", processNumber - 1)
				readInt(&id)
			}
			manager.StopProcess(id)
		case WAIT_ALL:
			manager.WaitProcessesTermination()
		case WAIT:
			var id int = -1
			for id >= processNumber && id < 0 {
				fmt.Printf("id of the process to wait [0;%d]: ", processNumber - 1)
				readInt(&id)
			}
			manager.GetRetval(id)
		}
	}
	manager.StartProcesses()
	//time.Sleep(time.Duration(5) * time.Second)
	manager.WaitProcessesTermination()

}
