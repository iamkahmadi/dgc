package app

import (
	"dgc/block"
	"dgc/blockchain"
	"dgc/types"
	"dgc/wallet"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	// First, check if the environment variable exists
	HTTP_PORT = getEnv("HTTP_PORT", "3001") // If HTTP_PORT is not set, it will default to "3001"
	bc        = blockchain.NewBlockchain()
	tp        = wallet.NewTransactionPool()
	p2pServer = NewP2pServer(bc)
)

// Helper function to get environment variable or fallback to default value
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Struct to capture the request body, now expecting a transaction
type MineRequest struct {
	TransactionData string `json:"transaction_data"` // Assuming transaction data is passed as a string
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(bc.Chain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mineBlock(w http.ResponseWriter, r *http.Request) {
	// Read the body of the request
	var req MineRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the JSON body into the MineRequest struct
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal the transaction data into a *types.Transaction
	var transaction *types.Transaction
	err = json.Unmarshal([]byte(req.TransactionData), &transaction)
	if err != nil {
		http.Error(w, "Error unmarshaling transaction data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Add the transaction to the transaction pool
	tp.UpdateOrAddTransaction(transaction) // Add this line

	// Add the new block to the blockchain with the transaction data
	// data := []*types.Transaction{transaction} // Create a slice containing the transaction
	// newBlock := bc.AddBlock(data)             // Call AddBlock with the transaction slice

	// just checking
	newBlock := bc.AddBlock(tp.Transactions)

	// Log the added block
	log.Printf("New block added: %s", block.ToString(&newBlock))

	// Respond with the newly mined block
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newBlock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Optionally, you can sync with other peers if you have a P2P server
	p2pServer.syncChains()
}

// direclyt adds transaction with putting it in pool
// func mineBlock(w http.ResponseWriter, r *http.Request) {
// 	// Read the body of the request
// 	var req MineRequest
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Parse the JSON body into the MineRequest struct
// 	err = json.Unmarshal(body, &req)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Here, you would typically unmarshal the transaction data into a *types.Transaction
// 	// For simplicity, let's assume the transaction data is passed as a JSON string
// 	var transaction *types.Transaction
// 	err = json.Unmarshal([]byte(req.TransactionData), &transaction)
// 	if err != nil {
// 		http.Error(w, "Error unmarshaling transaction data: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Add the new block to the blockchain with the transaction data
// 	data := []*types.Transaction{transaction} // Create a slice containing the transaction
// 	newBlock := bc.AddBlock(data)             // Call AddBlock with the transaction slice

// 	// Log the added block
// 	log.Printf("New block added: %s", block.ToString(&newBlock))

// 	// Respond with the newly mined block
// 	w.Header().Set("Content-Type", "application/json")
// 	err = json.NewEncoder(w).Encode(newBlock)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Optionally, you can sync with other peers if you have a P2P server
// 	p2pServer.syncChains()
// }

func getTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tp.Transactions) // Access the exported field
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	// Parse the request to get the transaction details
	var req struct {
		Recipient string `json:"recipient"`
		Amount    int    `json:"amount"`
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve or create a wallet (you may want to implement your own logic here)
	senderWallet := wallet.NewWallet() // This creates a new wallet each time; replace with your logic

	// Create a transaction
	transaction := wallet.CreateTransaction(senderWallet, req.Recipient, req.Amount, bc, tp)

	if transaction == nil {
		http.Error(w, "Failed to create transaction", http.StatusBadRequest)
		return
	}

	// Broadcast the transaction to peers
	p2pServer.BroadcastTransaction(transaction)

	// Respond with the created transaction
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func App() {
	// Routes
	http.HandleFunc("/blocks", getBlocks)
	http.HandleFunc("/mine", mineBlock)
	http.HandleFunc("/transactions", getTransactions)
	http.HandleFunc("/transact", createTransaction)

	// Show a message before starting the server
	log.Printf("Starting server on port %s...\n", HTTP_PORT)

	// Use sync.WaitGroup to keep the main thread alive
	var wg sync.WaitGroup
	wg.Add(2) // We have two goroutines: one for the HTTP server and one for the P2P server

	// Start the HTTP server in a separate goroutine
	go func() {
		defer wg.Done()
		log.Printf("Starting HTTP server on port %s...\n", HTTP_PORT)
		err := http.ListenAndServe(":"+HTTP_PORT, nil)
		if err != nil {
			log.Fatalf("Error starting HTTP server: %v\n", err)
		}
	}()

	// Start the P2P server (this is already running in its own goroutine in Listen())
	go func() {
		defer wg.Done()
		log.Printf("Starting P2P server on port %s...\n", P2P_PORT)
		p2pServer.Listen()
	}()

	// Wait for both goroutines to finish
	wg.Wait()
}
