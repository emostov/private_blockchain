package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(block.EncodeToJSON())
		}
	}
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
		fmt.Println("Im an id and am getting a post", SELFID[1])
		//run = false
		s := string(requestBody)
		data := HeartBeatData{}
		json.Unmarshal([]byte(s), &data)
		_, _ = w.Write([]byte(requestBody))
		block := DecodeFromJSON(string(data.BlockJSON))

		// if block does not exist

		if verifyNonce(block) {
			result := Bc.GetBlock(block.Header.Height, block.Header.Hash)
			resultParent := Bc.GetBlock(block.Header.Height-1, block.Header.ParentHash)
			// verify if block exists already
			if result == nil && resultParent != nil {
				fmt.Println("in HBRecieve, insert succeses: ", block.Header.Hash)
				Bc.Insert(*block)
			} else if result == nil && resultParent == nil { // block does not exist and need parent
				fmt.Println("in HBRecieve, need to ask for parent: ", block.Header.Hash)
				strheight := strconv.Itoa(int(block.Header.Height - 1))
				if askForParent(block.Header.ParentHash, strheight) {
					Bc.Insert(*block)
				}
			} else {
				fmt.Println("HB Recieve: Block and parent block already exist so no insert")
			}
		}
		mutex.Unlock()
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

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

//Start simply starts a thread for mining. Make sure to only call once!
func Start(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println("I am at port " + SELFID[1] + ", and I just got asked to start mining")
	go Bc.StartTryingNonces()
}
