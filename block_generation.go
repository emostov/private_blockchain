package main

import (
	"fmt"
	"log"
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
			if bc.Length == b.Header.Height {
				puzzleAnswer = makePuzzleAnswer(nonce, parentHash, blockValue)
				if checkPuzzleAnswerValid(target, puzzleAnswer) == false {
					nonce = generateNonce(nonce)
					run = true
					fmt.Println("Nonce not found")
				} else {
					b.Header.Nonce = nonce
					fmt.Println("block generation:About to insert after nonce is found")
					Bc.Insert(b)
					// delete this line when running
					stop = true
				}
			}

		}
		stop = true
	}
}

//Helper function for Start Trying Nonces

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
