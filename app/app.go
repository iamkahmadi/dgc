package app

import (
	"dgc/blockchain"
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

// Struct to capture the request body
type MineRequest struct {
	Data string `json:"data"`
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

	// Convert the string data to a slice of strings
	data := []string{req.Data}

	// Add the new block to the blockchain using the converted data
	block := bc.AddBlock(data)

	// Log the added block
	log.Printf("New block added: %s", block.ToString())

	// Respond with the newly mined block
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(block)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Optionally, you can sync with other peers if you have a P2P server
	p2pServer.syncChains()
}

func App() {
	// Routes
	http.HandleFunc("/blocks", getBlocks)
	http.HandleFunc("/mine", mineBlock)

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
