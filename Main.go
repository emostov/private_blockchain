package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Bc is the blochain instance

// other port to use 8080 switch self id and peer id 8080

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
	// JSONBc := EncodeBlockchainToJSON(bc)
	// bc2 := NewBlockChain()
	// bc2.DecodeBlockchainFromJSON(JSONBc)
	// JSONBc2 := EncodeToBlockchainJSON(bc2)
	// // printStringSlice(JSONBc2)
	// fmt.Println(JSONBc2)
	// fmt.Println("Length of the block chain is : ", bc2.Length)
}

// func test1() {
// 	// test making a genesis block and encoding of a single block
// 	bZero := makeGenesisBlock()
// 	// printBlock(bZero)
// 	encoded := bZero.EncodeToJSON()
// 	//fmt.Println(encoded)
// 	bZero2 := DecodeFromJSON(encoded)
// 	printBlock(bZero2)
// 	fmt.Println(bZero2.EncodeToJSON())
// }

// func test1() {
// 	// test making a genesis block and encoding of a single block
// 	bZero := makeGenesisBlock()
// 	// printBlock(bZero)
// 	encoded := bZero.EncodeToJSON()
// 	//fmt.Println(encoded)
// 	bZero2 := DecodeFromJSON(encoded)
// 	printBlock(bZero2)
// 	fmt.Println(bZero2.EncodeToJSON())
// }
