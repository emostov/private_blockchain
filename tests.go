import "fmt"

func testSetupBlockInsert() {
	var blocks = makeTenBlocks()
	Bc.Insert(blocks[0])
	Bc.Insert(blocks[1])
	Bc.Insert(blocks[2])
}

func testDecode() {
	bc := NewBlockChain()
	var blocks = makeTenBlocks()
	// genesis := makeGenesisBlock()
	// bc.Insert(genesis)
	bc.Insert(blocks[0])
	bc.Insert(blocks[1])
	//bc.Insert(blocks[2])
	fmt.Println("bc ", bc.Show())
	blocks1 := EncodeBlockchainToJSON(bc)
	Bc.DecodeBlockchainFromJSON(blocks1)
	fmt.Println("Bc ", Bc.Show())
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

func test1() {
	// block encoding and decoding test
	var b Block
	b.Initialize(int32(0), "parent hash", "merkle value")
	coded := b.EncodeToJSON()
	fmt.Println(" 1 coded block is ", string(coded))
	decoded := DecodeFromJSON(string(coded))
	fmt.Println("decoded block is ", decoded)
	coded2 := decoded.EncodeToJSON()
	fmt.Println("coded 2 ", string(coded2))

}

func test3() {
	// testing creating a blockchain, and block chain encoding and decoding
	bc := NewBlockChain()
	blocks := makeTenBlocks()
	for _, b := range blocks {
		bc.Insert(b)
	}
	fmt.Println(bc.Show())
}