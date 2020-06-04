# A couple of simple programs for testing libp2p transport performance

## Installation
```
$ git clone https://github.com/vyzo/libp2p-perf-test.git
$ cd libp2p-perf-test
$ go install ./...
```

## Usage

Programs:
- `test-server` is the server serving a file, over TCP and QUIC.
- `test-client` is the client connecting to the server and downloading the served file.

The server expects a file named `data` to serve (can override with `-data` option).
So let's generate a random 1GiB file and start the server:

```
$ dd bs=1024 count=1048576 if=/dev/urandom of=data
$ test-server
I am /ip4/127.0.0.1/udp/4001/quic/p2p/QmSfCAZ7kq8m3yuWip2gC6cE5YrpxtwDsC3wY4NMdj5a7Z
I am /ip4/192.168.2.4/udp/4001/quic/p2p/QmSfCAZ7kq8m3yuWip2gC6cE5YrpxtwDsC3wY4NMdj5a7Z
I am /ip4/127.0.0.1/tcp/4001/p2p/QmSfCAZ7kq8m3yuWip2gC6cE5YrpxtwDsC3wY4NMdj5a7Z
I am /ip4/192.168.2.4/tcp/4001/p2p/QmSfCAZ7kq8m3yuWip2gC6cE5YrpxtwDsC3wY4NMdj5a7Z

```

And now we can connect the client and download the data file, copying it to `/dev/null`.

With TCP:
```
$ test-client /ip4/192.168.2.4/tcp/4001/p2p/QmSfCAZ7kq8m3yuWip2gC6cE5YrpxtwDsC3wY4NMdj5a7Z
2020/06/04 20:59:09 Connecting to QmSfCAZ7kq8m3yuWip2gC6cE5YrpxtwDsC3wY4NMdj5a7Z
2020/06/04 20:59:09 Connected; requesting data...
2020/06/04 20:59:09 Transfering data...
2020/06/04 20:59:13 Received 1073741824 bytes in 4.273251937s
```

With QUIC:
```
$ test-client /ip4/192.168.2.4/udp/4001/quic/p2p/QmSfCAZ7kq8m3yuWip2gC6cE5YrpxtwDsC3wY4NMdj5a7Z
2020/06/04 20:59:28 Connecting to QmSfCAZ7kq8m3yuWip2gC6cE5YrpxtwDsC3wY4NMdj5a7Z
2020/06/04 20:59:28 Connected; requesting data...
2020/06/04 20:59:28 Transfering data...
2020/06/04 20:59:38 Received 1073741824 bytes in 9.798873952s
```
