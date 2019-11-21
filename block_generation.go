package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// StartTryingNonces ...
func (sbc *SyncBlockChain) StartTryingNonces() {
	fmt.Println("LOG: StartTryingNonces: started mining ")
	stop := false

	for !stop {
		parentBlock := sbc.GetLatestBlock()[0]
		// fmt.Println("Just created parent block ")
		parentHash := parentBlock.Header.Hash
		var b Block
		b.Initialize(sbc.BC.Length+1, parentHash, "test block value")
		blockValue := b.Value

		nonce := generateStartNonce(1)
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
					fmt.Println("LOG: StartTryingNonces: About to insert and nonce is found")
					SYNCBC.Insert(b)
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
	//peerMapJSON, err := PEERLIST.PeerListToJSON()
	peerMapJSON := EncodeIDListToJSON(PEERLIST.PeerIDs)

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	HBData := NewHeartBeatData(MINERID.Port, MINERID.Address, blockJSON, string(peerMapJSON))
	HBDataJSON, _ := HBData.HBDataToJSON()
	if len(PEERLIST.PeerIDs) >= 1 {
		for _, id := range PEERLIST.PeerIDs {
			fmt.Println("LOG: Sending message to peer", string(id.Port))
			resp, err := http.Post(id.Address+id.Port+"/heartbeat/recieve", "application/json", bytes.NewBuffer(HBDataJSON))
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Println("LOG: send heart beat - status ", resp.Status)
		}
	} else {
		fmt.Println("LOG: SendHeartBeat: No peers to send heartbeat to")
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
	for _, id := range PEERLIST.PeerIDs {
		fmt.Println("LOG: sending get request in askForParent", string(id.Port))
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(id.Address + id.Port + "/block/" + height + "/" + parentHash)
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

			result := SYNCBC.GetBlock(b.Header.Height, b.Header.Hash)
			if result != nil {
				fmt.Println("LOG: askForParent succeses: ", b.Header.Hash)
				SYNCBC.Insert(*b)
				return true
			}
			fmt.Println("LOG: asking for another parent block #askForParent")
			strheight := strconv.Itoa(int(b.Header.Height - 1))
			if askForParent(b.Header.ParentHash, strheight) {
				SYNCBC.Insert(*b)
				return true
			}
		}

	}
	fmt.Println("ask for parent does not work")
	return false
}

// DownloadChain goes to a node in peer list and asks for entire block
func DownloadChain() {
	// if len(PEERLIST.PeerIDs) >= 1 {
	// for _, id := range PEERLIST.PeerIDs {
	fmt.Println("LOG: DownloadChain: asking peer ", string(SID.Port), " for chain")
	resp, err := http.Get(SID.Address + SID.Port + "/Upload")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Log: in DonwloadChain status code is: " + string(resp.Status))
	if resp.StatusCode == http.StatusOK {
		SYNCBC.DecodeBlockchainFromJSON(string(body))
		return
	}

	// }
	// } else {
	// 	fmt.Println("LOG: DownloadChain: no peers to download from")
	// }

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
	// fmt.Println("target: ", target, "puzzleAnswer ", puzzleAnswer)
	return target == puzzleAnswer[:len(target)]
}
