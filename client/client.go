package main

import (
	"context"
	"fmt"
	"kvtest/kv"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	addr := "localhost:10001"
	fmt.Printf("Connecting to RPC Server\n")

	rpcConn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect to RPC: %v", err)
	}

	rpcClient := kv.NewRPCServiceClient(rpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	voter := &kv.Voter{
		Address: "localhost:10002",
		Id:      "node2",
	}
	ret, err := rpcClient.AddNode(ctx, voter)
	if err != nil {
		log.Fatalf("Could not AddNode: %v %v\n", err, voter)
	}
	log.Printf("Response from RPC Server for AddNode : %d %s", ret.GetCode(), ret.GetMessage())

	voter = &kv.Voter{
		Address: "localhost:10003",
		Id:      "node3",
	}
	ret, err = rpcClient.AddNode(ctx, voter)
	if err != nil {
		log.Fatalf("Could not AddNode: %v %v\n", err, voter)
	}
	log.Printf("Response from RPC Server for AddNode : %d %s", ret.GetCode(), ret.GetMessage())

	ep := &kv.EmptyParams{}
	ret, err = rpcClient.Monitor(ctx, ep)
	if err != nil {
		log.Fatalf("Could not Monitor: %v", err)
	}
	log.Printf("Response from RPC Server for Monitor : %d %s", ret.GetCode(), ret.GetMessage())

	kvi := &kv.KV{
		Key: "key1",
		Val: "value1",
	}
	fmt.Printf("Storing in DB via rpc %v\n", kvi)

	ret, err = rpcClient.StoreKV(ctx, kvi)

	if err != nil {
		log.Fatalf("Could not Store person: %v", err)
	}

	log.Printf("Response from RPC Server for Store Data : %d %s", ret.GetCode(), ret.GetMessage())

	pkey := &kv.Key{
		Key: "key1",
	}

	fmt.Printf("Retrieving data from DB via rpc %v\n", pkey)

	r_kv, err := rpcClient.GetVal(ctx, pkey)

	if err != nil {
		log.Fatalf("Could not retrieve person: %v", err)
	}

	fmt.Printf("Retrieved from DB via rpc %s %s\n", r_kv.GetKey(), r_kv.GetVal())
}
