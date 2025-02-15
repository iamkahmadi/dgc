package wallet

import (
	"fmt"
	"log"

	"dgc/blockchain"
	"dgc/types" // Import the new types package
	"dgc/util"
)

// NewWallet creates a new wallet with an initial balance and key pair
func NewWallet() *types.Wallet {
	privKey, pubKey, err := util.GenKeyPair()
	if err != nil {
		log.Fatalf("Error generating key pair: %v", err)
	}
	return &types.Wallet{
		Balance:   util.INITIAL_BALANCE,
		KeyPair:   privKey,
		PublicKey: util.PublicKeyToHex(pubKey),
	}
}

// String returns a string representation of the wallet
func String(w *types.Wallet) string {
	return fmt.Sprintf("Wallet -\nPublic Key: %s\nBalance: %d", w.PublicKey, w.Balance)
}

// Sign signs data using the wallet's private key
func Sign(w *types.Wallet, dataHash string) []byte {
	signature, err := util.SignData(w.KeyPair, dataHash)
	if err != nil {
		log.Fatalf("Error signing data: %v", err)
	}
	return signature
}

// CreateTransaction creates a new transaction or updates an existing one
func CreateTransaction(wallet *types.Wallet, recipient string, amount int, blockchain *blockchain.Blockchain, transactionPool *TransactionPool) *types.Transaction {
	wallet.Balance = CalculateBalance(wallet, blockchain)

	if amount > wallet.Balance {
		fmt.Printf("Amount: %d exceeds current balance: %d\n", amount, wallet.Balance)
		return nil
	}

	existingTransaction := ExistingTransaction(transactionPool, wallet.PublicKey)

	if existingTransaction != nil {
		UpdateTransaction(existingTransaction, wallet, recipient, amount)
		return existingTransaction
	}

	newTransaction := NewTransaction(wallet, recipient, amount)
	transactionPool.UpdateOrAddTransaction(newTransaction)

	return newTransaction
}

// CalculateBalance calculates the current balance of the wallet based on the blockchain
func CalculateBalance(wallet *types.Wallet, blockchain *blockchain.Blockchain) int {
	balance := wallet.Balance
	var transactions []*types.Transaction

	for _, block := range blockchain.Chain {
		for _, tx := range block.Data {
			transactions = append(transactions, tx)
		}
	}

	walletInputTs := filterTransactions(transactions, func(tx *types.Transaction) bool {
		return tx.Input.Address == wallet.PublicKey
	})

	var startTime int64 = 0

	if len(walletInputTs) > 0 {
		recentInputTx := walletInputTs[0]
		for _, tx := range walletInputTs {
			if tx.Input.Timestamp > recentInputTx.Input.Timestamp {
				recentInputTx = tx
			}
		}

		output := FindOutput(recentInputTx, wallet.PublicKey)
		if output != nil {
			balance = output.Amount
		}

		startTime = recentInputTx.Input.Timestamp
	}

	for _, tx := range transactions {
		if tx.Input.Timestamp > startTime {
			for _, output := range tx.Outputs {
				if output.Address == wallet.PublicKey {
					balance += output.Amount
				}
			}
		}
	}

	return balance
}

// BlockchainWallet creates a special wallet for the blockchain itself
func BlockchainWallet() *types.Wallet {
	blockchainWallet := NewWallet()
	blockchainWallet.PublicKey = "blockchain-wallet"
	return blockchainWallet
}

// filterTransactions filters transactions based on a predicate function
func filterTransactions(transactions []*types.Transaction, predicate func(*types.Transaction) bool) []*types.Transaction {
	var result []*types.Transaction
	for _, tx := range transactions {
		if predicate(tx) {
			result = append(result, tx)
		}
	}
	return result
}
