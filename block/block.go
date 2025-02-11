package block

import (
	"dgc/util"
	"fmt"
	"strings"
	"time"
)

// Block struct represents a block in the blockchain
type Block struct {
	Timestamp  int64
	LastHash   string
	Hash       string
	Data       []string
	Nonce      int
	Difficulty int
}

// ToString returns the string representation of a block
func (b *Block) ToString() string {
	return fmt.Sprintf("Block - Timestamp : %d LastHash : %s Hash : %s Nonce : %d Difficulty: %d Data : %v", b.Timestamp, b.LastHash, b.Hash, b.Nonce, b.Difficulty, b.Data)
}

// Genesis creates the genesis block
func Genesis() Block {
	// time.Now().Unix() for current time
	return Block{
		Timestamp:  63,
		LastHash:   "-----",
		Hash:       "GenesisHash",
		Data:       []string{},
		Nonce:      0,
		Difficulty: util.DIFFICULTY,
	}
}

// MineBlock mines a new block based on the previous block
func MineBlock(lastBlock Block, data []string) Block {
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
			fmt.Println("Nonce Generated upto: ", nonce)
		}

		if strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
			break
		}
	}

	return Block{
		Timestamp:  timestamp,
		LastHash:   lastHash,
		Hash:       hash,
		Data:       data,
		Nonce:      nonce,
		Difficulty: difficulty,
	}
}

// Hash generates the hash of a block
func Hash(timestamp int64, lastHash string, data []string, nonce int, difficulty int) string {
	return util.ChainUtilHash(fmt.Sprintf("%d%s%v%d%d", timestamp, lastHash, data, nonce, difficulty))
}

// BlockHash creates the block hash using SHA256 and is a method of the Block type
func (b *Block) BlockHash() string {
	return Hash(b.Timestamp, b.LastHash, b.Data, b.Nonce, b.Difficulty)
}

// AdjustDifficulty adjusts the difficulty based on the time taken to mine the block
func AdjustDifficulty(lastBlock Block, currentTime int64) int {
	difficulty := lastBlock.Difficulty
	if lastBlock.Timestamp+util.MINE_RATE > currentTime {
		difficulty++
	} else {
		difficulty--
	}
	return difficulty
}
