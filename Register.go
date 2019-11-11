package main

import (
	"encoding/json"
	"log"
)

// ServerRegisterData ...
type ServerRegisterData struct {
	ServerID    ID
	PeerMapJSON string
	PeerMap     []ID
}

// RegisterData ...
type RegisterData struct {
	PeerMapJSON string
	AssignedID  string
}

// ID ..
type ID struct {
	Address string
	Port    string
}

type PeerList struct {
	SelfID  ID
	PeerIDs []ID
	Length  int32
}

func (pl *PeerList) contains(otherID string) bool {
	for _, a := range pl.PeerIDs {
		if a == otherID {
			return true
		}
	}
	return false
}

//EncodePeerMapToJSON ...
func (srd *ServerRegisterData) EncodePeerMapToJSON() {
	peermapjs := "["
	for _, id := range srd.PeerMap {
		peermapjs += (id + ",")
	}
	peermapjs = peermapjs[:len(peermapjs)-1] + "]"
	srd.PeerMapJSON = peermapjs
}

//EncodeRegisterDataToJSON ...
func (rd *RegisterData) EncodeRegisterDataToJSON() []byte {
	encodedRD, err := json.Marshal(rd)
	if err != nil {
		log.Fatal(err)
	}
	return encodedRD
}

// EncodeIDToJSON ...
func (id *ID) EncodeIDToJSON() string {
	encodedID, err := json.Marshal(id)
	if err != nil {
		log.Fatal(err)
	}
	return string(encodedID)
}

// // RegisterMe
// func RegisterMe() {

// }
