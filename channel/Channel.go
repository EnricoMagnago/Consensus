package channel

import "fmt"


//--------MESSAGE----------

type Message struct {
	senderId   int
	receiverId int
	message    string
}

//------MESSAGESQUEUE-------------

type MessagesQueue struct {
	queue []*Message
}

func NewMessagesQueue() *MessagesQueue {
	return &MessagesQueue{make([]*Message, 0)}
}

func (messageQueue *MessagesQueue) IsEmpty() bool {
	return len(messageQueue.queue) == 0
}

func (messageQueue *MessagesQueue) Add(message *Message) {
	messageQueue.queue = append(messageQueue.queue, message)
}
func (messageQueue *MessagesQueue) Pop() *Message {
	if messageQueue.IsEmpty() {
		return nil
	}
	var message *Message = messageQueue.queue[0]
	messageQueue.queue = messageQueue.queue[1:]
	return message
}

//---------------CHANNEL--------------

type Channel struct {
	processesNumber int
	capacity        int
	gochannel       chan Message
	messagesBuffer  []MessagesQueue
}

func NewChannel(processNumber int, channelCapacity int) *Channel {
	if channelCapacity <= 0 {
		channelCapacity = processNumber * processNumber
	}
	var channel Channel = Channel{processNumber, channelCapacity, make(chan Message, channelCapacity), make([]MessagesQueue, processNumber)}
	for i := 0; i < processNumber; i++ {
		channel.messagesBuffer[i] = *NewMessagesQueue()
	}
	return &channel
}

func (channel *Channel) isEmpty() bool {
	return len(channel.gochannel) == 0
}

func (channel *Channel) isFull() bool {
	return len(channel.gochannel) == channel.capacity
}
/**
reads messages from the gochannel until empty or a message for the process with the specified id is found.
if id = -2, reads all the messages from the channel.
 */
func (channel *Channel) receive(id int) {
	var lastId int = -1
	for !channel.isEmpty() && lastId != id {
		var message Message = <-channel.gochannel
		lastId = message.receiverId
		channel.messagesBuffer[lastId].Add(&message)
	}
}

func (channel *Channel) receiveAll() {
	channel.receive(-2)
}

func (channel *Channel) send(message *Message) {
	if channel.isFull() {
		channel.receiveAll()
	}
	if channel.isFull() {
		fmt.Errorf("ERROR: Channel.send: can not __receive messages, channel still full.")
		return
	}
	channel.gochannel <- *message
}

func (channel *Channel) Deliver(processId int) *Message {
	if channel.messagesBuffer[processId].IsEmpty() {
		channel.receive(processId)
	}
	var res *Message = channel.messagesBuffer[processId].Pop()
	return res
}

