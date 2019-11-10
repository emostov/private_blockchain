package main

// GetLatestBlock returns slice of blocks at chains length
func (bc *BlockChain) GetLatestBlock() []Block {
	mutex.Lock()
	ret := bc.Get(bc.Length)
	mutex.Unlock()
	return ret
}

// GetParentBlock takes a block as a parameter, and returns its parent block
func (bc *BlockChain) GetParentBlock(b *Block) *Block {
	mutex.Lock()
	parentHeightBlocks := bc.Get(b.Header.Height)
	mutex.Unlock()
	for _, pBlock := range parentHeightBlocks {
		if pBlock.Header.Hash == b.Header.ParentHash {
			return &pBlock
		}
	}
	return nil
}
