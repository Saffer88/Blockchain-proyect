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

// Para pedir input de la llave privada en oculto por la consola
func Hide_private_key() string {
	fmt.Print("Ingrese la llave privada del sender: ")

	// Usa la función ReadPassword del paquete terminal para leer la contraseña de forma segura.
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	// Convierte la contraseña de bytes a un string.
	password := string(bytePassword)

	return password
}

func GetPublicKeyForUser(address string) (string, error) {
	// Abre la base de datos de cuentas
	db, err := leveldb.OpenFile("./accounts.db", nil)
	if err != nil {
		return "no se abrio correctamente la db", err
	}
	defer db.Close()

	// Busca el usuario por su dirección
	data, err := db.Get([]byte(address), nil)
	if err != nil {
		return "wateo en el get", err
	}

	// Decodifica los datos en una estructura de cuenta
	var account Account
	err = json.Unmarshal(data, &account)
	if err != nil {
		return "pasar a json falló", err
	}

	// Devuelve la clave pública del usuario
	return account.PublicKey, nil
}
