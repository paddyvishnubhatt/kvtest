package main

import (
	"flag"
	"kvtest/kv"
	"kvtest/server"
)

var (
	serverAddr = flag.String("address", "localhost:10001", "TCP host+port for this node")
	nodeId     = flag.String("id", "node1", "Node ID")
	bs         = flag.Bool("bs", false, "Bootstrap")
	port       = flag.Int("port", 9000, "Port running the http server on bs node")
)

func main() {
	flag.Parse()
	mykv := &kv.KeyVal{}

	if *bs {
		go server.RunServer(mykv, *serverAddr, *port)
	}

	mykv.InitRaft(*serverAddr, *nodeId, *bs)
}
