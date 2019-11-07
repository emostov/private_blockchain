package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// Bc is the blochain instance
var Bc = NewBlockChain()

// Upload ...
func Upload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(EncodeBlockchainToJSON(Bc)))
}

// AskForBlock ...
func AskForBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h := vars["height"]
	hash := vars["hash"]
	// height, err := strconv.Atoi(h)

	// if err == nil {
	// 	block := Bc.GetBlock(int32(height), hash)
	// 	if block == nil {
	// 		w.WriteHeader(http.StatusNotFound)
	// 		// Ask for parent block, insert current block into tree
	// 	} else {
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write(block.EncodeToJSON())
	// 	}
	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// }

	w.Write([]byte("height :" + h + ", hash: " + hash))
}

// HeartBeatRecieve ...
func HeartBeatRecieve(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		// results = append(results, string(body))
		mutex.Lock()
		run = false
		fmt.Println(string(body))
		// this is WRONG will be in form of heartbeat data
		s := string(body)
		data := HeartBeatData{}
		json.Unmarshal([]byte(s), &data)
		block := DecodeFromJSON(string(data.blockJSON))
		if verifyNonce(block) {
			Bc.Insert(*block)
		}
		mutex.Unlock()
		fmt.Fprint(w, "POST done")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

// ShowHandler ...
func ShowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Bc.Show()))
}

// HeartBeatReceive()
//  When a node receives a new block in HeartBeat, the node will first check if
//  the parent block of this new block exists in its own blockchain (the previous
// 	block is the block whose hash is the parentHash of the next block).
// If the previous block doesn't exist, the node will ask any peer in PeerList at
// "/block/{height}/{hash}" to download that block.
// After making sure the previous block exists, insert the block from HeartBeatData
//  to the current BlockChain.
//  Alter this function so that when it receives a HeartBeatData with a new block,
//  it verifies the nonce as described above.

// func HeartBeatRecieve() (w http.ResponseWriter, r *http.Request) {

// }
