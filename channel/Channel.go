package channel

import (
	"fmt"
	"math/rand"
	"time"
)


//--------MESSAGE----------

type MessageType int

const (
	REPORT MessageType = iota
	PROPOSAL
)

type Message struct {
	senderId    int
	receiverId  int
	messageType MessageType
	round       int
	estimate    int

}

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

type MessageDeliveryTime struct {
	message      *Message
	deliveryTime *time.Time
}

func newMessageDeliveryTime(message *Message, deliveryTime *time.Time) *MessageDeliveryTime {
	return &MessageDeliveryTime{message, deliveryTime}
}

//------MESSAGESQUEUE-------------

type MessagesQueue struct {
	queue []*MessageDeliveryTime
}

func NewMessagesQueue() *MessagesQueue {
	return &MessagesQueue{make([]*MessageDeliveryTime, 0)}
}

func (messageQueue *MessagesQueue) IsEmpty() bool {
	return len(messageQueue.queue) == 0
}

func (messageQueue *MessagesQueue) Add(message *MessageDeliveryTime) {
	messageQueue.queue = append(messageQueue.queue, message)
}
func (messageQueue *MessagesQueue) Pop() *Message {
	if messageQueue.IsEmpty() {
		return nil
	}
	var message *Message = nil
	if messageQueue.queue[0].deliveryTime.Before(time.Now()) {
		message = messageQueue.queue[0].message
		messageQueue.queue = messageQueue.queue[1:] // remove from the queue.
	}else{
		fmt.Println("sorry you have to wait")
	}
	return message
}

//---------------CHANNEL--------------

type Channel struct {
	processesNumber int
	capacity        int
	gochannel       chan Message
	messagesBuffer  []MessagesQueue
	mean            int64
	variance        int64
	count int
	countSend int
}

func NewChannel(processNumber int, mean int, variance int) *Channel {
	rand.Seed(int64(time.Now().Nanosecond()))
	var channelCapacity int = processNumber * processNumber

	var channel Channel = Channel{processNumber, channelCapacity, make(chan Message, channelCapacity), make([]MessagesQueue, processNumber), int64(mean), int64(variance),0,0}
	for i := 0; i < processNumber; i++ {
		channel.messagesBuffer[i] = *NewMessagesQueue()
	}
	return &channel
}

func (channel *Channel)generateDeliveryTime() *time.Time {
	var randInt int64 = int64(rand.NormFloat64() * float64(channel.variance) + float64(channel.mean))
	var delay time.Duration = time.Duration(randInt) * time.Millisecond
	var recieveTimestamp time.Time = time.Now().Add(delay)
	return &recieveTimestamp
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
		fmt.Printf("Rec %d",id);
		channel.count++
		var message Message = <-channel.gochannel
		lastId = message.GetReceiver()
		channel.messagesBuffer[lastId].Add(newMessageDeliveryTime(&message, channel.generateDeliveryTime()))
	}


}

func (channel *Channel) receiveAll() {
	channel.receive(-2)
}

func (channel *Channel) Send(message *Message) bool {
	if channel.isFull() {
		channel.receiveAll()
	}
	if channel.isFull() {
		fmt.Errorf("ERROR: Channel.send: can not __receive messages, channel still full.")
		return false
	}
	if message.GetSender() > -1 && message.GetSender() < channel.processesNumber &&  message.GetReceiver() > -1 && message.GetReceiver() < channel.processesNumber {
		channel.gochannel <- *message
		channel.countSend++
		return true
	}
	fmt.Errorf("WARNING: wrong process id, can not send message")
	return false
}

func (channel *Channel) BroadcastSend(message *Message) bool {
	var res bool = true
	for i := 0; i < channel.processesNumber; i++ {
		fmt.Printf("Send: %d\n",i)
		message.receiverId = i
		res = res && channel.Send(message)
		if !res{
			fmt.Errorf("ERROR broadcast send to %d failed",i)
		}
	}
	return res
}

func (channel *Channel) Deliver(processId int, round int) *Message {
	if channel.messagesBuffer[processId].IsEmpty() {
		//fmt.Printf("entrato1");
		channel.receive(processId)
		//fmt.Printf("entrato3");
		//if(channel.count!=0){
		//fmt.Printf("Send: %d, Rec: %d Round: %d\n",channel.countSend,channel.count, round)}
	}
	var res *Message = channel.messagesBuffer[processId].Pop()
	if res==nil{

	}
	return res
}

