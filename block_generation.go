package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// StartTryingNonces attempts nonces. It mines at a new height once the bc.length
// changes. When a block is found it sends out heartbeat data. Uses a random
// seed value as the initial nonce and then increments. Incrementing does not
// wrap around
func (sbc *SyncBlockChain) StartTryingNonces() {
	log.Println("LOG: StartTryingNonces: started mining ")
	stop := false
	for !stop {
		parentBlock := sbc.GetLatestBlock()[0]
		parentHash := parentBlock.Header.Hash
		var b Block
		b.Initialize(sbc.BC.Length+1, parentHash, "test block value", int32(0))
		blockValue := b.Value
		rand.Seed(time.Now().UnixNano()) // generate a random number for start nonce
		randStart := rand.Intn(300)
		nonce := generateStartNonce(randStart)
		puzzleAnswer := ""
		run := true
		for run {
			run = false
			if sbc.BC.Length+1 == b.Header.Height {

				puzzleAnswer = makePuzzleAnswer(nonce, parentHash, blockValue)
				if checkPuzzleAnswerValid(TARGET, puzzleAnswer) == false {
					nonce = generateNonce(nonce)
					run = true

				} else {
					b.Header.Nonce = nonce
					log.Println("LOG: StartTryingNonces: About to insert and nonce is found")
					SYNCBC.Insert(b)
					SendHeartBeat(string(b.EncodeToJSON()))
				}
			}
		}
	}
}

// SendHeartBeat to all peers in peerlist
func SendHeartBeat(blockJSON string) {
	peerMapJSON := EncodeIDListToJSON(PEERLIST.PeerIDs)
	HBData := NewHeartBeatData(MINERID.Port, MINERID.Address, blockJSON, string(peerMapJSON))
	HBDataJSON, _ := HBData.HBDataToJSON()
	if len(PEERLIST.PeerIDs) >= 1 {
		for _, id := range PEERLIST.PeerIDs {
			log.Println("LOG: Sending message to peer", string(id.Port))
			resp, err := http.Post(id.Address+id.Port+"/heartbeat/recieve", "application/json", bytes.NewBuffer(HBDataJSON))
			if err != nil {
				log.Print(err)
				return
			}
			log.Println("LOG: send heart beat - status ", resp.Status)
		}
	} else {
		log.Println("LOG: SendHeartBeat: No peers to send heartbeat to")
	}
}

// askForParent is a recursive function that asks for the parent block(s) of a block
// It asks peers in peerlist until a parent block is found.
func askForParent(parentHash string, height string) bool {
	intheight, err := strconv.ParseInt((height), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	if intheight < 0 {
		return false
	}
	for _, id := range PEERLIST.PeerIDs {
		log.Println("LOG: sending get request in askForParent", string(id.Port))
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(id.Address + id.Port + "/block/" + height + "/" + parentHash)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("LOG: About to check if need another parent #askForParent")
		if resp.StatusCode != http.StatusNotFound {
			b := DecodeFromJSON(string(body))
			log.Println("block json is", string(b.EncodeToJSON()))
			result := SYNCBC.GetBlock(b.Header.Height, b.Header.Hash)
			if result != nil || b.Header.Height == 0 {
				log.Println("LOG: askForParent succeses: ", b.Header.Hash)
				SYNCBC.Insert(*b)
				return true
			}
			log.Println("LOG: asking for another parent block #askForParent")
			strheight := strconv.Itoa(int(b.Header.Height - 1))
			if askForParent(b.Header.ParentHash, strheight) {
				log.Println("LOG: askForParent succeses: ", b.Header.Hash)
				SYNCBC.Insert(*b)
				return true
			}
		}
	}
	log.Println("ask for parent does not mpy succesful")
	return false
}

// DownloadChain goes to server node and asks for entire blockchain
func DownloadChain() {
	log.Println("LOG: DownloadChain: asking peer ", string(SID.Port), " for chain")
	resp, err := http.Get(SID.Address + SID.Port + "/Upload")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Log: in DonwloadChain status code is: " + string(resp.Status))
	if resp.StatusCode == http.StatusOK {
		SYNCBC.DecodeBlockchainFromJSON(string(body))
		return
	}
}

// takes int string, converts to int, adds 1 and converts back to int string
func generateNonce(previous string) string {
	previousInt, err := strconv.Atoi(previous)
	if err != nil {
		log.Fatal("error in generateNonce", err)
	}
	newNonce := strconv.Itoa(previousInt + 10)
	return newNonce
}

func verifyNonce(b *Block) bool {
	puzzleAnswer :=
		makePuzzleAnswer(b.Header.Nonce, b.Header.ParentHash, b.Value)
	return checkPuzzleAnswerValid(TARGET, puzzleAnswer)
}

func generateStartNonce(seed int) string {
	str := strconv.Itoa(seed)
	return str
}

func makePuzzleAnswer(nonce string, parentBlockHash string, blockValue string) string {
	puzzleString := parentBlockHash + nonce + blockValue
	puzzleAnswer := makeSha256Digest(puzzleString)
	return puzzleAnswer
}

func checkPuzzleAnswerValid(target string, puzzleAnswer string) bool {
	return target == puzzleAnswer[:len(target)]
}
