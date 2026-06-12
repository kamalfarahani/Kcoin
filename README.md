# Kcoin

A minimal blockchain and cryptocurrency implementation written in Go, built for
learning how the core pieces of a Bitcoin-style system fit together: proof-of-work
mining, a persistent block store, ECDSA-signed transactions, a UTXO model, and
Base58Check wallet addresses.

> **Status:** educational project. It is intentionally small and not intended for
> production or real value transfer.

---

## Table of contents

- [Features](#features)
- [How it works](#how-it-works)
- [Requirements](#requirements)
- [Getting started](#getting-started)
- [Using the library](#using-the-library)
- [Project layout](#project-layout)
- [API reference](#api-reference)
- [Configuration constants](#configuration-constants)
- [Limitations & known issues](#limitations--known-issues)
- [License](#license)

---

## Features

- **Proof-of-work mining** — blocks are sealed by finding a nonce whose SHA-256
  hash falls below a difficulty target.
- **Persistent storage** — blocks are stored in a [BoltDB](https://github.com/boltdb/bolt)
  file (`assets/blockchain.db`), so the chain survives restarts.
- **UTXO transaction model** — balances are derived by scanning unspent transaction
  outputs, the same model Bitcoin uses.
- **Digital signatures** — transactions are signed and verified with ECDSA on the
  P-256 curve.
- **Wallets & addresses** — public keys are hashed (SHA-256 → RIPEMD-160) and
  encoded as Base58Check addresses with a version byte and checksum.
- **Coinbase transactions** — every block's first transaction mints the mining
  subsidy to the miner.

## How it works

Kcoin stitches together five concepts:

1. **Wallets.** A wallet is an ECDSA key pair. The public key is hashed and
   Base58Check-encoded to produce a human-shareable **address**.
2. **Transactions.** A transaction spends previous **outputs** (referenced by
   inputs) and creates new outputs locked to recipient address hashes. The first
   transaction in a block is a **coinbase** that has no inputs and mints new coins.
3. **The UTXO set.** There is no account balance stored anywhere. To compute a
   balance, Kcoin walks the chain and collects every output locked to an address
   that has not yet been spent by a later input.
4. **Proof-of-work.** Before a block is accepted, `ProofOfWork.Mine` searches for a
   nonce so that `SHA-256(block data)` is numerically smaller than the difficulty
   target. The target is derived from `targetBits` (currently 24).
5. **The chain & store.** Blocks link to their predecessor via `PrevBlockHash`.
   They are serialized with `encoding/gob` and persisted in a BoltDB bucket, with a
   special key tracking the current chain tip.

```
Wallet ──address──▶ Transaction ──signed──▶ Block ──mined──▶ Blockchain ──▶ BoltDB
```

## Requirements

- **Go 1.26+** (the module targets the toolchain in `go.mod`)
- No other system dependencies — Go modules fetch everything else.

Third-party libraries (pulled automatically by `go mod`):

| Module | Purpose |
| --- | --- |
| `github.com/boltdb/bolt` | Embedded key/value store for blocks |
| `golang.org/x/crypto/ripemd160` | RIPEMD-160 hashing for address generation |

## Getting started

Clone the repository and build:

```bash
git clone git@github.com:kamalfarahani/Kcoin.git
cd Kcoin
go build ./...
```

Run the bundled demo, which creates a new wallet, initializes the blockchain
(mining the genesis block on first run), and prints the wallet's address and
balance:

```bash
go run ./cmd/kcoin
```

Example output:

```
New wallet address: 1A2b3C...your-address...
Balance of 1A2b3C...your-address...: 10
```

> The first run mines the **genesis block**, which awards the coinbase subsidy
> (10 coins) to the demo wallet. Mining performs proof-of-work, so it may take a
> moment. The chain is written to `assets/blockchain.db`; delete that file to start
> from a fresh genesis block.

## Using the library

The blockchain logic lives in the `kcoin` package
(`github.com/kamalfarahani/Kcoin/internal/kcoin`). A typical flow:

```go
package main

import (
	"fmt"

	"github.com/kamalfarahani/Kcoin/internal/kcoin"
)

func main() {
	// 1. Create wallets (ECDSA key pairs).
	miner := kcoin.NewWallet()
	recipient := kcoin.NewWallet()

	// 2. Initialize (or open) the blockchain, mining genesis to the miner.
	chain := kcoin.NewBlockchain(miner.GetAddress())
	manager := kcoin.NewBlockchainManager(chain)

	// 3. Send coins from the miner to another address.
	manager.SendCoin(*miner, recipient.GetAddress(), 5)

	// 4. Check balances (derived from the UTXO set).
	fmt.Println("Miner balance:    ", manager.GetBalance(miner.GetAddress()))
	fmt.Println("Recipient balance:", manager.GetBalance(recipient.GetAddress()))
}
```

> **Note:** `internal/kcoin` is an internal package, so it can only be imported by
> code inside this module (such as `cmd/`). If you want to consume Kcoin as a
> library from another module, move the package to `pkg/kcoin` and update imports.

## Project layout

```
Kcoin/
├── cmd/
│   └── kcoin/
│       └── main.go          # Demo entry point (wallet → chain → balance)
├── internal/
│   └── kcoin/               # Core library (package kcoin)
│       ├── base58.go            # Base58 encode/decode
│       ├── block.go             # Block type, mining a new block, genesis
│       ├── block_serialize.go   # gob (de)serialization for blocks
│       ├── blockchain.go        # Blockchain type, AddBlock, open/create
│       ├── blockchain_iterator.go # Walk the chain tip → genesis
│       ├── blockchain_manager.go  # Balances, UTXO lookup, build transactions
│       ├── database.go          # BoltDB persistence
│       ├── proofOfWork.go       # Proof-of-work mine/validate
│       ├── transaction.go       # Transaction & coinbase construction
│       ├── transaction_input.go # Inputs, signature/pubkey helpers
│       ├── transaction_output.go# Outputs, address locking
│       ├── transaction_manager.go # Sign & verify transactions
│       ├── transaction_serialize.go # gob (de)serialization for transactions
│       ├── utils.go             # Shared helpers (int↔bytes, etc.)
│       ├── wallet.go            # Wallet (ECDSA key pair)
│       └── wallet_address.go    # Address generation & decoding
├── assets/
│   └── blockchain.db        # BoltDB data file (created at runtime, git-ignored)
├── go.mod
└── go.sum
```

## API reference

The most useful exported types and functions in `package kcoin`:

### Wallets

| Symbol | Description |
| --- | --- |
| `NewWallet() *Wallet` | Generate a new ECDSA (P-256) key pair. |
| `(*Wallet).GetAddress() []byte` | Base58Check address derived from the public key. |

### Blockchain

| Symbol | Description |
| --- | --- |
| `NewBlockchain(address []byte) *Blockchain` | Open the DB, or create it and mine the genesis block to `address`. |
| `(*Blockchain).AddBlock(txs []*Transaction)` | Verify the transactions, mine a block, and append it. |
| `(*Blockchain).Iterator() *BlockchainIterator` | Iterate from the tip back to genesis. |
| `(*BlockchainIterator).Next() *Block` | Return the current block and step toward genesis. |

### Blockchain manager (high-level operations)

| Symbol | Description |
| --- | --- |
| `NewBlockchainManager(bc *Blockchain) *BlockchainManager` | Wrap a chain for balance/transaction operations. |
| `(*BlockchainManager).GetBalance(address []byte) int` | Sum unspent outputs for an address. |
| `(*BlockchainManager).SendCoin(w Wallet, to []byte, amount int)` | Build, sign, and mine a transfer in a new block. |
| `(*BlockchainManager).FindUnspentTransactions(address []byte) map[*Transaction][]uint` | UTXO lookup used to compute balances. |

### Transactions, blocks & proof-of-work

| Symbol | Description |
| --- | --- |
| `NewTransaction(inputs, outputs) *Transaction` | Construct a transaction and set its ID. |
| `NewCoinbaseTransaction(to []byte, data string) *Transaction` | Mint the subsidy to an address. |
| `(*Transaction).IsCoinbase() bool` | Whether the transaction is a coinbase. |
| `NewBlock(txs, prevHash) *Block` | Build and mine a block. |
| `NewProofOfWork(b *Block) *ProofOfWork` | Create a PoW solver for a block. |
| `(*ProofOfWork).Mine() (int64, []byte, error)` | Search for a valid nonce. |
| `(*ProofOfWork).Validate() bool` | Verify a block's proof-of-work. |

### Encoding

| Symbol | Description |
| --- | --- |
| `Base58Encode(input []byte) []byte` | Encode bytes as Base58. |
| `Base58Decode(input []byte) []byte` | Decode Base58 bytes. |

## Configuration constants

These tunables are defined as package constants in `internal/kcoin`:

| Constant | Value | Meaning |
| --- | --- | --- |
| `targetBits` | `24` | Mining difficulty (higher = harder). |
| `subsidy` | `10` | Coins awarded by a coinbase transaction. |
| `version` | `0x00` | Address version byte. |
| `addressChecksumLen` | `4` | Bytes of checksum appended to an address. |
| `dbFilePath` | `./assets/blockchain.db` | BoltDB data file location. |
| `genesisCoinbaseData` | `"Kamal Genesis Block"` | Arbitrary data embedded in the genesis coinbase. |

## Limitations & known issues

This is a learning project; some real-world concerns are deliberately out of scope:

- **No networking.** It is a single-node chain with no peer-to-peer layer or
  consensus across nodes.
- **No CLI.** `cmd/kcoin` is a fixed demo rather than a command-line wallet/node.
- **Single, fixed difficulty.** `targetBits` is constant; there is no difficulty
  retargeting.
- **Relative DB path.** The database is opened at `./assets/blockchain.db` relative
  to the working directory, so run commands from the repository root.
- **Errors are often panics.** Many failure paths call a panic helper rather than
  returning errors to the caller.

