package main

import (
	"dgc/blockchain"
	"fmt"
)

func main() {
	// Create a new Blockchain
	bc := blockchain.NewBlockchain()

	// Display the Genesis Block
	fmt.Println("Genesis Block: ")
	fmt.Println(bc.Chain[0].ToString())

	// Add a new block to the blockchain 1
	data1 := []string{"Some transaction data 1"}
	newBlock1 := bc.AddBlock(data1)
	fmt.Println("\nNewly mined block 1: ")
	fmt.Println(newBlock1.ToString())

	// Add a new block to the blockchain
	data2 := []string{"Some transaction data 2"}
	newBlock2 := bc.AddBlock(data2)
	fmt.Println("\nNewly mined block 2: ")
	fmt.Println(newBlock2.ToString())

	// Display the current blockchain
	fmt.Println("\nCurrent Blockchain: ")
	for _, blk := range bc.Chain {
		fmt.Println(blk.ToString())
	}

	// Validate the current chain
	if bc.IsValidChain(bc.Chain) {
		fmt.Println("\nBlockchain is valid.")
	} else {
		fmt.Println("\nBlockchain is invalid.")
	}
}
