package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"

	quic "github.com/libp2p/go-libp2p-quic-transport"
	tcp "github.com/libp2p/go-tcp-transport"
	ma "github.com/multiformats/go-multiaddr"
)

const TestProtocol = protocol.ID("/libp2p/test/data")

func main() {
	streams := flag.Int("streams", 1, "number of parallel download streams")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] peer", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	a, err := ma.NewMultiaddr(flag.Args()[0])
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

	if *streams == 1 {
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
	} else {
		var wg sync.WaitGroup
		var count int32

		dataStreams := make([]network.Stream, 0, *streams)
		for i := 0; i < *streams; i++ {
			s, err := host.NewStream(cctx, pi.ID, TestProtocol)
			if err != nil {
				log.Fatal(err)
			}
			defer s.Close()
			dataStreams = append(dataStreams, s)
		}

		log.Printf("Transferring data in %d parallel streams", *streams)

		start := time.Now()
		for i := 0; i < *streams; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				file, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				n, err := io.Copy(file, dataStreams[i])
				if err != nil {
					log.Printf("Error receiving data: %s", err)
				}
				atomic.AddInt32(&count, int32(n))
			}(i)
		}

		wg.Wait()
		end := time.Now()
		log.Printf("Received %d bytes in %s", count, end.Sub(start))

	}
}
