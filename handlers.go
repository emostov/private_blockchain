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

// Upload sends the entire json block chain
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Println("Log: succesful ask for block chain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(EncodeBlockchainToJSON(SYNCBC.BC)))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// AskForBlock reques
func AskForBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		h := vars["height"]
		hash := vars["hash"]
		height, err := strconv.ParseInt(h, 10, 64)
		fmt.Println("LOG: " + "ask get" + "height :" + h + ", hash: " + hash)
		if err == nil {
			block := SYNCBC.GetBlock(int32(height), hash)
			if block == nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(block.EncodeToJSON())
			}
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
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

		fmt.Println("LOG: HB Recieve: Im getting a post request ")
		//run = false
		s := string(requestBody)
		data := HeartBeatData{}
		json.Unmarshal([]byte(s), &data)
		_, _ = w.Write([]byte(requestBody))
		block := DecodeFromJSON(string(data.BlockJSON))

		if verifyNonce(block) {
			fmt.Println("LOG: HB Recieve: got post and nonce verified")
			result := SYNCBC.GetBlock(block.Header.Height, block.Header.Hash)
			resultParent := SYNCBC.GetBlock(block.Header.Height-1, block.Header.ParentHash)
			// verify if block exists already
			if result == nil && resultParent != nil {
				fmt.Println("LOG: HB Recieve: insert succeses: ", block.Header.Hash)
				SYNCBC.Insert(*block)
			} else if result == nil && resultParent == nil { // block does not exist and need parent
				fmt.Println("LOG: in HBRecieve, need to ask for parent: ", block.Header.Hash)
				strheight := strconv.Itoa(int(block.Header.Height - 1))
				if askForParent(block.Header.ParentHash, strheight) {
					SYNCBC.Insert(*block)
				}
			} else {
				fmt.Println("LOG: HB Recieve: Block and parent block already exist so no insert")
			}
		} else {
			fmt.Println("LOG: HB Recieve: got post and nonce NOT verified")
		}
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
	w.Write([]byte(SYNCBC.BC.Show()))
}

//Start simply starts a thread for mining. Make sure to only call once!
func Start(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println("LOG: I just got asked to start mining")
	DoMinerRegistration()
	DownloadChain()
	go SYNCBC.StartTryingNonces()
	w.Write([]byte("Mining Engaged"))
}

//Register returns registration information to node and updates SRD.PeerMap with new ID
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("LOG: This is the Register server and I just got a request to register")
		requestBody, err := readRequestBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
		regData := RegisterData{AssignedID: requestBody, PeerMapJSON: SRD.PeerMapJSON}
		w.Write(regData.EncodeRegisterDataToJSON())
		decodedID := DecodeIDFromJSON(requestBody)
		SRD.AddNewPeer(decodedID)
		SRD.EncodePeerMapToJSON()
	} else {
		fmt.Println("LOG: this is the Register server and I just got a BAD request")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
