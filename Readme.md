# Simple distributed k-v store using hashicorp raft for consensus

To build
1. Check out git code
2. make clean
3. make

To run kv/raft nodes:
1. ./kvserver --address=localhost:10001 --id=node1 --bs=true
2. ./kvserver --address=localhost:10002 --id=node2 
3. ./kvserver --address=localhost:10003 --id=node3 

To Use web/server:
1. curl localhost:9000/AddVoter/localhost:10002/node2
2. curl localhost:9000/AddVoter/localhost:10003/node3
3. curl localhost:9000/Get/key1 
4. curl localhost:9000/Put/key1/value1
5. curl localhost:9000/Get/key1 

To use sample rpc client

```
Usage of ./kvclient:
  -address string
    	TCP host+port for this node (default "localhost:10002")
  -command string
    	monitor, addnode, set(key,val), get(key) (default "help")
  -id string
    	Node ID (default "node2")
  -key string
    	Key to test (default "somekey1")
  -val string
    	Val to set and test for key (default "someval1")
```

% ./kvclient --command=addnode --address=localhost:10002 --id=node2
```
Connecting to RPC Server addnode
2022/02/22 14:50:46 Response from RPC Server for AddNode : 1 Success
```
% ./kvclient --command=addnode --address=localhost:10003 --id=node3
```
Connecting to RPC Server addnode
2022/02/22 14:50:52 Response from RPC Server for AddNode : 1 Success
```
% ./kvclient --command=set --key=key1 --val=value1                 
```
Connecting to RPC Server set
Storing in DB via rpc key:"key1" val:"value1" 
2022/02/22 14:51:08 Response from RPC Server for Store Data : 1 Success
```
 % ./kvclient --command=get --key=key1             
```
Connecting to RPC Server get
Retrieving data from DB via rpc key:"key1" 
Retrieved from DB via rpc key1 value1
```

# Introduction

1. main launches the raft/KV and http server
2. http server (get/put) is used to wrap the raft/kv commands
3. main also launches rpc server which listens on the same port as passed in the address

Each node is running an instance of raft under rpc. The raft part comes w/ hashicorp raft lib - the rpce is a wrapper on top. There's another http/web wrapper as well.

Each rpc service implements basic RAFT functions to 
1. AddNode (Voter)
2. Get
3. Put
4. Monitor

The FSM in this case is a KV - map (string,string) which is hidden inside kv and exposed via APIs above. 

Clients (like rpc client or http curl) interact w/ the rpc service and thus manipulate and get State via APIs.


Note: A lot of hardcoding in the code - needs to be weeded out, Implement snapshot restore etc.
