package main

import (
	"encoding/json"
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

// GetBlock ...
func (bc *BlockChain) GetBlock(height int32, hash string) *Block {
	blocksAtHeight := bc.Get(height)
	if blocksAtHeight != nil {
		for _, block := range blocksAtHeight {
			if block.Header.Hash == hash {
				return &block
			}
		}
	}
	return nil
}

//Insert inserts a block into a blockchain
func (bc *BlockChain) Insert(b Block) {
	//mutex.Lock()
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
	if b.Header.Height > bc.Length {
		bc.Length = b.Header.Height
	}
	//mutex.Unlock()
	fmt.Println("Insert Success")
}

//EncodeToJSON encodes all blocks in chain and puts them into a slice
func (bc *BlockChain) EncodeToJSON() []string {
	// takes a block chain instance and creates a slice of json block Data
	// returns a the slice of json blocks
	var JSONBlocks []string
	for _, blockSlice := range bc.Chain {
		for _, block := range blockSlice {
			encodedBlock := block.EncodeToJSON()
			JSONBlocks = append(JSONBlocks, string(encodedBlock))
		}
	}
	return JSONBlocks
}

//EncodeBlockchainToJSON encodes all blocks in chain and puts them into a slice
func EncodeBlockchainToJSON(bc *BlockChain) string {
	// takes a block chain instance and creates a slice of json block Data
	// returns a the slice of json blocks
	var JSONBlocks string
	JSONBlocks = "["
	for _, blockSlice := range bc.Chain {
		for _, block := range blockSlice {
			encodedBlock := block.EncodeToJSON()
			JSONBlocks += string(encodedBlock)
			JSONBlocks += ","
		}
	}
	JSONBlocks = JSONBlocks[:len(JSONBlocks)-1]
	JSONBlocks += "]"
	return JSONBlocks
}

// DecodeFromJSON takes a blockchain instance and a list of json blocks
// and inserts each block into the blochchain instance
// func (bc *BlockChain) DecodeFromJSON(JSONBlocks string) {
// 	//takes a blockchain instance and a list of json blocks and inserts each block
// 	// into the blochchain instance
// 	//JSONBlocks = []Block(JSONBlocks)
// 	for _, JSONB := range JSONBlocks {
// 		block := DecodeFromJSON(JSONB)
// 		bc.Insert(*block)
// 	}
// }

//DecodeBlockchainFromJSON ...
func (bc *BlockChain) DecodeBlockchainFromJSON(JSONBlocks string) {
	//takes a blockchain instance and a list of json blocks and inserts each block
	// into the blochchain instance
	blockList := []Block{}
	err := json.Unmarshal([]byte(JSONBlocks), &blockList)
	if err == nil {
		for _, JSONB := range blockList {
			// block := DecodeFromJSON(JSONB)
			// bc.Insert(*block)
			bc.Insert(JSONB)
		}

	}

}

// GetLatestBlock returns the list of blocks of height "BlockChain.length"
func (bc *BlockChain) GetLatestBlock() []Block {
	// fmt.Println("Length of bc: ", bc.Length)
	// fmt.Println("array of block at length", bc.Chain[bc.Length])
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
	// creates an array of  blocks of the same
	// starts with a genesis block as the only block at Height zero
	Heights := [9]int32{1, 1, 2, 2, 3, 3, 4, 4, 5}
	bZero := makeGenesisBlock()
	var blocks []Block
	blocks = append(blocks, bZero)
	for i := 1; i < 10; i++ {
		var b Block
		// Height := int32((i % 4) + 1)
		// naive parent Hash, not actually accurate to chain
		b.Initialize(Heights[1], bZero.Header.Hash, "test block value")
		blocks = append(blocks, b)
	}
	return blocks
}
