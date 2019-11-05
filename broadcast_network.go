package main

import (
	"fmt"
	"sort"
	"strconv"
)

func (blockchain *Blockchain) ShowCanonical() string {
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
