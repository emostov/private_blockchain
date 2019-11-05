	// StartTryingNonces(): This function starts a new thread that tries different
	// nonces to generate new blocks. Nonce is a string of 16 hexes such as
	// "1f7b169c846f218a". Initialize the rand when you start a new node with
	// something unique about each node, such as the current time or the port
	//number.

	// 	(1) Start a while loop.
	// (2) Get the latest block or one of the latest blocks to use as a parent block.
	// (3) Create an MPT.
	// (4) Randomly generate the first nonce, verify it with simple PoW algorithm to
	// see if SHA3(parentHash + nonce + mptRootHash) starts with 10 0's (or the number
	// 	 you modified into). Since we use one laptop to try different nonces, six
	// 	 to seven 0's could be enough. If the nonce failed the verification,
	// 	 increment it by 1 and try the next nonce.
	// (6) If a nonce is found and the next block is generated, forward that block to
	// all peers with an HeartBeatData;
	// (7) If someone else found a nonce first, and you received the new block through
	// your function ReceiveHeartBeat(), stop trying nonce on the current block,
	// continue to the while loop by jumping to the step(2).