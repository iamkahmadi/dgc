package types

import (
	"crypto/ecdsa"
)

// Transaction represents a cryptocurrency transaction
type Transaction struct {
	ID      string
	Input   *Input
	Outputs []Output
}

// Input represents the input of a transaction
type Input struct {
	Timestamp int64
	Amount    int
	Address   string
	Signature []byte
}

// Output represents an output of a transaction
type Output struct {
	Amount  int
	Address string
}

// Wallet represents a cryptocurrency wallet
type Wallet struct {
	Balance   int
	KeyPair   *ecdsa.PrivateKey
	PublicKey string
}

// Block struct represents a block in the blockchain
type Block struct {
	Timestamp  int64
	LastHash   string
	Hash       string
	Data       []*Transaction
	Nonce      int
	Difficulty int
}
