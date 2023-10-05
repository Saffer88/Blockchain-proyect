package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func CalculateHash(block Block) string {
	record := fmt.Sprintf("%d%d%s", block.Index, block.Timestamp, block.PrevHash)
	for _, transaction := range block.Transactions {
		record += fmt.Sprintf("%s%s%f%d", transaction.Sender, transaction.Receiver, transaction.Amount, transaction.Nonce)
	}
	h := sha256.New()
	h.Write([]byte(record))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
