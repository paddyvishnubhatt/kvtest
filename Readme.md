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

Note: A lot of hardcoding in the code - needs to be weeded out, Implement snapshot restore etc. Wrap this around grpc instead of http server etc
