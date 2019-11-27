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
	sbc.Mux.Lock()
	defer sbc.Mux.Unlock()
	if val, ok := sbc.BC.Chain[Height]; ok {
		return val
	}
	return nil
}

// GetBlock uses a lock and the rest of the retrieval methods for SyncBC
// use get to interact with the chain
func (sbc *SyncBlockChain) GetBlock(height int32, hash string) *Block {
	blocksAtHeight := sbc.Get(height)
	sbc.Mux.Lock()
	defer sbc.Mux.Unlock()
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
	ret := sbc.Get(sbc.BC.Length)
	return ret
}

// GetParentBlock takes a block as a parameter, and returns its parent block
func (sbc *SyncBlockChain) GetParentBlock(b *Block) *Block {
	if !(b.Header.Height >= int32(1)) {
		log.Println("LOG: GetParentBlock: err Block height zero, so no parent")
		return nil
	}
	parentHeightBlocks := sbc.Get(b.Header.Height - int32(1))
	for _, pBlock := range parentHeightBlocks {
		if pBlock.Header.Hash == b.Header.ParentHash {
			return &pBlock
		}
	}
	log.Println("LOG: GetParentBlock: calling askforparent could not find parent block")
	askForParent(b.Header.ParentHash, fmt.Sprint(b.Header.Height-int32(1)))
	for _, pBlock := range parentHeightBlocks {
		if pBlock.Header.Hash == b.Header.ParentHash {
			return &pBlock
		}
	}
	log.Println("Log: GetParentBlock: could not find paretn block")
	return nil
}

//Insert inserts a block into a blockchain, checks for duplicates and updates length
// Also updates the blocks difficulty based on the parents block difficulty
// If a parent block cannot be found no insert will happen
func (sbc *SyncBlockChain) Insert(b Block) {
	log.Println("Log: Insert: Begin insert attempt")
	if b.Header.Height >= 1 {
		if sbc.GetParentBlock(&b) != nil {
			b.Header.Difficulty = int32(len(TARGET)) + sbc.GetParentBlock(&b).Header.Difficulty
		} else {
			log.Println("LOG: Insert: could not insert because not parent found")
			return
		}
	}
	sbc.Mux.Lock()
	defer sbc.Mux.Unlock()
	val, ok := sbc.BC.Chain[b.Header.Height]
	if ok {
		for i := 0; i < len(val); i++ {
			if val[i] == b {
				log.Println("LOG: Insert: block was not inserted because duplicate")
				return
			}
		}
	}
	sbc.BC.Chain[b.Header.Height] = append(sbc.BC.Chain[b.Header.Height], b)
	log.Println("LOG: Insert: succesful insert for: " + b.Header.Hash)
	if b.Header.Height > sbc.BC.Length {
		sbc.BC.Length = b.Header.Height
	}
}

//DecodeBlockchainFromJSON ...
// takes a blockchain instance and a list of json blocks and inserts each block
// into the blochchain instance
func (sbc *SyncBlockChain) DecodeBlockchainFromJSON(JSONBlocks string) {
	log.Println("LOG: DecodeBlockchainFromJSON: about to decode and insert")
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
		log.Println("LOG: DecodeBlockchainFromJSON: Succesful decode and inserts ")
	} else {
		log.Fatalln(err)
	}
}
