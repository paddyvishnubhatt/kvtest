package main

import (
	"context"
	"flag"
	"fmt"
	"kvtest/kv"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
)

var (
	command    = flag.String("command", "help", "monitor, addnode, set(key,val), get(key)")
	serverAddr = flag.String("address", "localhost:10002", "TCP host+port for this node")
	nodeId     = flag.String("id", "node2", "Node ID")
	skey       = flag.String("key", "somekey1", "Key to test")
	sval       = flag.String("val", "someval1", "Val to set and test for key")
	port       = flag.String("port", "10001", "Port running rpcserver on")
)

func main() {
	flag.Parse()
	addr := "localhost:" + *port
	fmt.Printf("Connecting to RPC Server %v\n", *command)

	rpcConn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect to RPC: %v", err)
	}

	kvStore := kv.NewRPCServiceClient(rpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if strings.EqualFold(*command, "monitor") {
		ep := &kv.EmptyParams{}
		ret, err := kvStore.Monitor(ctx, ep)
		if err != nil {
			log.Fatalf("Could not Monitor: %v", err)
		}
		log.Printf("Response from RPC Server for Monitor : %d %s", ret.GetCode(), ret.GetMessage())
	} else if strings.EqualFold(*command, "help") {
		fmt.Println("Commands are help, addnode, set, get")
	} else if strings.EqualFold(*command, "addnode") {
		voter := &kv.Voter{
			Address: *serverAddr,
			Id:      *nodeId,
		}
		ret, err := kvStore.AddNode(ctx, voter)
		if err != nil {
			log.Fatalf("Could not AddNode: %v %v\n", err, voter)
		}
		log.Printf("Response from RPC Server for AddNode : %d %s", ret.GetCode(), ret.GetMessage())

	} else if strings.EqualFold(*command, "set") {
		kvi := &kv.KV{
			Key: *skey,
			Val: *sval,
		}
		fmt.Printf("Storing in DB via rpc %v\n", kvi)

		ret, err := kvStore.StoreKV(ctx, kvi)

		if err != nil {
			log.Fatalf("Could not Store KV: %v", err)
		}
		log.Printf("Response from RPC Server for Store Data : %d %s", ret.GetCode(), ret.GetMessage())
	} else if strings.EqualFold(*command, "get") {
		pkey := &kv.Key{
			Key: *skey,
		}

		fmt.Printf("Retrieving data from DB via rpc %v\n", pkey)

		r_kv, err := kvStore.GetVal(ctx, pkey)

		if err != nil {
			log.Fatalf("Could not retrieve person: %v", err)
		}

		fmt.Printf("Retrieved from DB via rpc %s %s\n", r_kv.GetKey(), r_kv.GetVal())
	} else {
		fmt.Println("Error - pls use help")
	}
}
