package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// StartTryingNonces ...
func (bc *BlockChain) StartTryingNonces() {
	stop := false

	for !stop {
		parentBlock := bc.GetLatestBlock()[0]
		//fmt.Println("Just created parent block ")
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

				} else {
					b.Header.Nonce = nonce
					fmt.Println("LOG: #blockgeneration: About to insert and nonce is found")
					Bc.Insert(b)
					SendHeartBeat(string(b.EncodeToJSON()))

					//stop = true // delete this line when running
				}
			}

		}
		//stop = true
	}
}

// SendHeartBeat ...
func SendHeartBeat(blockJSON string) {
	peerMapJSON, err := PEERLIST.PeerListToJSON()
	if err != nil {
		log.Fatalln(err)
	}
	HBData := NewHeartBeatData(SELFID, SELFADDRESS, blockJSON, string(peerMapJSON))
	HBDataJSON, _ := HBData.HBDataToJSON()
	for _, id := range PEERLIST.peerIDs {
		fmt.Println("LOG: Sending message to peer", string(id))
		resp, err := http.Post(string(id)+"/heartbeat/recieve", "application/json", bytes.NewBuffer(HBDataJSON))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("LOG: send heart beat - status ", resp.Status)
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
		fmt.Println("LOG: sending get request in askForParent", string(id))
		resp, err := http.Get(string(id) + "/block/" + height + "/" + parentHash)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("LOG: About to check if need another parent #askForParent")
		if resp.StatusCode != http.StatusNotFound {
			b := DecodeFromJSON(string(body))
			//fmt.Println("block string is ", string(body))
			fmt.Println("block json is", string(b.EncodeToJSON()))

			result := Bc.GetBlock(b.Header.Height, b.Header.Hash)
			if result != nil {
				fmt.Println("LOG: askForParent succeses: ", b.Header.Hash)
				Bc.Insert(*b)
				return true
			}
			fmt.Println("LOG: asking for another parent block #askForParent")
			strheight := strconv.Itoa(int(b.Header.Height - 1))
			if askForParent(b.Header.ParentHash, strheight) {
				Bc.Insert(*b)
				return true
			}
		}

	}
	fmt.Println("ask for parent does not work")
	return false
}

// DownloadChain goes to a node in peer list and asks for entire block
func DownloadChain() {
	for _, id := range PEERLIST.peerIDs {
		fmt.Println("LOG: #download Sending message to peer", string(id))
		resp, err := http.Get(string(id) + "/Upload")
		if err != nil {
			fmt.Println("asked for blockchain and there was an the err, ", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Log: in DonwloadChain status code is: " + string(resp.Status))
		if resp.StatusCode == http.StatusOK {
			fmt.Println("LOG: download chain status: ", resp.Status)
			fmt.Println("Copying block chain ")
			Bc.DecodeBlockchainFromJSON(string(body))
		}

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
