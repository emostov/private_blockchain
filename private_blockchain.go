package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// general utils
func makeSha256Digest(m string) string {
	//takes a message string and returns a string of messages sh256 digest
	h := sha512.New()
	h.Write([]byte(m))
	out_hash := hex.EncodeToString(h.Sum(nil))
	return out_hash
}

type JsonShape struct {
	// for creating proper form when encoding a block to json
	Height      int32
	Timestamp   int64 // Unix timestamp
	Parent_hash string
	Size        int32
	Hash        string
	Value       string
}

type Header struct {
	height      int32
	timestamp   int64 // Unix timestamp
	parent_hash string
	size        int32
	hash        string // hash_str := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + string(b.Header.Size) + b.Value
}

func NewHeader(height int32, p_hash string) *Header {
	// returns a header after given height and parent hash
	// called by the block initialization method
	time := int64(time.Now().Unix())
	return &Header{height: height, parent_hash: p_hash, timestamp: time, size: int32(32)}
}

type Block struct {
	Header *Header
	Value  string // root hash of merkle tree
}

func (b *Block) Initialize(height int32, parent_hash string, value string) {
	// takes a block, its height, its parent_hash, and value and intializes it with
	// header containing hash
	b.Value = value
	b.Header = NewHeader(height, parent_hash)
	hash_str := string(b.Header.height) + string(b.Header.timestamp) + b.Header.parent_hash + string(b.Header.size) + b.Value
	digest := makeSha256Digest(hash_str)
	b.Header.hash = digest
}

func (b *Block) EncodeToJson() string {
	// takes a block pointer and returns an json encoded string of the block
	shape := JsonShape{
		Height:      b.Header.height,
		Timestamp:   b.Header.timestamp,
		Parent_hash: b.Header.parent_hash,
		Size:        b.Header.size,
		Hash:        b.Header.hash,
		Value:       b.Value,
	}
	var encoded []byte
	encoded, err := json.Marshal(shape)
	if err != nil {
		fmt.Println("Error in Json Encoding, ", err)
	}
	return string(encoded)
}

func DecodeFromJson(m string) *Block {
	// takes in a string of a json encoded block and returns a block pointer
	var shape JsonShape
	json.Unmarshal([]byte(m), &shape)
	h := &Header{
		height:      shape.Height,
		timestamp:   shape.Timestamp,
		parent_hash: shape.Parent_hash,
		size:        shape.Size,
		hash:        shape.Hash,
	}
	b := &Block{Header: h, Value: shape.Value}
	return b
}

type BlockChain struct {
	Chain  map[int32][]Block
	Length int32
}

func NewBlockChain() *BlockChain {
	// use this to create and initializea block chain instance
	var bc BlockChain
	bc.Chain = make(map[int32][]Block)
	bc.Length = int32(0)
	return &bc
}

func (bc *BlockChain) Get(height int32) []Block {
	// takes an instance of a block chain and a height in int32
	// returns a slice containing the blocks at that that height or nil
	if val, ok := bc.Chain[height]; ok {
		return val
	}
	return nil
}

func (bc *BlockChain) Insert(b Block) {
	// takes a blockchain instance and inserts a block instance by height
	val, ok := bc.Chain[b.Header.height]
	if ok {
		for i := 0; i < len(val); i++ {
			if val[i] == b {
				return
			}
		}
		bc.Chain[b.Header.height] = append(bc.Chain[b.Header.height], b)
	} else {
		bc.Chain[b.Header.height] = append(bc.Chain[b.Header.height], b)
	}
	if b.Header.height > bc.Length {
		bc.Length = b.Header.height
	}
}

func (bc *BlockChain) EncodeToJson() []string {
	// takes a block chain instance and creates a slice of json block Data
	// returns a the slice of json blocks
	var json_blocks []string
	for _, block_slice := range bc.Chain {
		for _, block := range block_slice {
			encoded_block := block.EncodeToJson()
			json_blocks = append(json_blocks, encoded_block)
		}
	}
	return json_blocks
}

func (bc *BlockChain) DecodeFromJson(json_blocks []string) {
	//takes a blockchain instance and a list of json blocks and inserts each block
	// into the blochchain instance
	for _, json_b := range json_blocks {
		block := DecodeFromJson(json_b)
		bc.Insert(*block)
	}
}

func make_genesis_block() Block {
	//creates and returns a genesis block
	p_hash := makeSha256Digest("hash this")
	merkle_root_dummy := makeSha256Digest("root_dummy_hash")
	var b Block
	b.Initialize(0, p_hash, merkle_root_dummy)
	return b
}

// testing utils
func printBlock(b *Block) {
	// prints blocks fields for debugging and testing purposes only
	h := "height: " + string(b.Header.height) + ", timestamp: " + string(b.Header.timestamp) + ", parent hash: " + b.Header.parent_hash + ", size" + string(b.Header.size)
	value := "Block Value: " + b.Value
	fmt.Println(value)
	fmt.Println(h)
	fmt.Println("___Block End___")
}

func printStringSlice(slice []string) {
	// takes slice of json blocks and prints each one
	fmt.Println("about to print each json block in list")
	for _, json_block := range slice {
		fmt.Println(json_block)
	}
}

func makeTenBlocks() []Block {
	// creates an array of 10 blocks of 5 different heights
	// starts with a genesis block as the only block at height zero
	b_zero := make_genesis_block()
	var blocks []Block
	blocks = append(blocks, b_zero)
	for i := 1; i < 10; i++ {
		var b Block
		height := int32((i % 4) + 1)
		b.Initialize(height, blocks[i-1].Header.hash, "test block value")
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
	bc.Insert(make_genesis_block())
	json_bc := bc.EncodeToJson()
	printStringSlice(json_bc)
}

func test3() {
	// testing creating a blockchain, and block chain encoding and decoding
	bc := NewBlockChain()
	blocks := makeTenBlocks()
	for _, b := range blocks {
		bc.Insert(b)
	}
	json_bc := bc.EncodeToJson()
	bc2 := NewBlockChain()
	bc2.DecodeFromJson(json_bc)
	json_bc2 := bc2.EncodeToJson()
	printStringSlice(json_bc2)
	fmt.Println("Length of the block chain is : ", bc2.Length)
}

func test1() {
	// test making a genesis block and encoding of a single block
	b_zero := make_genesis_block()
	// printBlock(b_zero)
	encoded := b_zero.EncodeToJson()
	//fmt.Println(encoded)
	b_zero_2 := DecodeFromJson(encoded)
	printBlock(b_zero_2)
	fmt.Println(b_zero_2.EncodeToJson())
}
