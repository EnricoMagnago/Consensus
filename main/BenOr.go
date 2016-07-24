package main

import (
	"consensus/process"
	"consensus/util"
	"math/rand"
	"time"
	"consensus/channel"
	"fmt"
)

func newRandEstimate(maxVal int) int {
	return int(rand.Int()) % maxVal
}

func BenOr(conf *process.ProcessConfiguration, terminator *util.AtomicBool, retVal *util.RetVal) {
	rand.Seed(int64(time.Now().Nanosecond()))
	const maxVal int = 10 //from 0 to maxVal-1
	const F int = 0

	// process state init
	var ID int = conf.ProcessId
	var N int = conf.ProcessesNumber
	var ROUND int = 0
	var EST int = newRandEstimate(maxVal)
	var DECIDED bool = false

	for !DECIDED && !terminator.Get() {
		ROUND += 1
		//fmt.Printf("%d) est: %d; round: %d\n", ID, EST, ROUND);
		var message *channel.Message = channel.NewMessage(ID, -1, channel.REPORT, ROUND, EST)
		var broadCastres bool = conf.Channel.BroadcastSend(message)
		if !broadCastres {
			fmt.Errorf("ERROR BenOr in the broadcast send of report messages")
		}
		fmt.Printf("%d - waiting Majority\n", ID)
		var majority = waitMajority(N, F, conf.Channel, &ROUND, ID, maxVal, terminator)
		fmt.Printf("%d) found Majority\n", ID)
		switch(majority){
		case -1:
		//fmt.Printf("%d) no majority\n", ID)
		default:
		//fmt.Printf("%d) majority found on %d\n", ID, majority)
		}
		//DECIDED = true
		//majority can be -1 -> no majority.

		var msgProp *channel.Message = channel.NewMessage(ID, -1, channel.PROPOSAL, ROUND, majority)
		broadCastres = conf.Channel.BroadcastSend(msgProp) // can be -1
		if !broadCastres {
			fmt.Errorf("ERROR BenOr in the broadcast send of proposal messages")
		}
		fmt.Printf("%d - waiting Proposal\n", ID)
		var proposalRet []int = waitProposal(N, F, conf.Channel, &ROUND, ID, maxVal, terminator) //wait a proposal with a setted estimate (!= -1), returns the value if present, -1 if not
		fmt.Printf("%d) found Proposal\n", ID)
		var majorityEst int = proposalRet[0] // value of the found majority, -1 if not found.
		var counter int = proposalRet[1] // number of proposal received with majorityEst value.

		if majorityEst != -1 {
			EST = majorityEst
		} else {
			EST = newRandEstimate(maxVal)
		}
		if counter > F {
			DECIDED = true
		}
		//decide -> wait N-F proposal with the same value.
	}

	//fmt.Printf("%d) decided: %d, round: %d\n", ID, EST, ROUND);
	retVal.Set(EST)
}


/**
 * @param n number processes in the protocol.
 * @param f number of processes that can crash.
 * @param channel communication channel to read.
 * @param round current round of the protocol, ignore messages with different round.
 * @param id of the calling process.
 * @param maxVal the proposed values will be from 0 to maxVal-1 (included).
 * @param terminator termination requested by outside.
 * @return number of proposal message on estimate (which is setted properly).
 */
func waitProposal(n int, f int, chann *channel.Channel, round *int, processId int, maxVal int, terminator *util.AtomicBool) []int {
	var messageNumber int = 0
	var majorityCounter int = 0
	var majorityEst int = -1 // which value has the majority. -1 no majority.
	// exit when a majority is found or when delivered n-f messages with the correct round and type.
	for messageNumber < (n - f) && !terminator.Get() {
		var message *channel.Message = chann.Deliver(processId)
		if message != nil &&  message.GetMessageType() == channel.PROPOSAL {
			messageNumber++
			fmt.Printf("%d) proposal received %d from: %d\n", processId, message.GetEstimate(), message.GetSender())

			// message with different round or wrong type.
			if message.GetRound() != *round {
				if message.GetRound() > *round {
					fmt.Printf("%dp) round update\n", processId)
					*round = message.GetRound() // round update
					// restart the count of the messages (new round) and keep going.
				}
			}

			// the sender has seen a majority.
			if message.GetEstimate() != -1 {
				// first estimate != -1 seen by the current process.
				if majorityEst == -1 {
					majorityEst = message.GetEstimate() // set the value of the majority
					majorityCounter = 1 // start counting.
				} else {
					// there can not be 2 majorities.
					if majorityEst != message.GetEstimate() {
						fmt.Errorf("%dp) received 2 different estimate in proposal\n", processId)
						var res []int = make([]int, 2)
						res[0] = -1
						res[1] = -1
					}
					majorityCounter++
				}
			} // the sender has not seen a majority	-> do nothing.
		}
	}
	//fmt.Printf("%dp) majority counter: %d; est:%d\n", processId, majorityCounter, majorityEst)
	var res []int = make([]int, 2)
	res[0] = majorityEst
	res[1] = majorityCounter
	return res
}

/**
 * @param n number processes in the protocol.
 * @param f number of processes that can crash.
 * @param channel communication channel to read.
 * @param round current round of the protocol, ignore messages with different round.
 * @param id of the calling process.
 * @param maxVal the proposed values will be from 0 to maxVal-1 (included).
 * @param terminator termination requested by outside.
 * @return -1: no majority, k: majority found on k (k>=0)
 */
func waitMajority(n int, f int, chann *channel.Channel, round *int, processId int, maxVal int, terminator *util.AtomicBool) int {
	var counters []int = make([]int, maxVal, maxVal)
	var messageNumber int = 0
	for i := 0; i < maxVal; i++ {
		counters[i] = 0
	}
	var res int = -3
	// exit when a majority is found or when delivered n-f messages with the correct round and type.
	for res == -3 && !terminator.Get() {
		var message *channel.Message = chann.Deliver(processId)
		if message != nil &&  message.GetMessageType() == channel.REPORT {
			fmt.Printf("%d) received %d from: %d\n", processId, message.GetEstimate(), message.GetSender())
			// message with different round or wrong type.
			if message.GetRound() != *round {
				if message.GetRound() > *round {
					fmt.Printf("%d) round update\n", processId)
					*round = message.GetRound() // round update
					// restart the count of the messages (new round) and keep going.
					for i := 0; i < maxVal; i++ {
						counters[i] = 0
					}
				}
			}
			var est int = message.GetEstimate()
			counters[est]++
			// majority found
			if counters[est] >= (n + 1) / 2 {
				res = est
			} else {
				messageNumber++
				// n-f messages, still no majority.
				if messageNumber >= n - f {
					res = -1
				}
			}
		}
	}
	//fmt.Printf("%d) ret: %d\n", processId, res)
	return res
}
