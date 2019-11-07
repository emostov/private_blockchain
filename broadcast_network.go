package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
)

// ShowCanonical displays canonical chain
func (blockchain *BlockChain) ShowCanonical() string {
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
