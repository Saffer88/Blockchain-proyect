package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
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

func NewTransaction(sender string, receiver string, amount float64, privateKeyHex string, db *leveldb.DB) (*Transaction, error) {

	verify := VerifyBalance(sender, amount, db)

	if !verify {
		fmt.Println("\nNo existe el saldo suficiente en la cuenta.")
		return nil, nil
	}

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
		fmt.Println("Error al decodear la firma:", err)
		return false
	}

	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		fmt.Println("Error al decodear la llave publica: ", err)
		return false
	}

	publicKeyBytes = append([]byte{0x04}, publicKeyBytes...)

	if len(publicKeyBytes) != 65 {
		fmt.Println("Largo inválido")
		return false
	}

	publicKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		fmt.Println("Error en unmarshal", err)
		return false
	}

	isValidSig := crypto.VerifySignature(crypto.FromECDSAPub(publicKey), hash[:], sig[:len(sig)-1])

	if isValidSig {
		fmt.Print("\nLa firma es válida ")
	} else {
		fmt.Print("\nLa firma no es válida ")
	}
	return isValidSig
}

func CalculateBalance(address string, db *leveldb.DB) float64 {
	var balance float64

	iter := db.NewIterator(nil, nil)
	for iter.Next() {

		blockData := iter.Value()

		var block Block
		err := json.Unmarshal(blockData, &block)
		if err != nil {
			log.Fatal(err)
		}

		for _, tx := range block.Transactions {
			if tx.Sender == address {
				balance -= tx.Amount
			}
			if tx.Receiver == address {
				balance += tx.Amount
			}
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Fatal(err)
	}

	return balance
}

func VerifyBalance(address string, amount float64, db *leveldb.DB) bool {
	balance := CalculateBalance(address, db)

	if balance >= amount {
		return true
	}
	return false
}

func GetTransactionsByAddress(address string, db *leveldb.DB) ([]Transaction, error) {
	var transactions []Transaction
	iter := db.NewIterator(util.BytesPrefix([]byte("")), nil)
	for iter.Next() {
		var block Block
		err := json.Unmarshal(iter.Value(), &block)
		if err != nil {
			return nil, err
		}
		for _, tx := range block.Transactions {
			if tx.Sender == address || tx.Receiver == address {
				transactions = append(transactions, tx)
			}
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func DisplayTransactions(address string, db *leveldb.DB) error {
	transactions, err := GetTransactionsByAddress(address, db)
	if err != nil {
		return err
	}
	fmt.Printf("\nSender\t|\tReceiver\t|\tAmount\t|\tSignature\n")
	fmt.Printf("--------------------------------------------------------\n")
	for _, tx := range transactions {
		fmt.Printf("%s\t|\t%s\t|\t%f\t|\t%s\n", tx.Sender, tx.Receiver, tx.Amount, tx.Signature)
		fmt.Printf("--------------------------------------------------------\n")
	}
	return nil
}
