package main

import (
	"crypto/sha256"
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

func makeSha256Digest(m string) string {
	//takes a message string and returns a string of messages sh256 digest
	h := sha256.New()
	h.Write([]byte(m))
	outHash := hex.EncodeToString(h.Sum(nil))
	return outHash
}

// JSONShape is a struct used to help in conversion to json
type JSONShape struct {
	// for creating proper form when encoding a block to json
	Difficulty int32  `json:"difficulty"` // target
	Nonce      string `json:"nonce"`
	Height     int32  `json:"height"`
	Timestamp  int64  `json:"timestamp"` // Unix Timestamp
	ParentHash string `json:"parenthash"`
	Size       int32  `json:"size"`
	Hash       string `json:"hash"`
	Value      string `json:"value"`
}

//Header Struct describing the fields of the header
type Header struct {
	Difficulty int32  `json:"difficulty"` // target
	Nonce      string `json:"nonce"`
	Height     int32  `json:"height"`
	Timestamp  int64  `json:"timestamp"` // Unix Timestamp
	ParentHash string `json:"parenthash"`
	Size       int32  `json:"size"`
	Hash       string `json:"hash"` // HashStr := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + string(b.Header.Size) + b.Value
}

//Block struct with header pointer and value
type Block struct {
	Header Header
	Value  string // root Hash of merkle tree
}

//Header tools

// NewHeader is used to create and initialize a heder
// returns a header after given Height and parent Hash
// called by the block initialization method
func NewHeader(Height int32, pHash string, difficulty int32) Header {
	time := int64(time.Now().Unix())
	return Header{Height: Height, ParentHash: pHash, Timestamp: time, Size: int32(32), Difficulty: difficulty}
}

// Initialize an instance of a block
// takes a block, its Height, its ParentHash, and value and intializes it with
// header containing Hash
func (b *Block) Initialize(Height int32, ParentHash string, value string, difficulty int32) {
	b.Value = value
	b.Header = NewHeader(Height, ParentHash, difficulty)
	HashStr := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + string(b.Header.Size) + b.Value
	digest := makeSha256Digest(HashStr)
	b.Header.Hash = digest
}

// EncodeToJSON a block instance
// takes a block pointer and returns an json encoded string of the block
func (b *Block) EncodeToJSON() []byte {
	shape := JSONShape{
		Difficulty: b.Header.Difficulty,
		Nonce:      b.Header.Nonce,
		Height:     b.Header.Height,
		Timestamp:  b.Header.Timestamp,
		ParentHash: b.Header.ParentHash,
		Size:       b.Header.Size,
		Hash:       b.Header.Hash,
		Value:      b.Value,
	}
	var encoded []byte
	encoded, err := json.Marshal(shape)
	if err != nil {
		fmt.Println("Error in Json Encoding, ", err)
	}
	///return string(encoded)
	return encoded
}

// DecodeFromJSON takes a json string of a block and converts into a block instance
// takes in a string of a json encoded block and returns a block pointer
func DecodeFromJSON(m string) *Block {
	var shape JSONShape
	json.Unmarshal([]byte(m), &shape)
	h := Header{
		Difficulty: shape.Difficulty,
		Nonce:      shape.Nonce,
		Height:     shape.Height,
		Timestamp:  shape.Timestamp,
		ParentHash: shape.ParentHash,
		Size:       shape.Size,
		Hash:       shape.Hash,
	}
	b := &Block{Header: h, Value: shape.Value}
	return b
}
