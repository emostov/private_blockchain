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
	h := sha512.New()
	h.Write([]byte(m))
	out_hash := hex.EncodeToString(h.Sum(nil))
	return out_hash
}

func printBlock(b *Block) {
	h := "height: " + string(b.Header.height) + ", timestamp: " + string(b.Header.timestamp) + ", parent hash: " + b.Header.parent_hash + ", size" + string(b.Header.size)
	value := "Block Value: " + b.Value
	fmt.Println(value)
	fmt.Println(h)
	fmt.Println("___Block End___")
}

type JsonShape struct {
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
	time := int64(time.Now().Unix())
	return &Header{height: height, parent_hash: p_hash, timestamp: time, size: int32(32)}
}

type Block struct {
	Header *Header
	Value  string // root hash of merkle tree
}

// func NewBlock(height int32, p_hash string, value string) *Block {
// 	header := NewHeader(height, p_hash)
// 	return &Block{Header: header, Value: value}
// }

func (b *Block) Initialize(height int32, parent_hash string, value string) {
	b.Value = value
	b.Header = NewHeader(height, parent_hash)
	fmt.Println("block height", b.Header.height)
	hash_str := string(b.Header.height) + string(b.Header.timestamp) + b.Header.parent_hash + string(b.Header.size) + b.Value
	digest := makeSha256Digest(hash_str)
	b.Header.hash = digest
}

func (b *Block) EncodeToJson() string {
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
	var shape JsonShape
	json.Unmarshal([]byte(m), &shape)
	fmt.Println("shape time stamp", shape.Timestamp)
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

func (bc *BlockChain) Get(height int32) []Block {
	if val, ok := bc.Chain[height]; ok {
		return val
	}
	return nil
}

func (bc *BlockChain) Insert(b Block) {
	val, ok := bc.Chain[b.Header.height]
	if ok {
		for i := 0; i < len(val); i++ {
			if val[i] == b {
				return
			}
		}
	}
	bc.Chain[b.Header.height] = append(bc.Chain[b.Header.height], b)
}

func (bc *BlockChain) EncodeToJson() []string {
	var json_blocks []string
	for _, block_slice := range bc.Chain {
		for _, block := range block_slice {
			encoded_block := block.EncodeToJson()
			append(json_blocks, encoded_block)
		}
	}
	return json_block
}

func make_genesis_block() Block {
	p_hash := makeSha256Digest("hash this")
	merkle_root_dummy := makeSha256Digest("root_dummy_hash")
	var b Block
	b.Initialize(0, p_hash, merkle_root_dummy)
	return b
}

func main() {
	fmt.Println("hello world")
	test()
}

func printStringSlice(slice []string) {
	fmt.Println
	for _, json_block := range slice {
		fmt.Println(json_block)
	}
}

func makeTenBlocks() []*Block{
	b_zero := make_genesis_block()
  blocks = []*Blocks
  append(blocks, b_zero)
  for i:=1; i<10; i++{
    var b *Block
    b.Initialize((i%4) + 1, blocks[i-1].Header.hash, "test block value")
    append(blocks, b.Initialize)
  }
}


func test2(){
  var bc BlockChain
  blocks = makeTenBlocks()
}
func test1() {
	b_zero := make_genesis_block()
	// printBlock(b_zero)
	encoded := b_zero.EncodeToJson()
	//fmt.Println(encoded)

	b_zero_2 := DecodeFromJson(encoded)
	printBlock(b_zero_2)
	fmt.Println(b_zero_2.EncodeToJson())
}
