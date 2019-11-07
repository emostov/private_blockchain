package main

import (
	"encoding/json"
	"fmt"
)

// HeartBeatData stores bc and node information to send
type HeartBeatData struct {
	id          string // sender's id
	address     string // sender's address
	blockJSON   string
	peerMapJSON string
}

// NewHeartBeatData creates instance of heart beat
func NewHeartBeatData(id string, address string, blockJSON string, peerMapJSON string) *HeartBeatData {
	return &HeartBeatData{id, address, blockJSON, peerMapJSON}
}

// EncodeToJSON converts heart beat data to json string
func (hbd *HeartBeatData) EncodeToJSON() string {
	var encoded []byte
	encoded, err := json.Marshal(hbd)
	if err != nil {
		fmt.Println("Error in Json Encoding, ", err)
	}
	return string(encoded)
}
