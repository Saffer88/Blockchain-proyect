package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	Signature string
}

func AddTransaction(block *Block, transaction Transaction) {
	block.Transactions = append(block.Transactions, transaction)
}

func NewTransaction(sender string, receiver string, amount float64, privateKeyHex string) (*Transaction, error) {
	// Convert hex private key to *ecdsa.PrivateKey
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	txData := fmt.Sprintf("%s%s%f", sender, receiver, amount)
	hash := sha256.Sum256([]byte(txData))

	sig, err := crypto.Sign(hash[:], privateKey) // firma
	if err != nil {
		return nil, err
	}

	// Construct Transaction object
	tx := &Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Signature: hex.EncodeToString(sig),
	}

	return tx, nil
}

func VerifyTransaction(sender string, receiver string, amount float64, signature, publicKeyHex string) bool {

	txData := fmt.Sprintf("%s%s%f", sender, receiver, amount)
	hash := sha256.Sum256([]byte(txData))

	sig, err := hex.DecodeString(signature)
	if err != nil {
		fmt.Println("Error decoding signature:", err)
		return false
	}

	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		fmt.Println("Error decoding public key:", err)
		return false
	}

	// Agrega el byte 0x04 al inicio del publicKeyBytes
	publicKeyBytes = append([]byte{0x04}, publicKeyBytes...)

	if len(publicKeyBytes) != 65 {
		fmt.Println("Invalid public key length")
		return false
	}

	publicKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		fmt.Println("Error unmarshaling public key:", err)
		return false
	}

	isValidSig := crypto.VerifySignature(crypto.FromECDSAPub(publicKey), hash[:], sig[:len(sig)-1])
	return isValidSig
}
