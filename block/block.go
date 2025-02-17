package block

import (
	"dgc/types"
	"dgc/util"
	"fmt"
	"strings"
	"time"
)

// ToString returns the string representation of a block
func ToString(b *types.Block) string {
	return fmt.Sprintf("Block - Timestamp : %d LastHash : %s Hash : %s Nonce : %d Difficulty: %d Data : %v", b.Timestamp, b.LastHash, b.Hash, b.Nonce, b.Difficulty, b.Data)
}

// Genesis creates the genesis block
func Genesis() types.Block {
	// time.Now().Unix() for current time
	return types.Block{
		Timestamp:  63,
		LastHash:   "-----",
		Hash:       "GenesisHash",
		Data:       []*types.Transaction{}, // Initialize with empty slice of Transactions
		Nonce:      0,
		Difficulty: util.DIFFICULTY,
	}
}

// MineBlock mines a new block based on the previous block
func MineBlock(lastBlock types.Block, data []*types.Transaction) types.Block { // Change data parameter type
	var hash string
	var timestamp int64
	lastHash := lastBlock.Hash
	difficulty := lastBlock.Difficulty
	nonce := 0

	for {
		nonce++
		timestamp = time.Now().Unix()
		difficulty = AdjustDifficulty(lastBlock, timestamp)
		hash = Hash(timestamp, lastHash, data, nonce, difficulty)

		if nonce%1000 == 0 {
			fmt.Println("Nonce Generated up to: ", nonce)
		}

		if strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
			break
		}
	}

	return types.Block{
		Timestamp:  timestamp,
		LastHash:   lastHash,
		Hash:       hash,
		Data:       data,
		Nonce:      nonce,
		Difficulty: difficulty,
	}
}

// // Hash generates the hash of a block
// func Hash(timestamp int64, lastHash string, data []*types.Transaction, nonce int, difficulty int) string {
// 	return util.ChainUtilHash(fmt.Sprintf("%d%s%v%d%d", timestamp, lastHash, data, nonce, difficulty))
// }

func Hash(timestamp int64, lastHash string, data []*types.Transaction, nonce int, difficulty int) string {
	// Serialize the data explicitly by converting each transaction to a string
	var serializedData string
	for _, tx := range data {
		serializedData += fmt.Sprintf("%s%s%d%v", tx.ID, tx.Input.Address, tx.Input.Amount, tx.Outputs)
	}
	return util.ChainUtilHash(fmt.Sprintf("%d%s%s%d%d", timestamp, lastHash, serializedData, nonce, difficulty))
}

// BlockHash creates the block hash using SHA256 and is a method of the Block type (converted to function)
func BlockHash(b *types.Block) string {
	return Hash(b.Timestamp, b.LastHash, b.Data, b.Nonce, b.Difficulty)
}

// AdjustDifficulty adjusts the difficulty based on the time taken to mine the block
func AdjustDifficulty(lastBlock types.Block, currentTime int64) int {
	difficulty := lastBlock.Difficulty
	if lastBlock.Timestamp+util.MINE_RATE > currentTime {
		difficulty++
	} else {
		difficulty--
	}
	return difficulty
}
