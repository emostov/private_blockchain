Bare Bones Private Block Chain

Zeke Mostov
October, 14, 2019

To create a block declare and instance and then use the block initialize method.
To create a blockchain use the NewBlockChain() method create an instance that is
initialized with an empty map.

Blocks supports Initialize, DecodeFromJson and EncodeFromJson.

BlockChain supports Get, Insert, EncodeToJson, and DecodeFromJson. BlockChain
length starts at 1

make_genesis_block() is useful for making an arbitrary genesis block for a
Blockchain.
