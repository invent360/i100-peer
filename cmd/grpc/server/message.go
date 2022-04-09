package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	msg "github.com/i101-p2p/cmd/grpc/proto/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Message struct {
	msg.MessageServer
}

func NewMessage() *Message {
	return &Message{}

}

func (m *Message) MessagePeer(ctx context.Context, request *msg.MessageRequest) (*msg.MessageResponse, error) {
	log.Println("Handling MessagePeer")
	return &msg.MessageResponse{Feedback: "Ok"}, nil
}

func (m *Message) SubscribeToPeer(msgServer msg.Message_SubscribeToPeerServer) error {
	log.Println("Handling SubscribeToPeer")
	go func() {
		for {
			req, err := msgServer.Recv()
			if err == io.EOF {
				log.Println("Client has closed connection")
				break
			}

			if err != nil {
				log.Fatal("Unable to read from client", err)
			}

			log.Println("Handle client request", req)
		}
	}()

	for {
		if err := msgServer.Send(&msg.MessageResponse{Feedback: "haha"}); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
}

func RunServer(addr string) {
	grpcServer := grpc.NewServer()
	messageServer := NewMessage()
	msg.RegisterMessageServer(grpcServer, messageServer)

	reflection.Register(grpcServer)

	log.Println("...starting server @", fmt.Sprintf("%s:%s", strings.Split(addr, "/")[2], strings.Split(addr, "/")[4]))

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", strings.Split(addr, "/")[2], strings.Split(addr, "/")[4]))
	if err != nil {
		log.Fatal("Unable to listen", err)
	}
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Unable to serve", err)
	}
}
