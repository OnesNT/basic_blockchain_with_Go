package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index        int
	Timestamp    string
	Data         string
	Hash         string
	PrevHash     string
	Transactions []Transaction
	Validator    string
	TokensReward int
}

type Transaction struct {
	From   string
	To     string
	Amount float64
}

type SmartContract struct {
	ContractCode string
}

type Account struct {
	ID      string
	Balance float64
}

type Wallet struct {
	Owner  string
	Tokens int
}

var validators = []string{"validator1", "validator2", "validator3"}

var wallets = make(map[string]Wallet)

var Blockchain []Block

func calculateHash(block Block) string {
	// Concatenate block fields into a single string
	record := string(block.Index) + block.Timestamp + block.Data + block.PrevHash

	// Convert the concatenated string into a byte slice (ASCII)
	rb := []byte(record)

	// Create a new SHA-256 hash instance
	h := sha256.New()

	// Write the byte slice to the hash instance
	h.Write(rb)

	// Finalize the hash computation and get the resulting hash as a byte slice
	hashed := h.Sum(nil) // nil indicates that no extra data should be appended to the current hash state before finalization.

	// Encode the byte slice to a hexadecimal string and return it
	return hex.EncodeToString(hashed)
}

// Generate new block from previous block in blockchain except first block

func generateBlock(oldBlock Block, transactions []Transaction, validator string) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	newBlock.Validator = validator
	newBlock.Transactions = transactions

	// Reward the validator with tokens (for example, 1 token)
	newBlock.TokensReward = 1

	return newBlock, nil
}

// Assert that block is valid or not
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func createTransaction(from, to string, amount float64) Transaction {
	return Transaction{from, to, amount}
}

func executeSmartContract(contract SmartContract, transactions []Transaction, accounts map[string]*Account) {
	fmt.Println("Executing smart contract:", contract.ContractCode)
	fmt.Println("Transactions:")

	for _, tx := range transactions {
		fmt.Printf("%s -> %s : %.2f\n", tx.From, tx.To, tx.Amount)

		fromAccount := accounts[tx.From]
		toAccount := accounts[tx.To]

		if fromAccount.Balance >= tx.Amount {
			fromAccount.Balance -= tx.Amount
			toAccount.Balance += tx.Amount
			fmt.Printf("Transaction successful: %s -> %s : %.2f\n", tx.From, tx.To, tx.Amount)
		} else {
			fmt.Printf("Transaction failed: %s has insufficient balance\n", tx.From)
		}
	}

	fmt.Println("Final account balances:")
	for id, account := range accounts {
		fmt.Printf("Account %s: %.2f\n", id, account.Balance)
	}
}

func isValidValidator(validator string) bool {
	for _, v := range validators {
		if v == validator {
			return true
		}
	}
	return false
}

func addToWallet(owner string, amount int) {
	if wallet, ok := wallets[owner]; ok {
		wallets[owner] = Wallet{owner, wallet.Tokens + amount}
	} else {
		wallets[owner] = Wallet{owner, amount}
	}
}

func transferTokens(from, to string, amount int) {
	if wallet, ok := wallets[from]; ok && wallet.Tokens >= amount {
		wallets[from] = Wallet{from, wallet.Tokens - amount}
		addToWallet(to, amount)
		fmt.Printf("Transfer %d tokens from %s to %s successful.\n", amount, from, to)
	} else {
		fmt.Println("Transfer failed. Insufficient balance or invalid wallet.")
	}
}

func main() {
	// Initialize accounts
	accounts := map[string]*Account{
		"Alice":   {ID: "Alice", Balance: 100.0},
		"Bob":     {ID: "Bob", Balance: 50.0},
		"Charlie": {ID: "Charlie", Balance: 20.0},
	}

	// Create the genesis block
	t := time.Now()
	genesisBlock := Block{0, t.String(), "Genesis Block", "", "", nil, "", 0}
	genesisBlock.Hash = calculateHash(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)
	fmt.Println("Genesis Block Hash:", genesisBlock.Hash)

	// Create transactions
	transactions := []Transaction{
		createTransaction("Alice", "Bob", 10.0),
		createTransaction("Bob", "Charlie", 5.0),
	}

	// Execute a smart contract
	contract := SmartContract{ContractCode: "Transfer funds if balance is sufficient"}
	executeSmartContract(contract, transactions, accounts)

	// Add transactions to a new block
	validator := "validator1"
	newBlock, _ := generateBlock(Blockchain[len(Blockchain)-1], transactions, validator)

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)
		addToWallet(validator, newBlock.TokensReward)
		fmt.Println("Block added successfully.")
		fmt.Println("Hash:", newBlock.Hash)
	}

	// Print final state of the blockchain
	fmt.Println("Blockchain:", Blockchain)

	// Transfer tokens between wallets
	addToWallet("Alice", 10)
	addToWallet("Bob", 5)
	transferTokens("Alice", "Bob", 5)

	// Print final state of wallets
	fmt.Println("Final wallets state:")
	for owner, wallet := range wallets {
		fmt.Printf("Wallet %s: %d tokens\n", owner, wallet.Tokens)
	}
}
