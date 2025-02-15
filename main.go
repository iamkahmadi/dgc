package main

import "dgc/app"

// func main() {
// 	// Create a new wallet
// 	senderWallet := wallet.NewWallet()

// 	// Display the wallet details
// 	fmt.Println(wallet.String(senderWallet))

// 	// Define recipient address and amount for the transaction
// 	recipient := "6013c9a2a701482af704412814b483767db8ec53d1bbeaf8ddd4bbcba155128ee26ceb134e2fb81bb84913dcdac71ad61630b061636c4ee1aa54c21245bca6fc" // Replace with an actual public key

// 	amount := 50 // Example amount to send

// 	// Declare the transaction pool outside the block so it is accessible later
// 	txPool := wallet.NewTransactionPool()

// 	// Create a new transaction
// 	tx := wallet.NewTransaction(senderWallet, recipient, amount)

// 	if tx != nil {
// 		// Display transaction details
// 		fmt.Printf("Transaction created:\nID: %s\n", tx.ID)
// 		for _, output := range tx.Outputs {
// 			fmt.Printf("Output - Amount: %d, Address: %s\n", output.Amount, output.Address)
// 		}

// 		// Verify the transaction
// 		isValid := wallet.VerifyTransaction(tx)
// 		if isValid {
// 			fmt.Println("Transaction is valid.")
// 		} else {
// 			fmt.Println("Transaction is invalid.")
// 		}

// 		// Add the transaction to the pool
// 		txPool.UpdateOrAddTransaction(tx)
// 		fmt.Println("Added transaction to the pool.")

// 		// Check valid transactions in the pool
// 		validTxns := wallet.ValidTransactions(txPool)
// 		for _, validTxn := range validTxns {
// 			fmt.Printf("Valid Transaction in Pool ID: %s\n", validTxn.ID)
// 		}
// 	} else {
// 		fmt.Println("Transaction could not be created due to insufficient balance.")
// 	}

// 	// Example of updating the transaction (if needed)
// 	newRecipient := "new-recipient-public-key-example" // Replace with an actual public key
// 	newAmount := 20                                    // Example amount to send

// 	if tx != nil {
// 		wallet.UpdateTransaction(tx, senderWallet, newRecipient, newAmount)

// 		fmt.Printf("Transaction updated:\nID: %s\n", tx.ID)
// 		for _, output := range tx.Outputs {
// 			fmt.Printf("Output - Amount: %d, Address: %s\n", output.Amount, output.Address)
// 		}

// 		// Verify the updated transaction
// 		isValid := wallet.VerifyTransaction(tx)
// 		if isValid {
// 			fmt.Println("Updated transaction is valid.")
// 		} else {
// 			fmt.Println("Updated transaction is invalid.")
// 		}

// 		// Check valid transactions in the pool after update
// 		// Use txPool declared outside of blocks
// 		validTxns := wallet.ValidTransactions(txPool)
// 		for _, validTxn := range validTxns {
// 			fmt.Printf("Valid Transaction in Pool ID: %s\n", validTxn.ID)
// 		}
// 	}
// }

// for checking the api and p2p-server
func main() {
	// Start the application (uncomment if you have an app to start)
	app.App()
}

// for checking blockchain
// func main() {
// 	// Create a new Blockchain
// 	bc := blockchain.NewBlockchain()

// 	// Display the Genesis Block
// 	fmt.Println("Genesis Block: ")
// 	fmt.Println(bc.Chain[0].ToString())

// 	// Add a new block to the blockchain 1
// 	data1 := []string{"Some transaction data 1"}
// 	newBlock1 := bc.AddBlock(data1)
// 	fmt.Println("\nNewly mined block 1: ")
// 	fmt.Println(newBlock1.ToString())

// 	// Add a new block to the blockchain
// 	data2 := []string{"Some transaction data 2"}
// 	newBlock2 := bc.AddBlock(data2)
// 	fmt.Println("\nNewly mined block 2: ")
// 	fmt.Println(newBlock2.ToString())

// 	// Display the current blockchain
// 	fmt.Println("\nCurrent Blockchain: ")
// 	for _, blk := range bc.Chain {
// 		fmt.Println(blk.ToString())
// 	}

// 	// Validate the current chain
// 	if bc.IsValidChain(bc.Chain) {
// 		fmt.Println("\nBlockchain is valid.")
// 	} else {
// 		fmt.Println("\nBlockchain is invalid.")
// 	}
// }
