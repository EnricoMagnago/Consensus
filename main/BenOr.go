package main

import (
	"consensus/process"
	"consensus/util"
	"math/rand"
	"time"
	"consensus/channel"
)

func BenOr(conf *process.ProcessConfiguration, terminator *util.AtomicBool, retVal *util.RetVal) {
	rand.Seed(int64(time.Now().Nanosecond()))
	const maxVal int = 10 //from 0 to maxVal-1
	const F int = 1

	// process state init
	var ID int = conf.ProcessId
	var N int = conf.ProcessesNumber
	var ROUND int = 0
	var EST int = int(rand.Int()) % maxVal
	var DECIDED bool = false

	for !DECIDED && !terminator.Get() {
		ROUND += 1
		var message *channel.Message = channel.NewMessage(ID, -1, channel.REPORT, ROUND, EST)
		conf.Channel.BroadcastSend(message)

		//waitMajority() //wait n/2 REPORT messages with the same estimate value and returns it, if not present returns -1.
		//conf.Channel.Broadcast(majorityValue (can be -1))

		//EST, counter = waitProposal() //wait a proposal with a setted estimate (!= -1), returns the value if present, -1 if not
		//if EST = - 1 : EST = int(rand.Int()) % maxVal
		//if counter > F : decide

		//decide -> wait N-F proposal with the same value.
		/*
		var counters []int = make([]int, maxVal, maxVal)
		for i := 0; i < maxVal; i++ {
			counters[i] = 0
		}

		var message *channel.Message = nil
		var foundMajority bool = false
		var messagesNumber int = 0
		for !foundMajority && !terminator.Get() && messagesNumber < N - F {
			// wait a message
			for message == nil && !terminator.Get() {
				message = conf.Channel.Deliver(ID)
			}
			// check message
			if message.GetRound() == ROUND && message.GetMessageType() == channel.REPORT {
				messagesNumber += 1
				counters[message.GetEstimate()] += 1
				// found a majority
				if counters[message.GetEstimate()] >= N / 2 + 1 {
					var propose *channel.Message = channel.NewMessage(ID, -1, channel.PROPOSAL, ROUND, message.GetEstimate())
					conf.Channel.BroadcastSend(propose)
					foundMajority = true
				}
			}
			message = nil
		}
		if !foundMajority {
			var propose *channel.Message = channel.NewMessage(ID, -1, channel.PROPOSAL, ROUND, -1)
			conf.Channel.BroadcastSend(propose)
		}*/
	}

	retVal.Set(ID)
}
