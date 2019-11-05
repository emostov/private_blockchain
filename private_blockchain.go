package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

/*
To create a block declare and instance and then use the block initialize method.
To create a blockchain use the NewBlockChain() method create an instance that is
initialized with an empty map.

Blocks supports Initialize, DecodeFromJSON and EncodeFromJson.

BlockChain supports Get, Insert, EncodeToJSON, and DecodeFromJSON. BlockChain
length starts at 1

makeGenesisBlock() is useful for making an arbitrary genesis block for a
Blockchain.
*/

// general utils
func makeSha256Digest(m string) string {
	//takes a message string and returns a string of messages sh256 digest
	h := sha512.New()
	h.Write([]byte(m))
	outHash := hex.EncodeToString(h.Sum(nil))
	return outHash
}

// BlockChain struct is used to describe the structure of the blockchain
type BlockChain struct {
	Chain  map[int32][]Block
	Length int32 // length starts at 1
}

//JSONShape is a struct used to help in conversion to json
type JSONShape struct {
	// for creating proper form when encoding a block to json
	Height     int32
	Timestamp  int64 // Unix timestamp
	ParentHash string
	Size       int32
	Hash       string
	Value      string
}

//Header Struct describing the fields of the header
type Header struct {
	height     int32
	timestamp  int64 // Unix timestamp
	ParentHash string
	size       int32
	hash       string // hashStr := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + string(b.Header.Size) + b.Value
}

//Block struct with header pointer and value
type Block struct {
	Header *Header
	Value  string // root hash of merkle tree
}

//Header tools

// NewHeader is used to create and initialize a heder
func NewHeader(height int32, pHash string) *Header {
	// returns a header after given height and parent hash
	// called by the block initialization method
	time := int64(time.Now().Unix())
	return &Header{height: height, ParentHash: pHash, timestamp: time, size: int32(32)}
}

//Block Methods

// Initialize an instance of a block
func (b *Block) Initialize(height int32, ParentHash string, value string) {
	// takes a block, its height, its ParentHash, and value and intializes it with
	// header containing hash
	b.Value = value
	b.Header = NewHeader(height, ParentHash)
	hashStr := string(b.Header.height) + string(b.Header.timestamp) + b.Header.ParentHash + string(b.Header.size) + b.Value
	digest := makeSha256Digest(hashStr)
	b.Header.hash = digest
}

// EncodeToJSON a block instance
func (b *Block) EncodeToJSON() string {
	// takes a block pointer and returns an json encoded string of the block
	shape := JSONShape{
		Height:     b.Header.height,
		Timestamp:  b.Header.timestamp,
		ParentHash: b.Header.ParentHash,
		Size:       b.Header.size,
		Hash:       b.Header.hash,
		Value:      b.Value,
	}
	var encoded []byte
	encoded, err := json.Marshal(shape)
	if err != nil {
		fmt.Println("Error in Json Encoding, ", err)
	}
	return string(encoded)
}

//DecodeFromJSON takes a json string of a block and converts into a block instance
func DecodeFromJSON(m string) *Block {
	// takes in a string of a json encoded block and returns a block pointer
	var shape JSONShape
	json.Unmarshal([]byte(m), &shape)
	h := &Header{
		height:     shape.Height,
		timestamp:  shape.Timestamp,
		ParentHash: shape.ParentHash,
		size:       shape.Size,
		hash:       shape.Hash,
	}
	b := &Block{Header: h, Value: shape.Value}
	return b
}

//BlockChain Methods and tools

//NewBlockChain creates a new blockchain instance, initializing map
func NewBlockChain() *BlockChain {
	// use this to create and initializea block chain instance
	var bc BlockChain
	bc.Chain = make(map[int32][]Block)
	bc.Length = int32(0)
	return &bc
}

//Get returns a blockchain instance's height
func (bc *BlockChain) Get(height int32) []Block {
	// takes an instance of a block chain and a height in int32
	// returns a slice containing the blocks at that that height or nil
	if val, ok := bc.Chain[height]; ok {
		return val
	}
	return nil
}

//Insert inserts a block into a blockchain
func (bc *BlockChain) Insert(b Block) {
	// takes a blockchain instance and inserts a block instance by height
	val, ok := bc.Chain[b.Header.height]
	if ok {
		for i := 0; i < len(val); i++ {
			if val[i] == b {
				return
			}
		}
	}
	bc.Chain[b.Header.height] = append(bc.Chain[b.Header.height], b)
	// } else {
	// 	bc.Chain[b.Header.height] = append(bc.Chain[b.Header.height], b)
	// }
	if b.Header.height+1 > bc.Length {
		bc.Length = b.Header.height + 1
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

func makeGenesisBlock() Block {
	//creates and returns a genesis block
	pHash := makeSha256Digest("hash this")
	merkleRootDummy := makeSha256Digest("root_dummy_hash")
	var b Block
	b.Initialize(0, pHash, merkleRootDummy)
	return b
}

// Testing utils and functions
func printBlock(b *Block) {
	// prints blocks fields for debugging and testing purposes only
	h := "height: " + string(b.Header.height) + ", timestamp: " + string(b.Header.timestamp) + ", parent hash: " + b.Header.ParentHash + ", size" + string(b.Header.size)
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
	// creates an array of 10 blocks of 6 different heights
	// starts with a genesis block as the only block at height zero
	heights := [9]int32{1, 1, 2, 2, 3, 3, 4, 4, 5}
	bZero := makeGenesisBlock()
	var blocks []Block
	blocks = append(blocks, bZero)
	for i := 1; i < 10; i++ {
		var b Block
		// height := int32((i % 4) + 1)
		// naive parent hash, not actually accurate to chain
		b.Initialize(heights[i-1], blocks[i-1].Header.hash, "test block value")
		blocks = append(blocks, b)
	}
	return blocks
}

// tests
func main() {
	//fmt.Println("Beggining of main")
	//test3()
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
