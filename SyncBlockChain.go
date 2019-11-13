package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// SyncBlockChain adds a lock to use with the blockchaing
type SyncBlockChain struct {
	BC  *BlockChain `json:"bc"`
	Mux *sync.Mutex `json:"mux"`
}

//Get returns a blockchain instance's Height
// takes an instance of a block chain and a Height in int32
// returns a slice containing the blocks at that that Height or nil
func (sbc *SyncBlockChain) Get(Height int32) []Block {
	// sbc.Mux.Lock()
	// defer sbc.Mux.Unlock()
	if val, ok := sbc.BC.Chain[Height]; ok {
		return val
	}
	return nil
}

// GetBlock ...
func (sbc *SyncBlockChain) GetBlock(height int32, hash string) *Block {
	blocksAtHeight := sbc.Get(height)
	// sbc.Mux.Lock()
	// defer sbc.Mux.Unlock()
	if blocksAtHeight != nil {
		for _, block := range blocksAtHeight {
			if block.Header.Hash == hash {
				return &block
			}
		}
	}
	return nil
}

// NewSyncBlockChain returns an intialized SyncBlockChain pointer
func NewSyncBlockChain() *SyncBlockChain {
	return &SyncBlockChain{BC: NewBlockChain(), Mux: &sync.Mutex{}}
}

// GetLatestBlock returns slice of blocks at chains length
func (sbc *SyncBlockChain) GetLatestBlock() []Block {
	// sbc.Mux.Lock()
	// defer sbc.Mux.Unlock()
	ret := sbc.Get(sbc.BC.Length)
	return ret
}

// GetParentBlock takes a block as a parameter, and returns its parent block
func (sbc *SyncBlockChain) GetParentBlock(b *Block) *Block {
	// sbc.Mux.Lock()
	// defer sbc.Mux.Unlock()
	parentHeightBlocks := sbc.Get(b.Header.Height)

	for _, pBlock := range parentHeightBlocks {
		if pBlock.Header.Hash == b.Header.ParentHash {
			return &pBlock
		}
	}
	return nil
}

//Insert inserts a block into a blockchain, checks for duplicates and updates length
func (sbc *SyncBlockChain) Insert(b Block) {
	fmt.Println("Log: About to attempt insert into block chain")
	sbc.Mux.Lock()
	defer sbc.Mux.Unlock()
	val, ok := sbc.BC.Chain[b.Header.Height]
	if ok {
		for i := 0; i < len(val); i++ {
			if val[i] == b {
				fmt.Println("Log: In Insert and block was not inserted because duplicate")
				return
			}
		}
	}

	sbc.BC.Chain[b.Header.Height] = append(sbc.BC.Chain[b.Header.Height], b)

	if b.Header.Height > sbc.BC.Length {
		sbc.BC.Length = b.Header.Height
	}

	fmt.Println("LOG: post sbc.insert Show() below ")
	//fmt.Println(sbc.BC.Show())

}

//DecodeBlockchainFromJSON ...
func (sbc *SyncBlockChain) DecodeBlockchainFromJSON(JSONBlocks string) {
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
			sbc.Insert(b)
		}
		fmt.Println("LOG: DecodeBlockchainFromJSON: Succesful decode and inserts ")

	} else {
		log.Fatalln(err)
	}

}
