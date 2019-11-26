package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
)

// ShowCanonicalOld displays canonical chain
func (blockchain *BlockChain) ShowCanonicalOld() string {
	rs := ""
	var idList []int
	for id := range blockchain.Chain {
		idList = append(idList, int(id))
	}

	sort.Sort(sort.Reverse(sort.IntSlice(idList)))
	for _, id := range idList {
		var hashs []string
		for _, block := range blockchain.Chain[int32(id)] {
			hashs = append(hashs, "height="+strconv.Itoa(int(block.Header.Height))+
				", timestamp="+strconv.Itoa(int(block.Header.Timestamp))+
				", hash="+block.Header.Hash+
				", parentHash="+block.Header.ParentHash+
				", size="+strconv.Itoa(int(block.Header.Size)))
		}
		sort.Strings(hashs)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	return rs
}

// ShowCanonical displays and returs string form o canonical chain
func (sbc *SyncBlockChain) ShowCanonical() string {
	rs := ""
	canonicalChains := sbc.getCanonicalChain()
	if len(canonicalChains) > 1 {
		rs += "Fork not resolved yet. There are multiple eligible chains." + "\n"
	} else {
		rs += "There is a canoincal chain." + "\n"
	}
	for i, bSlice := range canonicalChains {
		istr := strconv.Itoa(i)
		rs += "Chain: " + istr + "\n"
		for _, bStr := range bSlice {
			rs += fmt.Sprintf("%s, ", bStr)
			//rs += bStr + "\n"
		}
	}
	return rs
}

// GetCanonicalChain returns array representing canonical chain
func (sbc *SyncBlockChain) getCanonicalChain() [][]string {
	maxHeight := int32(-1)
	var maxBlocks []Block
	for _, block := range sbc.Get(sbc.BC.Length) {
		if block.Header.Difficulty > maxHeight {
			maxHeight = block.Header.Difficulty
			maxBlocks = []Block{block}
		} else if block.Header.Difficulty == maxHeight {
			maxBlocks = append(maxBlocks, block)
		}
	}
	returnString := [][]string{}
	for _, block := range maxBlocks {
		hashs := []string{}
		for block.Header.Height >= 1 {
			hashs = append(hashs, "height="+strconv.Itoa(int(block.Header.Height))+
				", timestamp="+strconv.Itoa(int(block.Header.Timestamp))+
				", hash="+block.Header.Hash+
				", parentHash="+block.Header.ParentHash+
				", size="+strconv.Itoa(int(block.Header.Size)))
			block = *(sbc.GetParentBlock(&block))
		}
		returnString = append(returnString, hashs)
	}
	return returnString
}

//Show displays the blockchain
func (blockchain *BlockChain) Show() string {
	rs := ""
	var idList []int
	for id := range blockchain.Chain {
		idList = append(idList, int(id))
	}
	sort.Ints(idList)
	for _, id := range idList {
		var hashs []string
		for _, block := range blockchain.Chain[int32(id)] {
			hashs = append(hashs, block.Header.Hash+"<="+block.Header.ParentHash)
		}
		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	sum := sha256.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is the BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs
	return rs
}
