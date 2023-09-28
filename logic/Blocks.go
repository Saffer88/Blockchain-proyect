package logic

import (
	"fmt"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Limit        int
	Nonce        int
}

func GenerateBlock(db *leveldb.DB, transactions []Transaction, limit int) Block {
	prevBlock, err := GetLastBlock(db)
	if err != nil {
		fmt.Println("Error al obtener el último bloque:", err)
		panic(err)
	}
	var block Block
	block.Index = prevBlock.Index + 1
	block.Timestamp = time.Now().Unix()
	block.Transactions = transactions
	block.PrevHash = prevBlock.Hash
	block.Nonce = prevBlock.Nonce + 1
	block.Limit = limit
	block.Hash = CalculateHash(block)
	return block
}

func Limit(db *leveldb.DB) bool {

	lastBlock, err := GetLastBlock(db)
	if err != nil {
		fmt.Println("Error al obtener el último bloque:", err)
	}
	limit, err := GetTotalTransactions(db, lastBlock.Index)
	if err != nil {
		fmt.Println("Error al obtener la cantidad de transacciones:", err)
	}
	if lastBlock.Limit > limit {
		return false
	}

	return true
}

func Genesis() Block {
	var block Block
	block.Index = 1
	block.Timestamp = time.Now().Unix()
	block.Transactions = nil
	block.PrevHash = "1231231231231231231231231231231231231231231231231231231231231231"
	block.Nonce = 1
	block.Limit = 0
	block.Hash = CalculateHash(block)
	return block

}
