package main

import (
	"consensus/processManager"
	"consensus/process"
	"fmt"
	"consensus/util"
	"time"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}


/*
 process number [2,15] N
 crashable processes [0, (N-1)/2] F
 range [2,15]
 mean delay  {0, 25, 50, 75, 100}
 delay variance {0, 10, 20}
 */

func min(n0 int, n1 int) int {
	var res int = n0
	if n1 < n0 {
		res = n1
	}
	return res
}

func main() {
	const MaxProcessesNumber int = 4
	const MaxCrashableProcesses int = (MaxProcessesNumber - 1) / 2
	const MaxRangeValue int = 10

	const MaxMeanDelay int = 100
	const DelayStep int = 25

	const MaxVariance int = 20
	const VarianceStep int = 10
	fmt.Println("processNumber\tcrashableProcesses\tmaxRangeVal\tdelayMean\tvariance")

	for processNumber := 2; processNumber <= MaxProcessesNumber; processNumber++ {
		for crashableProcesses := 0; crashableProcesses <= min((processNumber - 1) / 2, MaxCrashableProcesses); crashableProcesses++ {
			for maxRangeVal := 2; maxRangeVal <= MaxRangeValue; maxRangeVal++ {
				for delayMean := 0; delayMean <= MaxMeanDelay; delayMean += DelayStep {
					for variance := 0; variance <= MaxVariance; variance += VarianceStep {
						fmt.Printf("\n%d\t%d\t%d\t%d\t%d", processNumber, crashableProcesses, maxRangeVal, delayMean, variance)

						// init model.
						var manager processManager.Manager = processManager.NewManager(processNumber, delayMean, variance, crashableProcesses, maxRangeVal)
						var workers []process.WorkerFunction = make([]process.WorkerFunction, 0, processNumber)
						for i := 0; i < processNumber; i++ {
							workers = append(workers, BenOr)
						}
						manager.AddProcesses(workers)
						startTime := time.Now()
						manager.StartProcesses()

						// stop the processes
						for processId := 0; processId < crashableProcesses; processId++ {
							manager.StopProcess(processId)
						}
						var retVals []*util.RetVal = manager.WaitProcessesTermination()
						deltaTime := time.Now().Unix() - startTime.Unix()

						// check for errors in the execution.
						var decidedValue int = retVals[crashableProcesses].Get()
						for i := crashableProcesses + 1; i < len(retVals); i++ {
							if decidedValue != retVals[i].Get() {
								//fmt.Errorf("%d\t%d\t%d\t%d\t%d\n\tI valori di ritorno non sono consistenti: ", processNumber, crashableProcesses, maxRangeVal, delayMean, variance)
								fmt.Errorf("\n\tI valori di ritorno non sono consistenti: ", processNumber, crashableProcesses, maxRangeVal, delayMean, variance)
								for i := 0; i < len(retVals); i++ {
									fmt.Errorf("%d, ", retVals[i].Get())
								}
								fmt.Println()
							}
						}

						// log result.
						file, err := os.Create("../data/" + strconv.Itoa(processNumber) + "_" + strconv.Itoa(crashableProcesses) + "_" + strconv.Itoa(maxRangeVal) + "_" + strconv.Itoa(delayMean) + "_" + strconv.Itoa(variance) + ".csv")
						check(err)
						defer file.Close()
						file.WriteString(strconv.FormatInt(deltaTime, 10) + "," + strconv.Itoa(manager.GetSendedMessageNumber()) + "," + strconv.Itoa(manager.GetDeliveredMessageNumber()) + "," + strconv.Itoa(decidedValue) + "\n")
					}
				}
			}
		}
	}
}
