package main

// import (
// 	"fmt"
// 	"golang/blockchain"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// func main() {

// 	var chain = blockchain.InitBlockChain()

// 	r := mux.NewRouter()
// 	r.HandleFunc("/addTransaction", chain.AddTransaction).Methods("POST")
// 	r.HandleFunc("/getTransaction", chain.GetTransaction).Methods("GET")
// 	r.HandleFunc("/getBlockchain", chain.GetBlockchain).Methods("GET")
// 	r.HandleFunc("/addNewBlock", chain.AddNewBlock).Methods("POST")
// 	// r.HandleFunc("/new", newBook).Methods("POST")

// 	go func() {

// 		for _, block := range chain.GetBlocks() {
// 			fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
// 			fmt.Printf("Transaction in Block: %v\n", blockchain.TransactionsToString(block.Transactions))
// 			fmt.Printf("Hash: %x\n", block.Hash)
// 			fmt.Println()
// 		}

// 	}()
// 	log.Println("Listening on port 3000")

// 	log.Fatal(http.ListenAndServe(":3000", r))
// }

import (
	"fmt"
	"golang/blockchain"
	"os"
)

func main() {
	var chain = blockchain.InitBlockChain()

	fmt.Println("Blockchain Console:")
	fmt.Println("--------------------")

	for {
		fmt.Println("\nSelect an operation:")
		fmt.Println("1. Add Transaction")
		fmt.Println("2. Add New Block")
		fmt.Println("3. Display Blockchain")
		fmt.Println("4. Display Transactions")
		fmt.Println("5. Exit")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			addTransaction(chain)
		case 2:
			addNewBlock(chain)
		case 3:
			displayBlockchain(chain)
		case 4:
			displayTransactions(chain)
		case 5:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func addTransaction(chain *blockchain.Blockchain) {
	var data string
	fmt.Print("Enter transaction data: ")
	fmt.Scanln(&data)
	trans := chain.GetTransactions()
	trans = append(trans, &blockchain.Transaction{Data: []byte(data)})
	chain.SetTransactions(trans)
	fmt.Println("Transaction added successfully.")
}

func addNewBlock(chain *blockchain.Blockchain) {
	if len(chain.GetBlocks()) == 0 {
		fmt.Println("Error: Cannot add a new block without any transactions.")
		return
	}
	chain.AddBlock(chain.GetTransactions())
	chain.SetTransactions([]*blockchain.Transaction{})
	fmt.Println("New block added successfully.")
}

func displayBlockchain(chain *blockchain.Blockchain) {
	blocks := chain.GetBlocks()

	for _, block := range blocks {
		fmt.Printf("\nPrev Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Transaction in Block: %v\n", blockchain.TransactionsToString(block.Transactions))
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}

func displayTransactions(chain *blockchain.Blockchain) {
	fmt.Println("Transactions:")
	transactions := blockchain.TransactionsToString(chain.GetTransactions())
	fmt.Println(transactions)
}
