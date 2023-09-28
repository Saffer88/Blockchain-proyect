package logic

import (
	"time"
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

func GenerateBlock(prevBlock Block, transactions []Transaction, limit int) Block {
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
