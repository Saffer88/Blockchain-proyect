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

func UpdateBlockHash(block *Block) {
	block.Hash = CalculateHash(*block)
}

// Cuenta genesis
// Private Key: 5389340a76f6ac7f16dd3accf1ba2fd8cc505451be96601955cf234d4f0915d3
// Public Key: 7d6866b740b19acdef6055398dfb2ace996153099471f12f6f1ff19d7856157ddf49cd8bca60f22467dbac63df248e06adbde6246d0c77385b3dd592e3ae31a1
// Address: 0e1dd7f2e5cb568ee13534424aaa978e484df040

func Genesis() Block {
	var block Block
	block.Index = 1
	block.Timestamp = time.Now().Unix()
	block.Transactions = []Transaction{
		{
			Sender:    "",
			Receiver:  "0e1dd7f2e5cb568ee13534424aaa978e484df040",
			Amount:    1000000,
			Signature: "",
			Nonce:     0,
		},
	}
	block.PrevHash = "0000000000000000000000000000000000000000000000000000000000000000"
	block.Limit = 0
	block.Hash = CalculateHash(block)
	return block
}
