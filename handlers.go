package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Bc ...
var Bc = NewBlockChain()

var mutex = &sync.Mutex{}

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
	height, err := strconv.ParseInt(h, 10, 64)
	fmt.Println("i am ", SELFID, " ask get", "height :"+h+", hash: "+hash)
	if err == nil {
		block := Bc.GetBlock(int32(height), hash)
		if block == nil {
			w.WriteHeader(http.StatusNotFound)
			// Ask for parent block, insert current block into tree
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(block.EncodeToJSON())
			// check if parent exists
			// if it doesnt ask for it?
		}
	}
	// else {
	// 	w.WriteHeader(http.StatusNotFound)
	// }

	//w.Write([]byte("height :" + h + ", hash: " + hash))
}

// HeartBeatRecieve ...
func HeartBeatRecieve(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		requestBody, err := readRequestBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)

		mutex.Lock()
		//w.Write([]byte("In HeartBeat Recieve  "))
		fmt.Println("Im an id and am getting a post", SELFID[1])
		//run = false
		s := string(requestBody)
		data := HeartBeatData{}
		json.Unmarshal([]byte(s), &data)
		_, _ = w.Write([]byte(requestBody))
		block := DecodeFromJSON(string(data.BlockJSON))
		// verify if block exists

		// if block does not exist

		if verifyNonce(block) {
			result := Bc.GetBlock(block.Header.Height, block.Header.Hash)
			if result != nil {
				fmt.Println("in HBRecieve, insert succeses: ", block.Header.Hash)
				Bc.Insert(*block)
			} else {
				fmt.Println("in HBRecieve, need to ask for parent: ", block.Header.Hash)
				strheight := strconv.Itoa(int(block.Header.Height - 1))
				if askForParent(block.Header.ParentHash, strheight) {
					Bc.Insert(*block)
				}
			}
		}
		mutex.Unlock()
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// make function that sends out heartbeat data with a post request
// HeartBeatSend()

func readRequestBody(r *http.Request) (string, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.New("cannot read request body")
	}
	defer r.Body.Close()
	return string(reqBody), nil
}

// ShowHandler ...
func ShowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Bc.Show()))
}
