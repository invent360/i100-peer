package network

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const DEBUG = "\033[31m[*][DEBUG]\033[0m"

const protocolTopicName = "/bitcoin-simulation/1.0.6"
const contentRepeat = 8
const periodicIHAVE = 15

type TopicNetwork struct {
	Messages chan *Message
	Blocks   map[string]string
	Headers  []string

	ctx   context.Context
	host  host.Host
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription
	self  peer.ID
}

type Message struct {
	MsgType int //0 data, 1 ihave, 2 iwant
	Sender  peer.ID
	Blocks  map[string]string
	IHAVE   []string
	IWANT   []string
}

//Join the GossipSub network
func JoinNetwork(ctx context.Context, host host.Host, ps *pubsub.PubSub, self peer.ID) (*TopicNetwork, error) {

	//Join the topic
	topic, err := ps.Join(protocolTopicName)
	if err != nil {
		return nil, err
	}

	//Subscribe to the topic
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	net := &TopicNetwork{
		ctx:      ctx,
		host:     host,
		ps:       ps,
		topic:    topic,
		sub:      sub,
		self:     self,
		Messages: make(chan *Message, 1024),
		Blocks:   make(map[string]string),
		Headers:  make([]string, 0),
	}

	log.Println("- Network joined successfully, topic:", topic.String())

	go net.ReadService()
	return net, nil

}

//Loop service for GossipSub that handles IHAVE messages
func (net *TopicNetwork) ReadService() {
	for {

		//Get next message in the topic
		received, err := net.sub.Next(net.ctx)
		if err != nil {
			close(net.Messages)
			return
		}

		//If I'm the sender, ignore the message
		if received.ReceivedFrom == net.self {
			log.Println("- I am the sender, ignoring the packet")
			continue
		}

		//Unmarshal the message
		message := new(Message)
		err = json.Unmarshal(received.Data, message)
		if err != nil {
			continue
		}

		printMessage(*message)

		//Handle the IHAVE message (IWANT and DATA are send directly, so GossipSub should not see them)
		if message.MsgType == 1 {
			iwant := make([]string, 0)
			for _, owned := range message.IHAVE {
				if _, found := net.Blocks[owned]; !found {
					iwant = append(iwant, owned)
				}
				//Set maximum request to 16 blocks
				if len(iwant) == 16 {
					break
				}
			}
			//If I need some blocks, ask for them with a direct IWANT message
			if len(iwant) > 0 {
				msg := &Message{
					MsgType: 2,
					Sender:  net.self,
					Blocks:  map[string]string{},
					IHAVE:   []string{},
					IWANT:   iwant,
				}
				net.directSend(message.Sender, *msg)
			} else {
				log.Println("- No blocks needed")
			}
			//Forward the IHAVE message in the network to see if someone else needs the blocks listed here
			log.Println("- Forwarding IHAVE message on the network")
			net.Messages <- message
		} else {
			fmt.Println(DEBUG, "This was not expected, message type", message.MsgType)
		}

	}
}

// Publish the message on the GossipSub network (used only for IHAVE messages)
func (net *TopicNetwork) Publish(ihave []string) error {

	//Create and marshal the message
	message := &Message{
		MsgType: 1,
		Sender:  net.self,
		Blocks:  map[string]string{},
		IHAVE:   ihave,
		IWANT:   []string{},
	}
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	//Publish the message on the GossipSub network
	err = net.topic.Publish(net.ctx, msg)
	if err != nil {
		return err
	}
	log.Printf("- Message IHAVE published (%d chars): %s\n", len(msg), string(msg))
	return nil

}

//Stream handler for DATA and IWANT messages
func (net *TopicNetwork) handleStream(s network.Stream) {

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	//Receive and respond message loop
	go func(rw *bufio.ReadWriter) {

		//Create the decoder for the stream
		decoder := json.NewDecoder(rw)
		var message Message

		for {

			//Get the message and decode it
			err := decoder.Decode(&message)

			if err != nil {
				log.Println("- Stream closed:", err)
				s.Reset()
				return
			}

			printMessage(message)

			//Handle the different messages
			switch message.MsgType {
			case 0: //Check the blocks received and store the ones not already stored

				for header, content := range message.Blocks {
					if _, found := net.Blocks[header]; !found {
						net.Blocks[header] = content
						net.Headers = append(net.Headers, header)
						log.Printf("- Stored block #%d %s: (x %d) %s\n", len(net.Headers), header, contentRepeat, content[:16])
					} else {
						log.Println("- I already have the block", header) //This should not happen
					}
				}
				log.Printf("- Now I have #%d blocks: %s\n", len(net.Headers), net.Headers)

			case 2: //Send the requested blocks with a DATA message

				//Create a map of requested block to send to the peer
				toSend := make(map[string]string)
				for _, wanted := range message.IWANT {
					if _, found := net.Blocks[wanted]; found {
						toSend[wanted] = net.Blocks[wanted]
					}
				}

				//Send DATA message directly to the peer that requested those blocks
				msg := &Message{
					MsgType: 0,
					Sender:  net.self,
					Blocks:  toSend,
					IHAVE:   []string{},
					IWANT:   []string{},
				}
				net.directSend(message.Sender, *msg)

			default:
				fmt.Println(DEBUG, "This was not expected, message type", message.MsgType)
			}
		}
	}(rw)

}

//Send directly from a peer to another
func (net *TopicNetwork) directSend(receiver peer.ID, msg Message) {

	//Open the stream to the receiver peer
	stream, err := net.host.NewStream(net.ctx, receiver, protocolTopicName)
	if err != nil {
		log.Printf("- Error opening stream to %s: %s\n", receiver, err)
		return
	}

	//Marshal the message to send it
	message, err := json.Marshal(msg)
	if err != nil {
		log.Println("- Error marshalling message to send:", err)
		return
	}

	//Write the message on the stream
	nWritten, err := stream.Write(message)
	if err != nil {
		log.Println("- Error sending message on stream:", err)
		return
	}

	if msg.MsgType == 0 {
		log.Printf("- Message DATA sent to %s (%d chars) with blocks %s", receiver, nWritten, getKeys(msg.Blocks))
	} else if msg.MsgType == 2 {
		log.Printf("- Message IWANT sent to %s: %s", receiver, string(message))
	}

}

//Print received message
func printMessage(msg Message) {
	switch msg.MsgType {
	case 0:
		log.Printf("- Message DATA received from %s, blocks %s", msg.Sender, getKeys(msg.Blocks))
	case 1:
		log.Printf("- Message IHAVE received from %s, it has %s", msg.Sender, msg.IHAVE)
	case 2:
		log.Printf("- Message IWANT received from %s, it wants %s", msg.Sender, msg.IWANT)
	}
}

//Get the keys of a map
func getKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

//PeriodicBroadcast send IHAVE messages on the network for newly entered peers
func PeriodicBroadcast(net *TopicNetwork) {
	for {
		time.Sleep(time.Second * periodicIHAVE)
		peers := net.ps.ListPeers(protocolTopicName)
		log.Printf("- Found %d other peers in the network: %s\n", len(peers), peers)

		messages := []string{"test 1", "test 2", "test 3"}

		if err := net.Publish(messages); err != nil {
			log.Println("- Error publishing IHAVE message on the network:", err)
		}
	}

}
