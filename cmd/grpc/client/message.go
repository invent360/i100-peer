package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/i101-p2p/cmd/grpc/proto/message"
	msg "github.com/i101-p2p/cmd/grpc/proto/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Message struct {
	client msg.MessageClient
}

func NewClient(conn grpc.ClientConnInterface) Message {
	return Message{
		client: msg.NewMessageClient(conn),
	}
}

func (m *Message) Subscribe(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := m.client.SubscribeToPeer(context.Background())

	if err != nil {
		return fmt.Errorf("create stream: %w", err)
	}

	req := message.MessageRequest{
		Title:   "Golang",
		Payload: "Some payload",
	}

	if err := stream.Send(&req); err != nil {
		return fmt.Errorf("error while calling RPC server: %w", err)
	}

	for {
		msgResponse, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("close and receive: %w", err)
		}

		fmt.Printf("Respons from GRPC server %v\n\n", msgResponse.GetFeedback())
	}
	return nil
}

func RunClient(addr string) {
	log.Println("...starting client @", fmt.Sprintf("%s:%s", strings.Split(addr, "/")[2], strings.Split(addr, "/")[4]))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", strings.Split(addr, "/")[2], strings.Split(addr, "/")[4]), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Fatalf("could not defer connection: %v", err)
		}
	}(conn)

	client := NewClient(conn)
	if err := client.Subscribe(context.Background()); err != nil {
		log.Fatalln(err)
	}
}
