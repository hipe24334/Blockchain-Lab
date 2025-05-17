package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Transaction represents a transaction in a block.
type Transaction struct {
	Data []byte
}
type TransactionDto struct {
	Data string
}

func (tran *Transaction) convertTransactionDto() *TransactionDto {
	return &TransactionDto{Data: string(tran.Data)}
}
func convertTransactionsDto(trans []*Transaction) []*TransactionDto {
	var transDto = []*TransactionDto{}
	for _, t := range trans {
		transDto = append(transDto, t.convertTransactionDto())
	}
	return transDto
}

var transactions []*Transaction = []*Transaction{}

func TransactionsToString(trans []*Transaction) string {
	var transactionStrings []string

	for _, t := range trans {
		transactionStrings = append(transactionStrings, string(t.Data))
	}

	return fmt.Sprintf("[%s]", strings.Join(transactionStrings, ","))
}

// Block represents a block in the blockchain.
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
}
type BlockDto struct {
	Timestamp     int64
	Transactions  []*TransactionDto
	PrevBlockHash []byte
	Hash          []byte
}

func (block *Block) convertBlockDto() *BlockDto {
	return &BlockDto{
		Timestamp:     block.Timestamp,
		Transactions:  convertTransactionsDto(block.Transactions),
		PrevBlockHash: block.PrevBlockHash,
		Hash:          block.Hash,
	}
}
func convertBlocksDto(blocks []*Block) []*BlockDto {
	var blockDto = []*BlockDto{}
	for _, t := range blocks {
		blockDto = append(blockDto, t.convertBlockDto())
	}
	return blockDto
}
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

func (b *Block) DeriveHash() {
	var transactionsData []byte
	for _, tx := range b.Transactions {
		transactionsData = append(transactionsData, []byte(tx.Data)...)
	}
	info := bytes.Join(
		[][]byte{
			[]byte(time.Unix(b.Timestamp, 0).Format(time.RFC3339)),
			transactionsData,
			b.PrevBlockHash,
		},
		[]byte{},
	)
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// CreateBlock creates a new block.
func CreateBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}
	block.DeriveHash()
	return block
}

// Genesis creates the genesis block.
func Genesis() *Block {
	return CreateBlock([]*Transaction{
		{Data: []byte("Genesis")},
	}, []byte{})
}

// Blockchain represents the blockchain.
type Blockchain struct {
	blocks []*Block
}

func (chain *Blockchain) GetBlocks() []*Block {
	return chain.blocks
}

// AddBlock adds a new block to the blockchain.
func (chain *Blockchain) AddBlock(transactions []*Transaction) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(transactions, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}

// Init.
func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

func (t *Transaction) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(t)

	Handle(err)

	return res.Bytes()
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// ///////////////////////////////////////////////////////////////
func (chain *Blockchain) GetTransactions() []*Transaction {
	return transactions
}
func (chain *Blockchain) SetTransactions(trans []*Transaction) {
	transactions = trans
}
func (chain *Blockchain) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jbytes, err := json.MarshalIndent(convertTransactionsDto(transactions), "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jbytes)
}
func (chain *Blockchain) AddTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not write Block: %v", err)
		w.Write([]byte("could not write block"))
		return
	}
	transactions = append(transactions, &Transaction{Data: []byte(data)})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
}

// //////////////////////////////////////////////////////////////
func (chain *Blockchain) GetBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jbytes, err := json.MarshalIndent(convertBlocksDto(chain.blocks), "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jbytes)
}
func (chain *Blockchain) AddNewBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(transactions) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Transactions is Empty"})
		return
	}

	chain.AddBlock(transactions)

	jbytes, err := json.MarshalIndent(chain.blocks[len(chain.blocks)-1].convertBlockDto(), "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}
	transactions = []*Transaction{}
	w.WriteHeader(http.StatusOK)
	w.Write(jbytes)
}
