# Bare Bones Block Chain

# Zeke Mostov
November 25, 2019

#INSTRUCTIONS TO RUN:
Before running any nodes, registration node set up is not required, but highly 
reccomended to gurantee that the new node will 'meet' other nodes in the network
and to ensure it has an up to date chain to start mining on.
Non-registration nodes are hardcoded to download there initial chain and peerlist
from a registration node. To start the registration simply start at command line
and pass no arguments:

"go run Block.go block_generation.go BlockChain.go handlers.go HeartBeat.go logger.go Main.go routes.go broadcast_network.go SyncBlockChain.go Register.go"

To start nodes, run at the command line and pass a valid port number that is not
in use. For, example to run at port 8001:
"go run Block.go block_generation.go BlockChain.go handlers.go HeartBeat.go logger.go Main.go routes.go broadcast_network.go SyncBlockChain.go Register.go 8001"

# ABOUT:
- Currently only a simple POW protocol is implemented, meaning the target value 
is not changed while mining. 
- Nodes start with no hard coded peers other then the ip address for the
registration node. The registration node starts with no hard coded peers. 
- Every time a node registers, the registration node responds with the nodes 
ID and a peermap. The registration node then updates the peermap with that nodes ID.
- There is a very simple gossip protocol. When a node sends a heart beat to a 
peer, the peer updates there peerlist with the senders id if they do not already
have. Next step will be to add functionality so the recieving node updates its 
peerlist with the senders entire peermap.
- Each node has a difficulty field, which is the sume of the difficulty of all
of its parents and the node itself. ShowCanonical() takes advantage of this field
to effeciently O(1) identify the canoninical chain (the chain with the most 
difficulty)
- The nonce is seeded at a random int between 0 and 300 and then increments. It
starts at a random int so it is less deterministic as to who will mine a block 
first
- The mining algorithmn is naive in that it just picks the first block at the 
current length of the chain. A next step would be to implement a function to find
the canonical chain and then mine on that
- askForParent() has a timeout of 5 seconds, so if it does not get a response in
5 seconds it will discontinue waiting.

# NEXT STEPS:
- Add full gossip protocol
- Make sure nodes mining on highest difficulty chain
- Dynamically adjust mining target
