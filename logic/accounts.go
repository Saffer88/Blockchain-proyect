package logic

import (
	"encoding/json"
	"fmt"
	"log"

	"syscall"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ssh/terminal"
)

type Account struct {
	Address   string
	PublicKey string
	Balance   float64
}

type AccountJSON struct {
	Mode      int    `json:"Mode"`
	PublicKey string `json:"public key"`
	Address   string `json:"Address"`
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

	seed := bip39.NewSeed(mnemonic, "")
	privateKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		fmt.Println("Hubo un error en la creación de las llaves ", err)
	} else {
		address := crypto.PubkeyToAddress(privateKey.PublicKey)

		publicKeyHex := fmt.Sprintf("%x%x", privateKey.PublicKey.X, privateKey.PublicKey.Y)
		addressHex := fmt.Sprintf("%x", address)

		err := SaveAccountToDB(addressHex, publicKeyHex, 1000)
		if err != nil {
			fmt.Println("Error al guardar la cuenta en la base de datos: ", err)
		}

		fmt.Printf("Private Key: %x\n", privateKey.D)
		fmt.Printf("Public Key: %s\n", publicKeyHex)
		fmt.Printf("Address: %s\n", addressHex)

		Accountjson := AccountJSON{
			Mode:      1,
			PublicKey: publicKeyHex,
			Address:   addressHex,
		}
		jsonBytes, err := json.Marshal(Accountjson)
		if err != nil {
			fmt.Println("Error al convertir la estructura en JSON:", err)
			return
		}
		jsonStr := string(jsonBytes)
		Broadcast(jsonStr)
	}
}

func SaveAccountToDB(address, publicKey string, balance float64) error {
	account := Account{
		Address:   address,
		PublicKey: publicKey,
		Balance:   balance,
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

func UpdateBalance(address string, newBalance float64) error {
	db, err := leveldb.OpenFile("./accounts.db", nil)
	if err != nil {
		return err
	}
	defer db.Close()

	data, err := db.Get([]byte(address), nil)
	if err != nil {
		return err
	}

	var account Account
	err = json.Unmarshal(data, &account)
	if err != nil {
		return err
	}

	account.Balance += newBalance

	updatedData, err := json.Marshal(account)
	if err != nil {
		return err
	}

	err = db.Put([]byte(address), updatedData, nil)
	if err != nil {
		return err
	}

	return nil
}

func GenesisAccount() error {
	account := Account{
		Address:   "0e1dd7f2e5cb568ee13534424aaa978e484df040",
		PublicKey: "7d6866b740b19acdef6055398dfb2ace996153099471f12f6f1ff19d7856157ddf49cd8bca60f22467dbac63df248e06adbde6246d0c77385b3dd592e3ae31a1",
		Balance:   1000000,
	}

	err := SaveAccountToDB(account.Address, account.PublicKey, account.Balance)

	if err != nil {
		return err
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
		fmt.Printf("Public Key: %.5f\n", account.Balance)
		fmt.Println("-------------------------")
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		return err
	}

	return nil
}

func HidePrivateKey() string {
	fmt.Print("\nIngrese la llave privada del sender: ")

	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	password := string(bytePassword)

	return password
}

func GetPublicKeyForUser(address string) (string, error) {
	db, err := leveldb.OpenFile("./accounts.db", nil)
	if err != nil {
		return "error abriendo db", err
	}
	defer db.Close()

	accountData, err := db.Get([]byte(address), nil)
	if err != nil {
		return "error geteando la cuenta", err
	}

	var account Account
	err = json.Unmarshal(accountData, &account)
	if err != nil {
		return "error json-codeando", err
	}

	return account.PublicKey, nil
}

func VerifyAccount(adress string) (bool, error) {

	db, err := leveldb.OpenFile("./accounts.db", nil)
	if err != nil {
		return false, err
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		var account Account
		err := json.Unmarshal(value, &account)
		if err != nil {
			return false, err
		}

		if string(key) == adress {
			return true, nil
		}
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		return false, err
	}
	return false, nil
}
