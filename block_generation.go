package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var target = "00"

// StartTryingNonces ...
func (bc *BlockChain) StartTryingNonces() {
	stop := false

	for !stop {
		parentBlock := bc.GetLatestBlock()[0]
		fmt.Println("Just created parent block ")
		parentHash := parentBlock.Header.Hash
		var b Block
		b.Initialize(bc.Length+1, parentHash, "test block value")
		blockValue := b.Value

		nonce := generateStartNonce(1)
		puzzleAnswer := ""
		run := true
		for run {
			run = false
			if bc.Length+1 == b.Header.Height {
				puzzleAnswer = makePuzzleAnswer(nonce, parentHash, blockValue)
				if checkPuzzleAnswerValid(target, puzzleAnswer) == false {
					nonce = generateNonce(nonce)
					run = true
					// fmt.Println("Nonce not found")
				} else {
					//mutex.Lock()
					b.Header.Nonce = nonce
					fmt.Println("block generation:About to insert after nonce is found")
					Bc.Insert(b)
					SendHeartBeat(string(b.EncodeToJSON()))
					// delete this line when running
					//stop = true
					//mutex.Unlock()
				}
			}

		}
		stop = true
	}
}

//Helper function for Start Trying Nonces

// SendHeartBeat ...
func SendHeartBeat(blockJSON string) {
	peerMapJSON, err := PEERLIST.PeerListToJSON()
	if err != nil {
		log.Fatalln(err)
	}
	HBData := NewHeartBeatData(SELFID, SELFADDRESS, blockJSON, string(peerMapJSON))
	HBDataJSON, _ := HBData.HBDataToJSON()
	for _, id := range PEERLIST.peerIDs {
		fmt.Println("Sending message to peer", string(id))
		resp, err := http.Post(string(id), "application/json", bytes.NewBuffer(HBDataJSON))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(resp.Status)
	}
}

func askForParent(parentHash string, height string) bool {
	intheight, err := strconv.ParseInt((height), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	if intheight < 0 {
		return false
	}
	for _, id := range PEERLIST.peerIDs {
		fmt.Println("sending get reques", string(id))
		resp, err := http.Get(string(id) + "/block/" + height + "/" + parentHash)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusNotFound {
			b := DecodeFromJSON(string(body))
			result := Bc.GetBlock(b.Header.Height, b.Header.Hash)
			if result != nil {
				Bc.Insert(*b)
				return true
			}
			strheight := strconv.Itoa(int(b.Header.Height))
			if askForParent(b.Header.ParentHash, strheight) {
				Bc.Insert(*b)
			}
		}

	}
	return false
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
	return checkPuzzleAnswerValid(target, puzzleAnswer)
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
	// fmt.Println("target: ", target, "puzzleAnswer ", puzzleAnswer)
	return target == puzzleAnswer[:len(target)]
}
