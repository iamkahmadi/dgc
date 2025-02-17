package blockchain

import (
	"dgc/block"
	"dgc/types"
	"fmt"
)

type Blockchain struct {
	Chain []types.Block
}

// NewBlockchain initializes a new Blockchain with the Genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain: []types.Block{block.Genesis()},
	}
}

// AddBlock adds a new block to the chain
func (bc *Blockchain) AddBlock(data []*types.Transaction) types.Block {
	lastBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := block.MineBlock(lastBlock, data)
	bc.Chain = append(bc.Chain, newBlock)
	return newBlock
}

// IsValidChain validates if the provided chain is a valid blockchain
func (bc *Blockchain) IsValidChain(chain []types.Block) bool {
	// Check if the first block is the Genesis block
	if fmt.Sprintf("%v", chain[0]) != fmt.Sprintf("%v", block.Genesis()) {
		fmt.Println("Genesis block is invalid.")
		fmt.Println("Expected Genesis Block Hash:", block.Genesis().Hash)
		fmt.Println("Actual Genesis Block Hash:", chain[0].Hash)
		return false
	}

	// Check the rest of the blocks
	for i := 1; i < len(chain); i++ {
		currentBlock := chain[i]
		lastBlock := chain[i-1]

		// Check that the current block's hash matches the expected one
		if currentBlock.LastHash != lastBlock.Hash || currentBlock.Hash != block.BlockHash(&currentBlock) {
			fmt.Printf("Block %d is invalid: LastHash: %s, Expected LastHash: %s\n", i, currentBlock.LastHash, lastBlock.Hash)
			return false
		}
	}

	return true
}

// ReplaceChain replaces the current blockchain with a new one if it's valid and longer
func (bc *Blockchain) ReplaceChain(newChain []types.Block) {
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
