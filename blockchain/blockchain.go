package blockchain

import (
	"dgc/block"
	"fmt"
)

type Blockchain struct {
	Chain []block.Block
}

// NewBlockchain initializes a new Blockchain with the Genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain: []block.Block{block.Genesis()},
	}
}

// AddBlock adds a new block to the chain
func (bc *Blockchain) AddBlock(data []string) block.Block {
	lastBlock := bc.Chain[len(bc.Chain)-1]
	block := block.MineBlock(lastBlock, data)
	bc.Chain = append(bc.Chain, block)
	return block
}

// IsValidChain validates if the provided chain is a valid blockchain
func (bc *Blockchain) IsValidChain(chain []block.Block) bool {
	// Check if the first block is the Genesis block
	if fmt.Sprintf("%v", chain[0]) != fmt.Sprintf("%v", block.Genesis()) {
		fmt.Println("Genesis block is invalid.")
		fmt.Println("Expected Genesis Block Hash:", block.Genesis().Hash)
		fmt.Println("Actual Genesis Block Hash:", chain[0].Hash)
		return false
	}

	// Check the rest of the blocks
	for i := 1; i < len(chain); i++ {
		block := chain[i]
		lastBlock := chain[i-1]

		// Check that the current block's hash matches the expected one
		if block.LastHash != lastBlock.Hash || block.Hash != block.BlockHash() {
			fmt.Printf("Block %d is invalid: LastHash: %s, Expected LastHash: %s\n", i, block.LastHash, lastBlock.Hash)
			return false
		}
	}

	return true
}

// ReplaceChain replaces the current blockchain with a new one if it's valid and longer
func (bc *Blockchain) ReplaceChain(newChain []block.Block) {
	// Ensure the new chain is longer than the current chain
	if len(newChain) <= len(bc.Chain) {
		fmt.Println("Received chain is not longer than the current chain.")
		return
	}

	// Ensure the new chain is valid
	if !bc.IsValidChain(newChain) {
		fmt.Println("The received chain is not valid.")
		return
	}

	// Replace the chain
	fmt.Println("Replacing blockchain with the new chain.")
	bc.Chain = newChain
}
