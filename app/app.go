package app

import (
	"dgc/block"
	"dgc/blockchain"
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
	HTTP_PORT   = getEnv("HTTP_PORT", "3001") // If HTTP_PORT is not set, it will default to "3001"
	bc          = blockchain.NewBlockchain()
	tp          = wallet.NewTransactionPool()
	p2pServer   = NewP2pServer(bc, tp)
	localWallet = wallet.NewWallet()
	localMiner  = NewMiner(bc, tp, localWallet, p2pServer)
)

// Helper function to get environment variable or fallback to default value
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(bc.Chain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mineBlock(w http.ResponseWriter, r *http.Request) {
	// Check if there are transactions in the transaction pool
	if len(tp.Transactions) == 0 {
		// Respond with a JSON message saying no transactions are available
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"error": "No transactions in transaction pool"}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Proceed with mining the block if there are transactions
	newBlock := bc.AddBlock(tp.Transactions)
	// Log the added block
	log.Printf("New block added: %s", block.ToString(&newBlock))
	tp.Clear()
	// Respond with the newly mined block in JSON format
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(newBlock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Optionally, sync with other peers if you have a P2P server
	p2pServer.syncChains()
}

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

func mineTransactions(w http.ResponseWriter, r *http.Request) {
	b := localMiner.Mine()

	// Log the added block
	log.Printf("New block added: %s", block.ToString(&b))

	// Respond with the newly mined block
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getPublicKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"publicKey": localMiner.Wallet.PublicKey})
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
	http.HandleFunc("/mine-transactions", mineTransactions)
	http.HandleFunc("/public-key", getPublicKey)

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
