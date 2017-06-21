package channel

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
)


//--------MESSAGE----------

type MessageType int

const (
	REPORT MessageType = iota
	PROPOSAL
	DECIDE
)

type Message struct {
	senderId    int
	receiverId  int
	messageType MessageType
	round       int // protocol round.
	estimate    int // actual payload.
}

/**
 * Creates a new message with the specified parameters, just a wrapper.
 */
func NewMessage(senderId int, receiverId int, messageType MessageType, round int, estimate int) *Message {
	return &Message{senderId, receiverId, messageType, round, estimate}
}

func (message *Message)GetSender() int {
	return message.senderId
}
func (message *Message)GetReceiver() int {
	return message.receiverId
}

func (message *Message)GetMessageType() MessageType {
	return message.messageType
}

func (message *Message)GetRound() int {
	return message.round
}

func (message *Message)GetEstimate() int {
	return message.estimate
}

// struct that adds to the message the time at which it can be delivered. Should never appear outside.
type MessageDeliveryTime struct {
	message      *Message
	deliveryTime *time.Time
}

func newMessageDeliveryTime(message *Message, deliveryTime *time.Time) *MessageDeliveryTime {
	return &MessageDeliveryTime{message, deliveryTime}
}

//------MESSAGESQUEUE-------------
/*
 * each process has his own queue of messages.
 * the messages are read from the shared channel and put on the corresponding queue (identifier of the process is the index of the queue).
 */
type MessagesQueue struct {
	queue []*MessageDeliveryTime
	mutex sync.Mutex
}

func NewMessagesQueue() *MessagesQueue {
	return &MessagesQueue{make([]*MessageDeliveryTime, 0), sync.Mutex{}}
}

/*
 * returns the number of messages inside the queue, private method.
 */
func (messageQueue *MessagesQueue) size() int {
	var res int = 0
	messageQueue.mutex.Lock()
	res = len(messageQueue.queue)
	messageQueue.mutex.Unlock()
	return res
}

/*
 * returns true only if there is at least 1 message in the queue.
 */
func (messageQueue *MessagesQueue) IsEmpty() bool {
	return messageQueue.size() == 0
}

/*
 * append a message in the queue.
 */
func (messageQueue *MessagesQueue) Add(message *MessageDeliveryTime) {
	messageQueue.mutex.Lock()
	messageQueue.queue = append(messageQueue.queue, message)
	messageQueue.mutex.Unlock()
}

/*
 * returns the first message in the queue, nil if the queue is empty.
 */
func (messageQueue *MessagesQueue) Pop() *Message {
	if messageQueue.IsEmpty() {
		return nil
	}
	var message *Message = nil
	messageQueue.mutex.Lock()
	if messageQueue.queue[0].deliveryTime.Before(time.Now()) {
		message = messageQueue.queue[0].message
		messageQueue.queue = messageQueue.queue[1:] // remove from the queue.
	} else {
		//fmt.Println("sorry you have to wait")
	}
	messageQueue.mutex.Unlock()
	return message
}

//---------------CHANNEL--------------

type Channel struct {
	mutex                sync.Mutex
	processesNumber      int             // total number of processes sharing the Channel.
	capacity             int             // fixed size queue.
	gochannel            chan Message    // actual channel of communication.
	messagesBuffer       []MessagesQueue // buffer of messages sent but not delivered.
	mean                 int64           // mean delay.
	variance             int64           // variance of the delay.
					     // meta-data on the channel, useful for usage stats.
	sendMessagesCount    int
	deliverMessagesCount int
	nilMessagesCount     int
	lastClient           int
}

func NewChannel(processNumber int, mean int, variance int) *Channel {
	rand.Seed(int64(time.Now().Nanosecond()))
	var channelCapacity int = processNumber * processNumber

	var channel Channel = Channel{sync.Mutex{}, processNumber, channelCapacity, make(chan Message, channelCapacity), make([]MessagesQueue, processNumber), int64(mean), int64(variance), 0, 0, 0, -1}
	for i := 0; i < processNumber; i++ {
		channel.messagesBuffer[i] = *NewMessagesQueue()
	}
	return &channel
}
/*
 * returns a gaussian generated delay generated using the given mean and variance.
 */
