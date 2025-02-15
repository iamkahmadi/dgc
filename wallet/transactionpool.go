package wallet

import (
	"dgc/types"
	"fmt"
)

// TransactionPool represents a pool of transactions
type TransactionPool struct {
	Transactions []*types.Transaction // Exported field
}

// NewTransactionPool creates a new instance of TransactionPool
func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		Transactions: make([]*types.Transaction, 0),
	}
}

// UpdateOrAddTransaction updates an existing transaction or adds a new one
func (tp *TransactionPool) UpdateOrAddTransaction(newTransaction *types.Transaction) {
	for i, existingTransaction := range tp.Transactions {
		if existingTransaction.ID == newTransaction.ID {
			tp.Transactions[i] = newTransaction // Update existing transaction
			return
		}
	}
	tp.Transactions = append(tp.Transactions, newTransaction) // Add new transaction
}

// ExistingTransaction returns an existing transaction for a given address
func ExistingTransaction(tp *TransactionPool, address string) *types.Transaction {
	for _, t := range tp.Transactions {
		if t.Input.Address == address {
			return t // Return the found transaction
		}
	}
	return nil // No existing transaction found
}

// ValidTransactions returns a slice of valid transactions
func ValidTransactions(tp *TransactionPool) []*types.Transaction {
	validTxns := make([]*types.Transaction, 0)

	for _, t := range tp.Transactions {
		outputTotal := 0
		for _, output := range t.Outputs {
			outputTotal += output.Amount
		}

		if t.Input.Amount != outputTotal {
			fmt.Printf("Invalid transaction from %s.\n", t.Input.Address)
			continue
		}

		if !VerifyTransaction(t) {
			fmt.Printf("Invalid signature from %s.\n", t.Input.Address)
			continue
		}

		validTxns = append(validTxns, t) // Add valid transaction to the slice
	}

	return validTxns // Return all valid transactions
}

// Clear removes all transactions from the pool
func (tp *TransactionPool) Clear() {
	tp.Transactions = []*types.Transaction{} // Reset the Transactions slice
}
