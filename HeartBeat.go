package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

// HeartBeatData stores bc and node information to send
type HeartBeatData struct {
	ID          string `json:"id"`      // sender's id
	Address     string `json:"address"` // sender's address
	BlockJSON   string `json:"blockjson"`
	PeerMapJSON string `json:"peermapjson"`
}

// NewHeartBeatData creates instance of heart beat
func NewHeartBeatData(id string, address string, blockJSON string, peerMapJSON string) *HeartBeatData {
	return &HeartBeatData{id, address, blockJSON, peerMapJSON}
}

// HBDataToJSON ...
func (hbd *HeartBeatData) HBDataToJSON() ([]byte, error) {
	value, err := json.Marshal(hbd)
	if err != nil {
		return []byte{}, errors.New("Cannot encode HBData ßto Json")
	}
	return value, nil
}

// PeerListToJSON ....
func (pl *PeerList) PeerListToJSON() ([]byte, error) {
	value, err := json.Marshal(pl)
	if err != nil {
		return []byte{}, errors.New("Cannot encode PeerList to Json")
	}
	return value, nil
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
