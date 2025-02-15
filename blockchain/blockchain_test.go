package blockchain

import (
	"dgc/block"
	"dgc/types"
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

	// Create a transaction for testing
	transaction := &types.Transaction{
		ID: "1",
		Input: &types.Input{
			Timestamp: 0,
			Amount:    10,
			Address:   "test_address",
			Signature: []byte("signature"),
		},
		Outputs: []types.Output{
			{Amount: 10, Address: "recipient_address"},
		},
	}

	data := []*types.Transaction{transaction} // Use []*types.Transaction instead of []string
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

	// Create a transaction for testing
	transaction := &types.Transaction{
		ID: "1",
		Input: &types.Input{
			Timestamp: 0,
			Amount:    10,
			Address:   "test_address",
			Signature: []byte("signature"),
		},
		Outputs: []types.Output{
			{Amount: 10, Address: "recipient_address"},
		},
	}

	data := []*types.Transaction{transaction} // Use []*types.Transaction instead of []string
	bc.AddBlock(data)

	// Validate the current chain
	if !bc.IsValidChain(bc.Chain) {
		t.Error("Blockchain should be valid")
	}

	// Modify the blockchain to invalidate it
	bc.Chain[1].Data[0].Outputs[0].Address = "tampered_address" // Tamper with a transaction

	if bc.IsValidChain(bc.Chain) {
		t.Error("Blockchain should not be valid after tampering")
	}
}

// TestReplaceChain tests replacing the current chain with a new valid chain
func TestReplaceChain(t *testing.T) {
	bc := NewBlockchain()

	transaction1 := &types.Transaction{
		ID: "1",
		Input: &types.Input{
			Timestamp: 0,
			Amount:    10,
			Address:   "test_address",
			Signature: []byte("signature"),
		},
		Outputs: []types.Output{
			{Amount: 10, Address: "recipient_address"},
		},
	}

	transaction2 := &types.Transaction{
		ID: "2",
		Input: &types.Input{
			Timestamp: 0,
			Amount:    20,
			Address:   "test_address_2",
			Signature: []byte("signature_2"),
		},
		Outputs: []types.Output{
			{Amount: 20, Address: "recipient_address_2"},
		},
	}

	bc.AddBlock([]*types.Transaction{transaction1})

	newChain := append([]types.Block{}, bc.Chain...)
	newBlock := block.MineBlock(bc.Chain[len(bc.Chain)-1], []*types.Transaction{transaction2})
	newChain = append(newChain, newBlock)

	if !bc.IsValidChain(newChain) {
		t.Error("New chain is not valid")
	}

	bc.ReplaceChain(newChain)

	if len(bc.Chain) != len(newChain) {
		t.Errorf("Expected blockchain length of %d, got %d", len(newChain), len(bc.Chain))
	}

	if !bc.IsValidChain(bc.Chain) {
		t.Error("Blockchain should be valid after replacement")
	}
}

// TestReplaceInvalidChain tests replacing the chain with an invalid one
func TestReplaceInvalidChain(t *testing.T) {
	bc := NewBlockchain()

	transaction := &types.Transaction{
		ID: "1",
		Input: &types.Input{
			Timestamp: 0,
			Amount:    10,
			Address:   "test_address",
			Signature: []byte("signature"),
		},
		Outputs: []types.Output{
			{Amount: 10, Address: "recipient_address"},
		},
	}

	bc.AddBlock([]*types.Transaction{transaction})

	newChain := append([]types.Block{}, bc.Chain...)
	newBlock := block.MineBlock(bc.Chain[len(bc.Chain)-1], []*types.Transaction{transaction})
	newBlock.Data[0].Outputs[0].Address = "tampered_address" // Tamper with a transaction

	newChain = append(newChain, newBlock)

	bc.ReplaceChain(newChain)

	if len(bc.Chain) != 2 {
		t.Errorf("Blockchain length should still be 2, got %d", len(bc.Chain))
	}
}
