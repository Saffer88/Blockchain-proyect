package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	Signature string
}

func NewTransaction(sender, receiver string, amount float64) Transaction {
	transaction := Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Signature: SignTransaction(sender, receiver, amount),
	}
	return transaction
}

func SignTransaction(sender string, receiver string, amount float64) string {
	transactionDetails := fmt.Sprintf("%s%s%f", sender, receiver, amount)
	h := sha256.New()
	h.Write([]byte(transactionDetails))
	hash := hex.EncodeToString(h.Sum(nil))
	// Firmar el hash con la clave privada del emisor
	signature := sender + hash
	return signature
}

func AddTransaction(block *Block, transaction Transaction) {
	block.Transactions = append(block.Transactions, transaction)
}
