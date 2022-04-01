package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ipfs/go-log/v2"
)

var logger = log.Logger("rendezvous")

func main() {

	log.SetAllLoggers(log.LevelWarn)
	if err := log.SetLogLevel("rendezvous", "info"); err != nil {
		return
	}

	ctx := context.Background()

	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	node, err := NewNode(ctx, config)

	if err != nil {
		panic(err)
	}

	go node.AdvertiseAndFindPeers(ctx, config)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	logger.Debug("Received signal, shutting down...")

	// shut the node down
	if err := node.Host.Close(); err != nil {
		panic(err)
	}
}
