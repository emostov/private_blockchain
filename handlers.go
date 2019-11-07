package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

// TestTryNonces ...
func (bc *BlockChain) TestTryNonces() string {
	parentBlock := bc.GetLatestBlock()[0]
	fmt.Println(parentBlock)
	parentHash := parentBlock.Header.Hash
	var b Block
	b.Initialize(bc.Length+1, parentHash, "test block value")
	blockValue := b.Value
	target := "000000" // six 0's
	nonce := generateStartNonce(1)
	run := true
	puzzleAnswer := ""
	for run {
		run = false
		puzzleAnswer = makePuzzleAnswer(nonce, parentHash, blockValue)
		if checkPuzzleAnswerValid(target, puzzleAnswer) == false {
			nonce = generateNonce(nonce)
			run = true
		}
	}
	b.Header.Nonce = nonce
	bc.Insert(b)
	return nonce
}

// StartTryingNonces ...
func (bc *BlockChain) StartTryingNonces() {
	tryNonce := true
	for tryNonce {
		//turn tryNonce to false if recieve another block
		// b := generateBlock
		parentBlock := bc.GetLatestBlock()[0]
		parentHash := parentBlock.Header.Hash
		var b Block
		b.Initialize(bc.Length+1, parentHash, "test block value")
		blockValue := b.Value
		target := "00000" // ten 0's
		nonce := generateStartNonce(1)

		run := true
		puzzleAnswer := ""
		for run == true {
			run = false
			puzzleAnswer = makePuzzleAnswer(nonce, parentHash, blockValue)
			checkPuzzleAnswerValid(target, puzzleAnswer)
			run = false
			// if !(checkPuzzleAnswerValid(target, puzzleAnswer)) {
			// 	nonce = generateNonce(nonce)
			// 	fmt.Println("nonce is", nonce)
			// 	run = true
			// }
		}
		// Broadcast Node with new nonce with heartbeat data
		// Or add recieved valid node to block chain

	}

}

//Helper function for Start Trying Nonces

// takes hexadecimal string, converts to int, adds 1 and converts back to int
func generateNonce(previous string) string {
	previousInt, err := strconv.Atoi(previous)
	if err != nil {
		log.Fatal("error in generateNonce", err)
	}

	previousPlusOne := previousInt + 10

	newNonce := strconv.Itoa(previousPlusOne)
	return newNonce
}

func generateStartNonce(seed int) string {
	str := strconv.Itoa(seed)
	return str
}

func makePuzzleAnswer(nonce string, parentBlockHash string, blockValue string) string {
	puzzleString := parentBlockHash + nonce + blockValue
	puzzleAnswer := makeSha256Digest(puzzleString)
	return puzzleAnswer
}

func checkPuzzleAnswerValid(target string, puzzleAnswer string) bool {
	// fmt.Println("target: ", target, "puzzleAnswer ", puzzleAnswer)
	return target == puzzleAnswer[:len(target)]
}

// Handlers Below

// func handlers() {
// 	router := mux.NewRouter().StrictSlash(true)
// }

// Bc is the blochain instance
var Bc = NewBlockChain()
var blocks = makeTenBlocks()

func testSetup() {
	Bc.Insert(blocks[0])
}

// Upload ...
func Upload(w http.ResponseWriter, r *http.Request) {

	// fmt.Println(Bc.EncodeToJSON())
	// if err != nil {
	// 	w.WriteHeader(http.StatusServiceUnavailable)
	// }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(EncodeBlockchainToJSON(Bc)))

}

// AskForBlock ...
func AskForBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	h := vars["height"]
	hash := vars["hash"]
	height, err := strconv.Atoi(h)
	if err == nil {
		block := Bc.GetBlock(int32(height), hash)
		if block == nil {
			w.WriteHeader(http.StatusNotFound)
			// Ask for parent block, insert current block into tree
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(block.EncodeToJSON())
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

	w.Write([]byte("height :" + string(height) + ", hash: " + hash))
}

// ShowHandler ...
func ShowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Bc.Show()))
}

// func HeartBeatRecieve() (w http.ResponseWriter, r *http.Request) {

// }

// HeartBeatReceive()
//  When a node receives a new block in HeartBeat, the node will first check if
//  the parent block of this new block exists in its own blockchain (the previous
// 	block is the block whose hash is the parentHash of the next block).
// If the previous block doesn't exist, the node will ask any peer in PeerList at
// "/block/{height}/{hash}" to download that block.
// After making sure the previous block exists, insert the block from HeartBeatData
//  to the current BlockChain.
//  Alter this function so that when it receives a HeartBeatData with a new block,
//  it verifies the nonce as described above.

/////////////

// type stringArray struct {
// 	strArray []string `json:"strArray"`
// }

// func NewStringArray() stringArray {
// 	return stringArray{
// 		make([]string, 0),
// 	}
// }
