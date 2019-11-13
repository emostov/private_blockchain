package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// ServerRegisterData ...
type ServerRegisterData struct {
	ServerID    ID     `json:"serverid"`
	PeerMapJSON string `json:"peermapjson"`
	PeerMap     []ID   `json:"peermap"`
}

// RegisterData ...
type RegisterData struct {
	PeerMapJSON string `json:"peermapjson"`
	AssignedID  string `json:"assignedid"` //string of id struct
}

// ID ..
type ID struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

// PeerList ...
type PeerList struct {
	SelfID  ID    `json:"selfid"`
	PeerIDs []ID  `json:"peerids"`
	Length  int32 `json:"length"`
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
	pl.Length = int32(len(pl.PeerIDs))
}

// RemoveSelfFromPeerList ...
func (pl *PeerList) RemoveSelfFromPeerList() {
	for i, id := range pl.PeerIDs {
		if id == MINERID {
			pl.PeerIDs = append(pl.PeerIDs[:i], pl.PeerIDs[i+1:]...)
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
	if len(srd.PeerMap) >= 1 {

		for _, id := range srd.PeerMap {
			peermapjs += (id.EncodeIDToJSON() + ",")
		}
		peermapjs = peermapjs[:len(peermapjs)-1] + "]"
	} else {
		peermapjs += "]"
	}
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
	resp, err := http.Post(SID.Address+SID.Port+"/peer", "application/json", bytes.NewBuffer([]byte(IDJSON)))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == http.StatusOK {
		rd := DecodeRegisterDataFromJSON(string(body))
		peerids := DecodePeerMapFromJSON(rd.PeerMapJSON)
		PEERLIST.AddNewPeers(peerids)
		PEERLIST.RemoveSelfFromPeerList()
		fmt.Println("LOG: I just registered and my peer list is: ", PEERLIST.PeerIDs)
	}

}

// peerlist insert peerids
