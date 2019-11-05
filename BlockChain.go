package main

import (
	"fmt"
)

// BlockChain struct is used to describe the structure of the blockchain
type BlockChain struct {
	Chain  map[int32][]Block
	Length int32 // length starts at 1
}

//NewBlockChain creates a new blockchain instance, initializing map
func NewBlockChain() *BlockChain {
	// use this to create and initializea block chain instance
	var bc BlockChain
	bc.Chain = make(map[int32][]Block)
	bc.Length = int32(0)
	return &bc
}

//Get returns a blockchain instance's Height
func (bc *BlockChain) Get(Height int32) []Block {
	// takes an instance of a block chain and a Height in int32
	// returns a slice containing the blocks at that that Height or nil
	if val, ok := bc.Chain[Height]; ok {
		return val
	}
	return nil
}

//Insert inserts a block into a blockchain
func (bc *BlockChain) Insert(b Block) {
	// takes a blockchain instance and inserts a block instance by Height
	val, ok := bc.Chain[b.Header.Height]
	if ok {
		for i := 0; i < len(val); i++ {
			if val[i] == b {
				return
			}
		}
	}
	bc.Chain[b.Header.Height] = append(bc.Chain[b.Header.Height], b)
	// } else {
	// 	bc.Chain[b.Header.Height] = append(bc.Chain[b.Header.Height], b)
	// }
	if b.Header.Height+1 > bc.Length {
		bc.Length = b.Header.Height + 1
	}
}

//EncodeToJSON encodes all blocks in chain and puts them into a slice
func (bc *BlockChain) EncodeToJSON() []string {
	// takes a block chain instance and creates a slice of json block Data
	// returns a the slice of json blocks
	var JSONBlocks []string
	for _, blockSlice := range bc.Chain {
		for _, block := range blockSlice {
			encodedBlock := block.EncodeToJSON()
			JSONBlocks = append(JSONBlocks, encodedBlock)
		}
	}
	return JSONBlocks
}

// DecodeFromJSON takes a blockchain instance and a list of json blocks
// and inserts each block into the blochchain instance
func (bc *BlockChain) DecodeFromJSON(JSONBlocks []string) {
	//takes a blockchain instance and a list of json blocks and inserts each block
	// into the blochchain instance
	for _, JSONB := range JSONBlocks {
		block := DecodeFromJSON(JSONB)
		bc.Insert(*block)
	}
}

// GetLatestBlock returns the list of blocks of height "BlockChain.length"
func (bc *BlockChain) GetLatestBlock() []Block {
	return bc.Chain[bc.Length]
}

// GetParentBlock takes a block as a parameter, and returns its parent block
func (bc *BlockChain) GetParentBlock(b *Block) *Block {
	parentHeightBlocks := bc.Get(b.Header.Height)
	for _, pBlock := range parentHeightBlocks {
		if pBlock.Header.Hash == b.Header.ParentHash {
			return &pBlock
		}
	}
	return nil
}

func makeGenesisBlock() Block {
	//creates and returns a genesis block
	pHash := makeSha256Digest("Hash this")
	merkleRootDummy := makeSha256Digest("root_dummy_Hash")
	var b Block
	b.Initialize(0, pHash, merkleRootDummy)
	return b
}

// Testing utils and functions
func printBlock(b *Block) {
	// prints blocks fields for debugging and testing purposes only
	h := "Height: " + string(b.Header.Height) + ", Timestamp: " + string(b.Header.Timestamp) + ", parent Hash: " + b.Header.ParentHash + ", Size" + string(b.Header.Size)
	value := "Block Value: " + b.Value
	fmt.Println(value)
	fmt.Println(h)
	fmt.Println("___Block End___")
}

func printStringSlice(slice []string) {
	// takes slice of json blocks and prints each one
	fmt.Println("about to print each json block in list")
	for _, JSONBlock := range slice {
		fmt.Println(JSONBlock)
	}
}

func makeTenBlocks() []Block {
	// creates an array of 10 blocks of 6 different Heights
	// starts with a genesis block as the only block at Height zero
	Heights := [9]int32{1, 1, 2, 2, 3, 3, 4, 4, 5}
	bZero := makeGenesisBlock()
	var blocks []Block
	blocks = append(blocks, bZero)
	for i := 1; i < 10; i++ {
		var b Block
		// Height := int32((i % 4) + 1)
		// naive parent Hash, not actually accurate to chain
		b.Initialize(Heights[i-1], blocks[i-1].Header.Hash, "test block value")
		blocks = append(blocks, b)
	}
	return blocks
}

// tests
func main() {
	fmt.Println("Beggining of main")
	test3()
}

func test2() {
	// testing insertion of a block into the block chain
	bc := NewBlockChain()
	bc.Insert(makeGenesisBlock())
	JSONBc := bc.EncodeToJSON()
	printStringSlice(JSONBc)
}

func test3() {
	// testing creating a blockchain, and block chain encoding and decoding
	bc := NewBlockChain()
	blocks := makeTenBlocks()
	for _, b := range blocks {
		bc.Insert(b)
	}
	JSONBc := bc.EncodeToJSON()
	bc2 := NewBlockChain()
	bc2.DecodeFromJSON(JSONBc)
	JSONBc2 := bc2.EncodeToJSON()
	printStringSlice(JSONBc2)
	fmt.Println("Length of the block chain is : ", bc2.Length)
}

func test1() {
	// test making a genesis block and encoding of a single block
	bZero := makeGenesisBlock()
	// printBlock(bZero)
	encoded := bZero.EncodeToJSON()
	//fmt.Println(encoded)
	bZero2 := DecodeFromJSON(encoded)
	printBlock(bZero2)
	fmt.Println(bZero2.EncodeToJSON())
}
