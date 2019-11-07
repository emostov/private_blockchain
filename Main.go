package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"./uri"
)

// Bc is the blochain instance

var mutex = &sync.Mutex{}
var target = "0000" // 5 0's
var run = true      // for loop conditional in StartTryingNonces()
var id = "6688"
var peerids = []string{"8080"}
var peerlist = PeerList{selfid: id, peerIDs: peerids, length: Bc.Length}

func main() {
	go miner1()
	miner2()
}
func miner1() {

	testSetup()

	//router := NewRouter()
	// go Bc.StartTryingNonces()
	Bc.StartTryingNonces()
	fmt.Println(handlers.Bc.Show())
	router := uri.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6689", router))
	}
}

func miner2() {

	testSetup()

	//router := NewRouter()
	// go Bc.StartTryingNonces()
	Bc.StartTryingNonces()
	fmt.Println(Bc.Show())
	router := NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":8080", router))
	}

	// resp, error := http.Get("http://localhost:6689/block/0/d7c768e1ac640682475a2a6ed935d788e2f8f5fb2c14e9927e97a9fcbb69a7b7386b9b2f6993548bb5a2de899054d4a8428ec40116be12ca872cfcedb0cf6bdc")
	// if error != nil {
	// 	log.Fatalln(error)
	// }
	// body, error := ioutil.ReadAll(resp.Body)
	// if error != nil {
	// 	log.Fatalln(error)
	// }
	// fmt.Println("miner2 got ", body)

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
