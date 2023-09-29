package logic

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tyler-smith/go-bip39"
)

type Account struct {
	Address   string
	PublicKey string
}

func CreateAccount() {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		log.Fatal(err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mnemonic: ", mnemonic)

	seed := bip39.NewSeed(mnemonic, "")
	privateKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		fmt.Println("Hubo un error en la creación de las llaves ", err)
	} else {
		address := crypto.PubkeyToAddress(privateKey.PublicKey)

		publicKeyHex := fmt.Sprintf("%x%x", privateKey.PublicKey.X, privateKey.PublicKey.Y)
		addressHex := fmt.Sprintf("%x", address)

		err := SaveAccountToDB(addressHex, publicKeyHex)
		if err != nil {
			fmt.Println("Error al guardar la cuenta en la base de datos: ", err)
		}

		fmt.Printf("Private Key: %x\n", privateKey.D)
		fmt.Printf("Public Key: %s\n", publicKeyHex)
		fmt.Printf("Address: %s\n", addressHex)

	}
}

func SaveAccountToDB(address, publicKey string) error {
	account := Account{
		Address:   address,
		PublicKey: publicKey,
	}

	accountData, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("Error al pasar a json: %v", err)
	}

	db, err := leveldb.OpenFile("./accounts.db", nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Put([]byte(address), accountData, nil)
	if err != nil {
		return fmt.Errorf("No se pudo almacenar la cuenta en la base de datos: %v", err)
	}

	return nil
}

func ShowAllAccounts() error {

	db, err := leveldb.OpenFile("./accounts.db", nil)
	if err != nil {
		return err
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		var account Account
		err := json.Unmarshal(value, &account)
		if err != nil {
			return err
		}

		fmt.Printf("Address: %s\n", key)
		fmt.Printf("Public Key: %s\n", account.PublicKey)
		fmt.Println("-------------------------")
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		return err
	}

	return nil
}
