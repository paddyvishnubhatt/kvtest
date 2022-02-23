package main

import (
	"flag"
	"kvtest/kv"
	"kvtest/server"
	"strings"
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
		s := strings.Split(*serverAddr, ":")
		srv := server.HTTPServer{}
		go srv.RunServer(mykv, strings.Trim(s[0], " "), *port)
	}

	mykv.InitRaft(*serverAddr, *nodeId, *bs)
}