func (channel *Channel)generateDeliveryTime() *time.Time {
	var randInt int64 = int64(rand.NormFloat64() * float64(channel.variance) + float64(channel.mean))
	var delay time.Duration = time.Duration(randInt) * time.Millisecond
	var receiveTimestamp time.Time = time.Now().Add(delay)
	return &receiveTimestamp
}
// wrapper on the inner gochannel.
func (channel *Channel) isEmpty() bool {
	return len(channel.gochannel) == 0
}

// wrapper on the inner gochannel.
func (channel *Channel) isFull() bool {
	return len(channel.gochannel) == channel.capacity
}

/*
 * reads messages from the gochannel until empty or a message for the process with the specified id is found.
 * if id = -2, reads all the messages from the channel.
 */
func (channel *Channel) receive(id int) {
	channel.lastClient = id
	var lastId int = -1
	//channel.mutex.Lock()
	for !channel.isEmpty() && lastId != id {
		var message Message = <-channel.gochannel // read message from go channel.
		lastId = message.GetReceiver() // who is the receiver.
		channel.messagesBuffer[lastId].Add(newMessageDeliveryTime(&message, channel.generateDeliveryTime()))
	}
	//channel.mutex.Unlock()
}

// wrapper on receive.
func (channel *Channel) receiveAll() {
	channel.receive(-2)
}

/*
 * send a new message through the channel.
 * returns false if the message is malformed.
 */
func (channel *Channel) Send(message *Message) bool {
	channel.lastClient = message.senderId
	channel.mutex.Lock()
	if channel.isFull() {
		channel.receiveAll()
	}

	if channel.isFull() {
		fmt.Errorf("ERROR: Channel.send: can not __receive messages, channel still full.")
		channel.mutex.Unlock()
		return false
	}

	if message.GetSender() > -1 && message.GetSender() < channel.processesNumber &&  message.GetReceiver() > -1 && message.GetReceiver() < channel.processesNumber {
		//fmt.Printf("%d sending to %d\n", message.GetSender(), message.GetReceiver())

		channel.gochannel <- *message
		channel.sendMessagesCount++
		channel.mutex.Unlock()
		return true
	}
	channel.mutex.Unlock()
	fmt.Errorf("WARNING: wrong process id, can not send message")
	return false
}

/*
 * send the message to all the processes that share the channel.
 */
func (channel *Channel) BroadcastSend(message *Message) bool {
	var res bool = true
	for i := 0; i < channel.processesNumber; i++ {
		// simply perform a Send to every process. No need to handle if the process fails while sending.
		message.receiverId = i
		res = res && channel.Send(message)
		if !res {
			fmt.Errorf("ERROR broadcast send to %d failed", i)
		}
	}
	return res
}

/*
 * deliver a message for the given process identifier.
 * returns nil if there is no message.
 */
func (channel *Channel) Deliver(processId int) *Message {
	channel.lastClient = processId
	if channel.messagesBuffer[processId].IsEmpty() {
		channel.mutex.Lock()
		channel.receive(processId)
		channel.mutex.Unlock()
	}

	var res *Message = channel.messagesBuffer[processId].Pop()
	if res != nil {
		channel.deliverMessagesCount++
		//fmt.Printf("%d received from %d\n", res.GetReceiver(), res.GetSender())
	} else {
		//fmt.Printf("nil message to %d\n", processId)
		channel.nilMessagesCount++
	}
	return res
}

func (channel *Channel) GetSendedMessagesNumber() int {
	return channel.sendMessagesCount
}

func (channel *Channel) GetDeliveredMessagesNumber() int {
	return channel.deliverMessagesCount
}

func (channel *Channel) PrintState() {
	fmt.Printf("last client %d\nsended: %d; delivered: %d; nilReplys: %d\nQueue sizes:", channel.lastClient, channel.GetSendedMessagesNumber(), channel.GetDeliveredMessagesNumber(), channel.nilMessagesCount)
	for process := 0; process < channel.processesNumber; process++ {
		fmt.Printf("\n\tprocess %d : %d", process, channel.messagesBuffer[process].size())
	}

	fmt.Print("\n")
}