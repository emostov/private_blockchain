package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Upload sends the entire json block chain
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Println("Log: succesful ask for block chain")
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
		log.Println("LOG: " + "ask get" + "height :" + h + ", hash: " + hash)
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

		log.Println("LOG: HB Recieve: Im getting a post request ")
		//run = false
		s := string(requestBody)
		data := HeartBeatData{}
		json.Unmarshal([]byte(s), &data)
		_, _ = w.Write([]byte(requestBody))
		log.Println("LOG: HB Recieve: address, port ", data.Address, data.ID)

		// below would implement gossip protocol
		// peermap := DecodePeerMapFromJSON(data.PeerMapJSON)
		// TODO double check data type for HeartBeat Peer Map
		otherID := ID{Port: data.ID, Address: data.Address}
		// peermap = append(peermap, otherID)
		PEERLIST.AddNewPeers([]ID{otherID})
		PEERLIST.RemoveSelfFromPeerList()

		block := DecodeFromJSON(string(data.BlockJSON))
		if verifyNonce(block) {
			log.Println("LOG: HB Recieve: got post and nonce verified")
			result := SYNCBC.GetBlock(block.Header.Height, block.Header.Hash)
			resultParent := SYNCBC.GetBlock(block.Header.Height-1, block.Header.ParentHash)
			// verify if block exists already
			if result == nil && resultParent != nil {
				log.Println("LOG: HB Recieve: insert succeses: ", block.Header.Hash)
				SYNCBC.Insert(*block)
			} else if result == nil && resultParent == nil { // block does not exist and need parent
				log.Println("LOG: in HBRecieve, need to ask for parent: ", block.Header.Hash)
				strheight := strconv.Itoa(int(block.Header.Height - 1))
				if askForParent(block.Header.ParentHash, strheight) {
					SYNCBC.Insert(*block)
				}
			} else {
				log.Println("LOG: HB Recieve: Block and parent block already exist so no insert")
			}
		} else {
			log.Println("LOG: HB Recieve: got post and nonce NOT verified")
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
	w.Write([]byte(SYNCBC.BC.ShowCanonical()))
}

//Start simply starts a thread for mining. Make sure to only call once!
func Start(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println("LOG: Start: I just got asked to start mining")
	if MINERID.Port != "6688" {
		DoMinerRegistration()
		DownloadChain()
	}
	go SYNCBC.StartTryingNonces()

	w.Write([]byte("LOG: Start: Mining Engaged"))

}

// Starttest just for doing test ask right now
func Starttest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println("LOG: I am server, I just got asked to start mining")
	DoMinerRegistration()
	// DownloadChain()
	// go SYNCBC.StartTryingNonces()
	w.Write([]byte("Server Mining Engaged"))
	testAsk()
}

//Register returns registration information to node and updates SRD.PeerMap with new ID
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("LOG: This is the Register server and I just got a request to register")
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
		log.Println("LOG: this is the Register server and I just got a BAD request")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
