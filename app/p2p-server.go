package app

import (
	"dgc/block"
	"dgc/blockchain"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	CHAIN              MessageType = "CHAIN"
	TRANSACTION        MessageType = "TRANSACTION"
	CLEAR_TRANSACTIONS MessageType = "CLEAR_TRANSACTIONS"
)

var P2P_PORT = "5001"
var peers []string

func init() {
	// Initialize peers from environment variable
	if peersEnv := os.Getenv("PEERS"); peersEnv != "" {
		peers = strings.Split(peersEnv, ",")
	}
	if p2pPort := os.Getenv("P2P_PORT"); p2pPort != "" {
		P2P_PORT = p2pPort
	}
}

type P2pServer struct {
	bc      *blockchain.Blockchain
	sockets []*websocket.Conn
}

func NewP2pServer(bc *blockchain.Blockchain) *P2pServer {
	return &P2pServer{
		bc:      bc,
		sockets: []*websocket.Conn{},
	}
}

func (server *P2pServer) Listen() {
	// Start WebSocket server
	http.HandleFunc("/", server.handleWebSocketConnection)

	go func() {
		// Listen for incoming WebSocket connections
		log.Printf("Listening for peer-to-peer connections on: %s", P2P_PORT)
		if err := http.ListenAndServe(":"+P2P_PORT, nil); err != nil {
			log.Fatalf("Error starting WebSocket server: %v", err)
		}
	}()

	// Connect to peers
	server.connectToPeers()
}

func (server *P2pServer) connectToPeers() {
	for _, peer := range peers {
		go server.connectToPeer(peer)
	}
}

func (server *P2pServer) connectToPeer(peer string) {
	log.Printf("Attempting to connect to peer: %s\n", peer)
	for {
		conn, _, err := websocket.DefaultDialer.Dial(peer, nil)
		if err != nil {
			log.Printf("Error connecting to peer %s: %v, retrying...", peer, err)
			time.Sleep(2 * time.Second) // Retry after a 2-second delay
			continue
		}
		server.connectSocket(conn)
		break // Break out of loop once successfully connected
	}
}

func (server *P2pServer) connectSocket(socket *websocket.Conn) {
	server.sockets = append(server.sockets, socket)
	log.Println("Socket connected")

	// Handle incoming messages
	go server.messageHandler(socket)

	// Send the blockchain to the connected peer
	server.sendChain(socket)
}

func (server *P2pServer) handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	server.connectSocket(conn)
}

func (server *P2pServer) messageHandler(socket *websocket.Conn) {
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		var data struct {
			Type  MessageType   `json:"type"`
			Chain []block.Block `json:"chain"`
		}

		err = json.Unmarshal(message, &data)
		if err != nil {
			log.Println("Error unmarshaling message:", err)
			continue
		}

		// match the chain for debugging purpose
		// log.Printf("Current chain length: %d", len(server.bc.Chain))
		// log.Printf("Received chain length: %d", len(data.Chain))
		// for i, block := range data.Chain {
		// 	log.Printf("Received Block %d: Hash = %s", i, block.Hash)
		// }

		switch data.Type {
		case CHAIN:
			server.bc.ReplaceChain(data.Chain)
		}

	}
}

func (server *P2pServer) sendChain(socket *websocket.Conn) {
	message := struct {
		Type  MessageType   `json:"type"`
		Chain []block.Block `json:"chain"`
	}{
		Type:  CHAIN,
		Chain: server.bc.Chain,
	}

	err := socket.WriteJSON(message)
	if err != nil {
		log.Println("Error sending chain:", err)
	}
}

func (server *P2pServer) syncChains() {
	for _, socket := range server.sockets {
		server.sendChain(socket)
	}
}
