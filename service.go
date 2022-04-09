package main

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/i101-p2p/protocol"
	"github.com/ipfs/go-log/v2"
)

var logger = log.Logger("rendezvous")

func RunClient() {
	book := protocol.Book{
		Name: "Kamasutra, the ultimate sex guide",
		Isbn: 12334,
	}

	data, err := proto.Marshal(&book)

	fmt.Println("Book in bytes:", data)
	if err != nil {
		logger.Error("error marshaling message:", err)
	}

	sendData(data)
}

func RunServer() {
	listener, err := net.Listen("tcp", "127.0.0.1:8085")
	if err != nil {
		logger.Fatal(err)
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			logger.Fatal(err)
		}

		go func(c net.Conn) {

			defer c.Close()
			data, err := ioutil.ReadAll(connection)

			if err != nil {
				logger.Fatal(err.Error())
			}

			book := protocol.Book{}

			err = proto.Unmarshal(data, &book)
			if err != nil {
				logger.Fatal(err.Error())
			}

			fmt.Println(&book)
		}(connection)
	}
}

func sendData(data []byte) {
	connection, err := net.Dial("tcp", "127.0.0.1:8085")

	if err != nil {
		logger.Fatal(err)
	}
	defer connection.Close()

	write, err := connection.Write(data)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info(write)
}
