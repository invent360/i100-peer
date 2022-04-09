package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	cfg "github.com/i101-p2p/cmd/config"
	"github.com/i101-p2p/network"
	p2p "github.com/i101-p2p/node"
	"github.com/ipfs/go-log"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var logger = log.Logger("rendezvous")

func main() {

	ctx := context.Background()

	config, err := cfg.ParseFlags()
	if err != nil {
		panic(err)
	}

	node, err := p2p.NewNode(ctx, config)

	if err != nil {
		panic(err)
	}

	//go server.RunServer(node.Host.Addrs()[0].String())

	go node.AdvertiseAndFindPeers(ctx, config)

	//------------
	//Create a new PubSub service using GossipSub routing and join the topic
	ps, err := pubsub.NewGossipSub(ctx, node.Host)
	if err != nil {
		panic(err)
	}
	topic, err := network.JoinNetwork(ctx, node.Host, ps, node.Host.ID())
	if err != nil {
		panic(err)
	}

	go network.PeriodicBroadcast(topic)
	//------------
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	logger.Debug("Received signal, shutting down...")

	// shut the node down
	if err := node.Host.Close(); err != nil {
		panic(err)
	}
}
