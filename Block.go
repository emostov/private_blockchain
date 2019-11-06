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

func makeSha256Digest(m string) string {
	//takes a message string and returns a string of messages sh256 digest
	h := sha512.New()
	h.Write([]byte(m))
	outHash := hex.EncodeToString(h.Sum(nil))
	return outHash
}

//JSONShape is a struct used to help in conversion to json
type JSONShape struct {
	// for creating proper form when encoding a block to json
	Nonce      string
	Height     int32
	Timestamp  int64 // Unix Timestamp
	ParentHash string
	Size       int32
	Hash       string
	Value      string
}

//Header Struct describing the fields of the header
type Header struct {
	Nonce      string
	Height     int32
	Timestamp  int64 // Unix Timestamp
	ParentHash string
	Size       int32
	Hash       string // HashStr := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + string(b.Header.Size) + b.Value
}

//Block struct with header pointer and value
type Block struct {
	Header *Header
	Value  string // root Hash of merkle tree
}

//Header tools

// NewHeader is used to create and initialize a heder
func NewHeader(Height int32, pHash string) *Header {
	// returns a header after given Height and parent Hash
	// called by the block initialization method
	time := int64(time.Now().Unix())
	return &Header{Height: Height, ParentHash: pHash, Timestamp: time, Size: int32(32)}
}

//Block Methods

// Initialize an instance of a block
func (b *Block) Initialize(Height int32, ParentHash string, value string) {
	// takes a block, its Height, its ParentHash, and value and intializes it with
	// header containing Hash
	b.Value = value
	b.Header = NewHeader(Height, ParentHash)
	HashStr := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + string(b.Header.Size) + b.Value
	digest := makeSha256Digest(HashStr)
	b.Header.Hash = digest
}

// EncodeToJSON a block instance
func (b *Block) EncodeToJSON() string {
	// takes a block pointer and returns an json encoded string of the block
	shape := JSONShape{
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
	return string(encoded)
}

//DecodeFromJSON takes a json string of a block and converts into a block instance
func DecodeFromJSON(m string) *Block {
	// takes in a string of a json encoded block and returns a block pointer
	var shape JSONShape
	json.Unmarshal([]byte(m), &shape)
	h := &Header{
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

// SetNonce recomputes hash and sets nonce for block
func (b *Block) SetNonce(nonce string) {
	b.Header.Nonce = nonce
	HashStr := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + string(b.Header.Size) + b.Value
	digest := makeSha256Digest(HashStr)
}
