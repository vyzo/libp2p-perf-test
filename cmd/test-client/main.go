package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"

	quic "github.com/libp2p/go-libp2p-quic-transport"
	tcp "github.com/libp2p/go-tcp-transport"
	ma "github.com/multiformats/go-multiaddr"
)

const TestProtocol = protocol.ID("/libp2p/test/data")

func main() {
	if len(os.Args) != 2 {
		log.Fatal("expected one argument; the peer multiaddr")
	}

	a, err := ma.NewMultiaddr(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	pi, err := peer.AddrInfoFromP2pAddr(a)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	host, err := libp2p.New(ctx,
		libp2p.NoListenAddrs,
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(quic.NewTransport),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connecting to %s", pi.ID.Pretty())

	cctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err = host.Connect(cctx, *pi)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected; requesting data...")

	s, err := host.NewStream(cctx, pi.ID, TestProtocol)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	file, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.Printf("Transfering data...")

	start := time.Now()
	n, err := io.Copy(file, s)
	if err != nil {
		log.Printf("Error receiving data: %s", err)
	}
	end := time.Now()

	log.Printf("Received %d bytes in %s", n, end.Sub(start))
}
