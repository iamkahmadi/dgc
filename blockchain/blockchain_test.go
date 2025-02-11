package blockchain

import (
	"dgc/block"
	"fmt"
	"testing"
)

// TestNewBlockchain tests the creation of a new Blockchain with the Genesis block
func TestNewBlockchain(t *testing.T) {
	bc := NewBlockchain()

	// Check if the blockchain is initialized with only the Genesis block
	if len(bc.Chain) != 1 {
		t.Errorf("Expected blockchain length of 1, got %d", len(bc.Chain))
	}

	// Check if the first block is the Genesis block
	if fmt.Sprintf("%v", bc.Chain[0]) != fmt.Sprintf("%v", block.Genesis()) {
		t.Errorf("Expected Genesis block, got %v", bc.Chain[0])
	}
}

// TestAddBlock tests adding a new block to the blockchain
func TestAddBlock(t *testing.T) {
	bc := NewBlockchain()

	// Add a block to the blockchain
	data := []string{"block 1 data"}
	newBlock := bc.AddBlock(data)

	// Check if the blockchain length increased
	if len(bc.Chain) != 2 {
		t.Errorf("Expected blockchain length of 2, got %d", len(bc.Chain))
	}

	// Check if the newly added block is valid
	if newBlock.LastHash != bc.Chain[0].Hash {
		t.Errorf("Expected block's LastHash to be %s, got %s", bc.Chain[0].Hash, newBlock.LastHash)
	}
}

// TestIsValidChain tests if the IsValidChain function works correctly
func TestIsValidChain(t *testing.T) {
	bc := NewBlockchain()

	// Add a block to the blockchain
	data := []string{"block 1 data"}
	bc.AddBlock(data)

	// Validate the current chain
	if !bc.IsValidChain(bc.Chain) {
		t.Error("Blockchain should be valid")
	}

	// Modify the blockchain to invalidate it
	invalidChain := bc.Chain
	invalidChain[1].Data = []string{"tampered data"}

	// Validate the tampered chain
	if bc.IsValidChain(invalidChain) {
		t.Error("Blockchain should not be valid")
	}
}

func TestReplaceChain(t *testing.T) {
	bc := NewBlockchain()

	// Add the first block
	data1 := []string{"block 1 data"}
	bc.AddBlock(data1)

	// Create a new, longer chain with one additional block
	newChain := append([]block.Block{}, bc.Chain...)
	newBlock := block.MineBlock(bc.Chain[len(bc.Chain)-1], []string{"block 2 data"})
	newChain = append(newChain, newBlock)

	// Validate the new chain before replacing
	fmt.Println("New chain data before replacement:")
	for _, blk := range newChain {
		fmt.Println(blk.ToString())
	}

	if !bc.IsValidChain(newChain) {
		t.Error("New chain is not valid")
	}

	// Replace the current blockchain with the new chain
	bc.ReplaceChain(newChain)

	// Ensure the blockchain is replaced with the new chain
	if len(bc.Chain) != len(newChain) {
		t.Errorf("Expected blockchain length of %d, got %d", len(newChain), len(bc.Chain))
	}

	// Check if the new blockchain is valid
	if !bc.IsValidChain(bc.Chain) {
		t.Error("Blockchain should be valid after replacement")
	}
}

// TestReplaceInvalidChain tests replacing the chain with an invalid one
func TestReplaceInvalidChain(t *testing.T) {
	bc := NewBlockchain()

	// Add a block to the blockchain
	data := []string{"block 1 data"}
	bc.AddBlock(data)

	// Create a new chain with an invalid block (tampered data)
	newChain := append([]block.Block{}, bc.Chain...)
	newChain = append(newChain, block.MineBlock(bc.Chain[len(bc.Chain)-1], []string{"tampered block data"}))
	newChain[1].Data = []string{"tampered data"}

	// Attempt to replace the blockchain with the invalid chain
	bc.ReplaceChain(newChain)

	// Ensure the blockchain is not replaced with an invalid chain
	if len(bc.Chain) != 2 {
		t.Errorf("Blockchain length should still be 2, got %d", len(bc.Chain))
	}
}
