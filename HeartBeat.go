package main

import (
	"encoding/json"
	"fmt"
)

type HeartBeatData struct {
	id          string // sender's id
	address     string // sender's address
	blockJSON   string
	peerMapJSON string
}

func NewHeartBeatData(id string, address string, blockJSON string, peerMapJSON string) *HeartBeatData {
	return &HeartBeatData{id, address, blockJSON, peerMapJSON}
}

func (hbd *HeartBeatData) EncodeToJson() string {
	var encoded []byte
	encoded, err := json.Marshal(hbd)
	if err != nil {
		fmt.Println("Error in Json Encoding, ", err)
	}
	return string(encoded)
}
