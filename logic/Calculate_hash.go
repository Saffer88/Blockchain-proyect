package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func CalculateHash(block Block) string {
	record := fmt.Sprintf("%d%d%s%d", block.Index, block.Timestamp, block.PrevHash, block.Nonce)
	for _, transaction := range block.Transactions {
		record += fmt.Sprintf("%s%s%f", transaction.Sender, transaction.Receiver, transaction.Amount)
	}
	h := sha256.New()
	h.Write([]byte(record))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
