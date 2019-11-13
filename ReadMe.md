Bare Bones Private Block Chain

Zeke Mostov
November, 13, 2019

INSTRUCTIONS TO RUN:
Before running any nodes, registration server set up is required. This is because
currently the node will log a fatal error if it sends request to server and it is
not up. To start the registration simply start at command line and pass no arguments:

"go run Block.go block_generation.go BlockChain.go handlers.go HeartBeat.go logger.go Main.go routes.go broadcast_network.go SyncBlockChain.go Register.go"

To start nodes, run at the command line and pass a valid port number that is not
in use. For, example to run at port 8001:
"go run Block.go block_generation.go BlockChain.go handlers.go HeartBeat.go logger.go Main.go routes.go broadcast_network.go SyncBlockChain.go Register.go 8001"

ABOUT:
- Currently only a simple POW protocol is implemented, meaning the target value is
not changed while mining. 
- Nodes start with no hard coded peers. The registration
server also starts with no hard coded peers. 
- Every time a node registers, the registration server responds with the nodes 
ID and a peermap. The server then updates the peermap with that nodes ID.
- When a node sends a heart beat to a peer, the peer updates there peerlist with
the senders id if they do not already have. Next step will be to add functionality
so the recieving node updates its peerlist with the senders entire peermap.
