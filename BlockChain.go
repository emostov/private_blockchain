package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// BlockChain struct is used to describe the structure of the blockchain
type BlockChain struct {
	Chain  map[int32][]Block `json:"chain"`
	Length int32             `json:"length"` // length starts at 0
}

//NewBlockChain creates a new blockchain instance, initializing map
func NewBlockChain() *BlockChain {
	// use this to create and initializea block chain instance
	var bc BlockChain
	bc.Chain = make(map[int32][]Block)
	bc.Length = int32(0)
	return &bc
}

//NonSyncGet  returns a blockchain instance's Height
func (bc *BlockChain) NonSyncGet(Height int32) []Block {
	// takes an instance of a block chain and a Height in int32
	// returns a slice containing the blocks at that that Height or nil
	if val, ok := bc.Chain[Height]; ok {
		return val
	}
	return nil
}

// NonSyncGetBlock ...
func (bc *BlockChain) NonSyncGetBlock(height int32, hash string) *Block {
	blocksAtHeight := bc.NonSyncGet(height)
	if blocksAtHeight != nil {
		for _, block := range blocksAtHeight {
			if block.Header.Hash == hash {
				return &block
			}
		}
	}
	return nil
}

// NonSyncInsert inserts a block into a blockchain
func (bc *BlockChain) NonSyncInsert(b Block) {
	fmt.Println("Log: About to attempt insert into block chain")
	val, ok := bc.Chain[b.Header.Height]
	if ok {
		for i := 0; i < len(val); i++ {
			if val[i] == b {
				fmt.Println("Log: In Insert and block was not inserted because duplicate")
				return
			}
		}
	}

	bc.Chain[b.Header.Height] = append(bc.Chain[b.Header.Height], b)

	if b.Header.Height > bc.Length {
		bc.Length = b.Header.Height
	}

	fmt.Println("LOG: post bc.insert Show() below ")
	fmt.Println(bc.Show())

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

// type JsonType struct {
// 	Array []string
// }

//NonSyncDecodeBlockchainFromJSON ...
func (bc *BlockChain) NonSyncDecodeBlockchainFromJSON(JSONBlocks string) {
	//takes a blockchain instance and a list of json blocks and inserts each block
	// into the blochchain instance
	var blockList []JSONShape
	err := json.Unmarshal([]byte(JSONBlocks), &blockList)
	if err == nil {
		for _, shape := range blockList {
			// create block from json
			h := Header{
				Nonce:      shape.Nonce,
				Height:     shape.Height,
				Timestamp:  shape.Timestamp,
				ParentHash: shape.ParentHash,
				Size:       shape.Size,
				Hash:       shape.Hash,
			}
			b := Block{Header: h, Value: shape.Value}
			bc.NonSyncInsert(b)
		}
		fmt.Println("no error")

	} else {
		log.Fatalln(err)
	}

}

// NonSyncGetLatestBlock returns the list of blocks of height "BlockChain.length"
func (bc *BlockChain) NonSyncGetLatestBlock() []Block {
	return bc.Chain[bc.Length]
}

// NonSyncGetParentBlock takes a block as a parameter, and returns its parent block
func (bc *BlockChain) NonSyncGetParentBlock(b *Block) *Block {
	parentHeightBlocks := bc.NonSyncGet(b.Header.Height)
	for _, pBlock := range parentHeightBlocks {
		if pBlock.Header.Hash == b.Header.ParentHash {
			return &pBlock
		}
	}
	return nil
}

// makeGenesisBlock makes a dummy block
func makeGenesisBlock() Block {
	//creates and returns a genesis block
	pHash := makeSha256Digest("Hash this")
	merkleRootDummy := makeSha256Digest("root_dummy_Hash")
	var b Block
	b.Initialize(0, pHash, merkleRootDummy, 0)
	return b
}

// Testing utils and functions

func printStringSlice(slice []string) {
	// takes slice of json blocks and prints each one
	fmt.Println("about to print each json block in list")
	for _, JSONBlock := range slice {
		fmt.Println(JSONBlock)
	}
}

func makeTenBlocks() []Block {
	// creates an array of  blocks of the same height
	// starts with a genesis block as the only block at Height zero
	bZero := makeGenesisBlock()
	var blocks []Block
	blocks = append(blocks, bZero)
	for i := 1; i < 10; i++ {
		var b Block
		b.Initialize(int32(1), bZero.Header.Hash, "test block value", int32(i))
		blocks = append(blocks, b)
	}
	return blocks
}
