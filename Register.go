package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

// PeerList ...
type PeerList struct {
	SelfID  ID
	PeerIDs []ID
	Length  int32
}

func contains(PeerIDs []ID, otherID ID) bool {
	for _, a := range PeerIDs {
		if a == otherID {
			return true
		}
	}
	return false
}

// AddNewPeers ...
func (pl *PeerList) AddNewPeers(newPeers []ID) {
	for _, pID := range newPeers {
		if !contains(pl.PeerIDs, pID) {
			pl.PeerIDs = append(pl.PeerIDs, pID)
		}
	}
}

// AddNewPeer ...
func (srd *ServerRegisterData) AddNewPeer(id ID) {
	// for _, pID := range newPeers {
	if !contains(srd.PeerMap, id) {
		srd.PeerMap = append(srd.PeerMap, id)
	}
	// }
	srd.EncodePeerMapToJSON()
}

//EncodePeerMapToJSON ...
func (srd *ServerRegisterData) EncodePeerMapToJSON() {
	peermapjs := "["
	for _, id := range srd.PeerMap {
		peermapjs += (id.EncodeIDToJSON() + ",")
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

// DecodeRegisterDataFromJSON ...
func DecodeRegisterDataFromJSON(rdjson string) RegisterData {
	var rd RegisterData
	if err := json.Unmarshal([]byte(rdjson), &rd); err != nil {
		log.Fatal(err)
	}
	return rd
}

// EncodeIDToJSON ...
func (id *ID) EncodeIDToJSON() string {
	encodedID, err := json.Marshal(id)
	if err != nil {
		log.Fatal(err)
	}
	return string(encodedID)
}

// DecodeIDFromJSON ...
func DecodeIDFromJSON(idjs string) ID {
	var id ID
	if err := json.Unmarshal([]byte(idjs), &id); err != nil {
		panic(err)
	}
	return id
}

// DecodePeerMapFromJSON ...
func DecodePeerMapFromJSON(peermapjs string) []ID {
	var peerlist []ID
	if err := json.Unmarshal([]byte(peermapjs), &peerlist); err != nil {
		panic(err)
	}
	return peerlist
}

// DoMinerRegistration register the miner with the server and updates the peer list
func DoMinerRegistration() {
	IDJSON := MINERID.EncodeIDToJSON()
	resp, err := http.Post(MINERID.Address+MINERID.Port+"/peer", "application/json", bytes.NewBuffer([]byte(IDJSON)))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		rd := DecodeRegisterDataFromJSON(string(body))
		peerids := DecodePeerMapFromJSON(rd.PeerMapJSON)
		PEERLIST.AddNewPeers(peerids)
	}
}

// peerlist insert peerids
