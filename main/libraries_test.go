package main

import (
	"consensus/processManager"
	"consensus/process"
	"strconv"
)

func main(){
	var processNumber int = 3
	var delayMean int = 100
	var variance int = 10
	var manager processManager.Manager = processManager.NewManager(processNumber, delayMean, variance)
	var configurations [processNumber]process.ProcessConfiguration
	for i := 0; i < processNumber; i++{
		configurations[i] = process.ProcessConfiguration{"process" + strconv.Itoa(i), manager.channel}
	}
}
