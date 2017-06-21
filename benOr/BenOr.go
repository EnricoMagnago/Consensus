package benOr

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
	var maxVal int = conf.MaxVal //from 0 to maxVal-1
	var F int = conf.F

	// process state init
	var ID int = conf.ProcessId
	var N int = conf.ProcessesNumber
	var ROUND int = 0
	var EST int = newRandEstimate(maxVal)
	var DECIDED bool = false
	// if true : broadcast your decision, false: avoid flooding.
	var send_decide bool = true

	for !DECIDED && !terminator.Get() {
		ROUND += 1
		//fmt.Printf("%d) est: %d; round: %d\n", ID, EST, ROUND);
		var message *channel.Message = channel.NewMessage(ID, -1, channel.REPORT, ROUND, EST)
		var broadCasted bool = conf.Channel.BroadcastSend(message)
		if !broadCasted {
			fmt.Errorf("ERROR BenOr in the broadcast send of report messages")
		}
		//fmt.Printf("%d - waiting Majority\n", ID)
		var majority_arr = waitMajority(N, F, conf.Channel, &ROUND, ID, maxVal, terminator)
		if (majority_arr[1] == 1) {
			// decide message received.
			EST = majority_arr[0]
			send_decide = false // avoid flooding.
			DECIDED = true
		} else {
			// process majority.
			//fmt.Printf("%d) found Majority\n", ID)
			var majority int = majority_arr[0]
			switch(majority){
			case -3:
				// terminated by request, exit with -2
				EST = -2
				DECIDED = true
			default:
				// -1: majority not found else majority found
				//majority can be -1 -> no majority.

				var msgProp *channel.Message = channel.NewMessage(ID, -1, channel.PROPOSAL, ROUND, majority)
				broadCasted = conf.Channel.BroadcastSend(msgProp) // can be -1
				if !broadCasted {
					fmt.Errorf("ERROR BenOr in the broadcast send of proposal messages")
				}
				//fmt.Printf("%d - waiting Proposal\n", ID)
				var proposalRet []int = waitProposal(N, F, conf.Channel, &ROUND, ID, maxVal, terminator) //wait a proposal with a setted estimate (!= -1), returns the value if present, -1 if not
				//fmt.Printf("%d) found Proposal\n", ID)

				var majorityEst int = proposalRet[0] // value of the found majority, -1 if not found.
				var counter int = proposalRet[1] // number of proposal received with majorityEst value.
				if terminator.Get() || (counter == -1 && majorityEst == -1) {
					// terminated by request or error in waitProposal, exit with -2
					EST = -2
					DECIDED = true
				} else {
					if counter == -1 {
						// received a DECIDE message
						EST = majorityEst
						DECIDED = true
						send_decide = false
					} else {
						if majorityEst != -1 {
							EST = majorityEst
						} else {
							EST = newRandEstimate(maxVal)
						}
						if counter > F {
							DECIDED = true
						}
					}
				}
			}
			//decide -> wait N-F proposal with the same value.
		}
	}

	if EST > -1 && send_decide {
		// set DECIDE message, tell all others a decision has been reached.
		var msgDecide *channel.Message = channel.NewMessage(ID, -1, channel.DECIDE, ROUND, EST)
		if !conf.Channel.BroadcastSend(msgDecide) {
			fmt.Errorf("ERROR BenOr in the broadcast send of decide messages")
		}
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
 * @return number of proposal message on estimate (which is setted properly); if counter is set to -1 it means the estimate has been decided.
 */
func waitProposal(n int, f int, chann *channel.Channel, round *int, processId int, maxVal int, terminator *util.AtomicBool) []int {
	var messageNumber int = 0
	var majorityCounter int = 0
	var majorityEst int = -1 // which value has the majority. -1 no majority.
	// exit when a majority is found or when delivered n-f messages with the correct round and type.
	for messageNumber < (n - f) && !terminator.Get() {
		var message *channel.Message = chann.Deliver(processId)
		if message != nil {
			switch message.GetMessageType(){
			case channel.DECIDE:
				if majorityEst != -1 && majorityEst != message.GetEstimate() {
					fmt.Errorf("%dp) received a decide with a different value\n", processId)
					var res []int = make([]int, 2)
					res[0] = -1
					res[1] = -1
				}
				majorityEst = message.GetEstimate()
				majorityCounter = -1
				messageNumber = n // exit loop
			//break
			case channel.PROPOSAL:
				messageNumber++
				//fmt.Printf("%d) (wProposal) received %d from: %d\n", processId, message.GetEstimate(), message.GetSender())

				// message with different round or wrong type.
				if message.GetRound() != *round {
					if message.GetRound() > *round {
						//fmt.Printf("%dp) round update\n", processId)
						*round = message.GetRound() // round update
						// restart the count of the messages (new round) and keep going.
						messageNumber = 0
						majorityCounter = 0
						majorityEst = -1
						// send a copy of the message.
						var message *channel.Message = channel.NewMessage(processId, -1, channel.PROPOSAL, *round, message.GetEstimate())
						var broadCasted bool = chann.BroadcastSend(message)
						if !broadCasted {
							fmt.Errorf("ERROR BenOr in the broadcast send of report messages")
						}
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
			//break
			case channel.REPORT:
				if message.GetRound() > *round {
					chann.Send(message)
				}
			//break
			default:
				fmt.Errorf("Unknown message type\n")
			}
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
 * @param terminator termination requested from outside.
 * @return -1: no majority, k: majority found on k (k>=0)
 */
func waitMajority(n int, f int, chann *channel.Channel, round *int, processId int, maxVal int, terminator *util.AtomicBool) []int {
	var counters []int = make([]int, maxVal, maxVal)
	var messageNumber int = 0
	for i := 0; i < maxVal; i++ {
		counters[i] = 0
	}
	var is_decide int = 0
	var res int = -3
	// exit when a majority is found or when delivered n-f messages with the correct round and type.
	for res == -3 && !terminator.Get() {
		var message *channel.Message = chann.Deliver(processId)
		if message != nil {
			switch message.GetMessageType(){
			case channel.DECIDE:
				res = message.GetEstimate()
				is_decide = 1
			case channel.REPORT:
				//fmt.Printf("%d) (wMajority) received %d from: %d\n", processId, message.GetEstimate(), message.GetSender())
				// message with different round or wrong type.
				if message.GetRound() != *round {
					if message.GetRound() > *round {
						//fmt.Printf("%d) round update\n", processId)
						*round = message.GetRound() // round update
						// restart the count of the messages (new round) and keep going.
						for i := 0; i < maxVal; i++ {
							counters[i] = 0
						}
						// send a copy of the message.
						var message *channel.Message = channel.NewMessage(processId, -1, channel.REPORT, *round, message.GetEstimate())
						var broadCasted bool = chann.BroadcastSend(message)
						if !broadCasted {
							fmt.Errorf("ERROR BenOr in the broadcast send of report messages")
						}

					}
				}
				var est int = message.GetEstimate()
				counters[est]++
				// majority found
				if counters[est] >= (n + 1) / 2 {
					res = est // change res, exit loop.
				} else {
					messageNumber++
					// n-f messages, still no majority.
					if messageNumber >= n - f {
						res = -1 // change res: -1 signals no majority found.
					}
				}
			//break
			case channel.PROPOSAL:
				if message.GetRound() >= *round {
					chann.Send(message) // re-enqueue
				}
			//break
			default:
				fmt.Errorf("Unknown message type\n")
			}
		}
	}
	//fmt.Printf("%d) ret: %d\n", processId, res)
	var tmp []int = make([]int, 2)
	tmp[0] = res
	tmp[1] = is_decide
	return tmp
}
