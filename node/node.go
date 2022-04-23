package node

import (
	"context"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"strings"
	"sync"

	"github.com/i101-p2p/cmd/config"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libp2pquic "github.com/libp2p/go-libp2p-quic-transport"
	"github.com/libp2p/go-tcp-transport"
	"github.com/multiformats/go-multiaddr"
)

type Node struct {
	Host   host.Host
	KadDHT *dht.IpfsDHT
}

func NewNode(ctx context.Context, config config.Config) (*Node, error) {
	var r io.Reader

	r = mrand.New(mrand.NewSource(config.Seed))

	privateKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", config.Port))
	p2pHost, err := libp2p.New(
		libp2p.ListenAddrs(addr),
		libp2p.Identity(privateKey),
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.DefaultMuxers,
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.FallbackDefaults,
	)

	if err != nil {
		fmt.Printf("errpr creating p2p host %s", err)
	}

	log.Println("###########")
	log.Println("I am:", p2pHost.ID())
	for _, addr := range p2pHost.Addrs() {
		fmt.Printf("%s: %s/p2p/%s", "I am @:", addr, p2pHost.ID().Pretty())
		fmt.Println()
	}
	log.Println("###########")
	//#########################################
	kadDHT, err := dht.New(ctx, p2pHost)
	if err != nil {
		panic(err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	log.Println("bootstrapping the DHT")
	if err = kadDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}
	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range config.BootstrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p2pHost.Connect(ctx, *peerInfo); err != nil {
				log.Println(err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerInfo)
			}
		}()
	}
	wg.Wait()
	return &Node{KadDHT: kadDHT, Host: p2pHost}, nil
}

func (node Node) AdvertiseAndFindPeers(ctx context.Context, cfg config.Config) {
	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	log.Println("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(node.KadDHT)
	discovery.Advertise(ctx, routingDiscovery, cfg.Rendezvous)
	log.Println("Successfully announced!")

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	for {
		peers, err := routingDiscovery.FindPeers(ctx, cfg.Rendezvous)
		if err != nil {
			log.Printf("error finding peers: ", err)
		}
		for p := range peers {
			if p.ID == node.Host.ID() {
				continue
			}
			if len(p.Addrs) > 0 {
				log.Println("found peers: ")
				for _, addr := range p.Addrs {
					fmt.Printf("%s/p2p/%s", addr, p.ID.Pretty())
					fmt.Println()
					//fmt.Println("--------------NewMultiaddr-----------------")
					//fmt.Println(ma.NewMultiaddr(addr.String()))
					//fmt.Println("--------------AddrInfoFromP2pAddr-----------------")
					//fmt.Println(peer.AddrInfoFromP2pAddr(addr))
				}
				status := node.Host.Network().Connectedness(p.ID)
				if status == network.CanConnect || status == network.NotConnected {
					fmt.Println(p.Addrs[1].String())
					fmt.Println(fmt.Sprintf("%s:%s", strings.Split(p.Addrs[0].String(), "/")[2], strings.Split(p.Addrs[0].String(), "/")[4]))
					if err := node.Host.Connect(ctx, p); err != nil {
						log.Printf("Connection failed:", err)
					} else {
						log.Println("connected to peer: ", p.ID)
						//go client.RunClient(p.Addrs[0].String())
					}
				}
			}

		}
	}
}
