package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mutex = &sync.Mutex{}
var target = "0000" // 5 0's
var run = true

func main() {

	testSetup()

	router := NewRouter()
	// go Bc.StartTryingNonces()
	Bc.StartTryingNonces()
	fmt.Println(Bc.Show())
	Bc.StartTryingNonces()
	fmt.Println(Bc.Show())
	log.Fatal(http.ListenAndServe(":8080", router))

	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/Upload", Upload)
	// router.HandleFunc("/block/{height}/{hash}", AskForBlock)
	// router.HandleFunc("/Show", ShowHandler)
	// router.HandleFunc("/heartbeat/recieve", HeartBeatRecieve)
	// log.Fatal(http.ListenAndServe(":8080", router))

}

func testSetup() {
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
