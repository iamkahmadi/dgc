

# Digital Governance Chain (DGC)

**DGC** is a blockchain for enabling transparent, decentralized digital governance systems. It provides core functionality for peer-to-peer networking, block validation, mining, and wallet-based transactions—all implemented in Go.

> 🚀 **Project by:** [iamkahmadi](https://github.com/iamkahmadi)

---

## 🧩 Description

The **Digital Governance Chain** (DGC) is a lightweight and extensible blockchain platform designed to support governance-related operations such as voting, resource allocation, identity management, and transparent public auditing.

With a modular architecture written in Go, DGC includes core blockchain logic, a transaction pool, mining functionality, and wallet support. The peer-to-peer network enables nodes to communicate and propagate blocks and transactions efficiently.

---

## 🛠️ Features

- ⛓️ Custom Blockchain Implementation  
- 🧾 Transaction Handling and Wallets  
- 🧠 Mining and Proof-of-Work Algorithm  
- 📡 P2P Server for Node Communication  
- ✅ Modular Codebase with Unit Tests  
- 🧪 Built-in Test Coverage for Core Components  

---

## 📁 Project Structure

```bash
.
├── app/                    # Node and application-level startup logic
│   ├── app.go
│   └── check files
├── block/                  # Block structure and utilities
│   ├── block.go
│   └── block_test.go
├── blockchain/             # Blockchain core logic
│   ├── blockchain.go
│   └── blockchain_test.go
├── miner/                  # Mining and Proof-of-Work
│   ├── miner.go
│   └── check files
├── p2p-server/             # Peer-to-peer networking
│   ├── p2p-server.go
│   └── check files
├── types/                  # Common type definitions
│   └── types.go
├── util/                   # Utility and config files
│   ├── chainutil.go
│   └── config.go
├── wallet/                 # Wallets and transaction pools
│   ├── wallet.go
│   ├── transaction.go
│   └── transactionpool.go
├── main.go                 # Entry point
├── go.mod                  # Go module definition
└── commands.txt / read.py # Miscellaneous helper scripts
```

---

## 🚀 Getting Started

### Prerequisites

- Go
- Git

### Installation

```bash
git clone https://github.com/iamkahmadi/dgc.git
cd dgc
go mod tidy
```

### Run the Node

```bash
go run main.go
```

### Run Tests

```bash
go test ./...
```

---

## 📡 P2P Network

DGC includes a basic peer-to-peer networking stack in the `p2p-server` package, allowing nodes to:

- Discover peers
- Share blocks and transactions
- Synchronize blockchain state

---

## 🔐 Wallets & Transactions

Users can create wallets and broadcast transactions securely using the `wallet` module. The transaction pool temporarily stores unconfirmed transactions until they are mined into blocks.

---

## ⛏️ Mining

The `miner` package implements a simplified Proof-of-Work consensus mechanism. Mining validates transactions and appends blocks to the blockchain.

---

## 🧪 Testing

Unit tests are available for the core packages. You can run all tests with:

```bash
go test ./...
```

---

## 📜 License

MIT License. See `LICENSE` for more information.

---

## 🤝 Contributing

Feel free to fork, submit pull requests, and open issues to contribute ideas or fixes.
