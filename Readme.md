# Simple k-v store using hashicorp raft

1. Check out git code
2. make clean
3. make
4. ./kvserver --address=localhost:10001 --id=node1 --bs=true
5. ./kvserver --address=localhost:10002 --id=node2 
6. ./kvserver --address=localhost:10003 --id=node3 
7. curl localhost:9000/AddVoter/localhost:10002/node2
8. curl localhost:9000/AddVoter/localhost:10003/node3
9. curl localhost:9000/Get/key1 
10. curl localhost:9000/Put/key1/value1
11. curl localhost:9000/Get/key1 

1. main launches the raft/KV and http server
2. http server (get/put) is used to wrap the raft/kv commands
3. main also launches rpc server which listens on the same port as passed in the address
4. sample rpc client output

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

% ./kvclient --command=addnode --address=localhost:10002 --id=node2
Connecting to RPC Server addnode
2022/02/22 14:50:46 Response from RPC Server for AddNode : 1 Success

% ./kvclient --command=addnode --address=localhost:10003 --id=node3
Connecting to RPC Server addnode
2022/02/22 14:50:52 Response from RPC Server for AddNode : 1 Success

% ./kvclient --command=set --key=key1 --val=value1                 
Connecting to RPC Server set
Storing in DB via rpc key:"key1" val:"value1" 
2022/02/22 14:51:08 Response from RPC Server for Store Data : 1 Success

 % ./kvclient --command=get --key=key1             
Connecting to RPC Server get
Retrieving data from DB via rpc key:"key1" 
Retrieved from DB via rpc key1 value1


Note: A lot of hardcoding in the code - needs to be weeded out, Implement snapshot restore etc.
