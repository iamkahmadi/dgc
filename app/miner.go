package app

import (
	"dgc/blockchain"
	"dgc/types"
	"dgc/wallet"
	"log"
)

// Miner struct represents a miner in the blockchain network
type Miner struct {
	Blockchain      *blockchain.Blockchain
	TransactionPool *wallet.TransactionPool
	Wallet          *types.Wallet
	P2PServer       *P2pServer
}

// NewMiner creates a new miner instance
func NewMiner(bc *blockchain.Blockchain, txPool *wallet.TransactionPool, w *types.Wallet, p2p *P2pServer) *Miner {
	return &Miner{
		Blockchain:      bc,
		TransactionPool: txPool,
		Wallet:          w,
		P2PServer:       p2p,
	}
}

// Mine mines a new block with valid transactions
func (m *Miner) Mine() types.Block {
	// Get valid transactions
	validTransactions := wallet.ValidTransactions(m.TransactionPool)

	// Add mining reward transaction
	rewardTx := wallet.RewardTransaction(m.Wallet, wallet.BlockchainWallet())
	validTransactions = append(validTransactions, rewardTx)

	// Mine a new block
	block := m.Blockchain.AddBlock(validTransactions)

	// Sync blockchain with peers
	m.P2PServer.syncChains()

	// Clear the transaction pool
	m.TransactionPool.Clear()

	// Broadcast to clear transactions
	m.P2PServer.BroadcastClearTransactions()

	log.Println("New block mined successfully:", block)

	return block
}
