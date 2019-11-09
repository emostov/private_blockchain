package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Bc is the blochain instance

// other port to use 8080 switch self id and peer id 8080 6689

var peerID = []string{"http://localhost:6689"}

// PEERLIST right now is just hard coded
var PEERLIST = PeerList{selfID: SELFID, peerIDs: peerID, length: Bc.Length}

// SELFID ...
var SELFID = []string{"http://localhost:", "8080"}

// SELFADDRESS ...
var SELFADDRESS = "http://localhost:" + SELFID[1]

func main() {
	miner1()
}
func miner1() {

	testSetupBlockInsert()

	//router := NewRouter()
	// go Bc.StartTryingNonces()
	Bc.StartTryingNonces()

	fmt.Println(Bc.Show())
	testAsk()
	fmt.Println("I am at port ", SELFID[1])
	router := NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":"+SELFID[1], router))
	}
}

func testSetupBlockInsert() {
	var blocks = makeTenBlocks()
	Bc.Insert(blocks[0])
	Bc.Insert(blocks[1])
	Bc.Insert(blocks[2])
}

func test4() {
	bc := NewBlockChain()

	blocks := makeTenBlocks()
	for _, b := range blocks {
		bc.Insert(b)
	}
	bc.GetLatestBlock()
	// fmt.Println(generateStartNonce(100000000000))
	// fmt.Println(bc.TestTryNonces())

}

func testAsk() {
	askForParent("97611bc0a6e098f600d4c776252ffc16173058b4d5e2ae4a7d336fe18eb7f11326b2ca4e40be4c5572800ca76c6cdf4b65931b297f098f73a256f497c8907736", "1")
}

// func test2() {
// 	// testing insertion of a block into the block chain
// 	bc := NewBlockChain()
// 	bc.Insert(makeGenesisBlock())
// 	JSONBc := bc.EncodeToJSON()
// 	printStringSlice(JSONBc)
// }

func test3() {
	// testing creating a blockchain, and block chain encoding and decoding
	bc := NewBlockChain()
	blocks := makeTenBlocks()
	bc.Show()
	for _, b := range blocks {
		bc.Insert(b)
	}
	fmt.Println(bc.Show())
}
