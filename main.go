package main

import (
	"Blockchain-proyect/logic"
	//"encoding/hex"
	"fmt"
	//"log"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	logic.Def_handler()

	db, err := leveldb.OpenFile("./blocks.db", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	logic.CreateGenesis(db) // Evaluamos si es necesario crear el bloque genesis
	for {                   // despliegue de las opciones
		{
			fmt.Print("\033[H\033[2J")
			fmt.Println(`
            ~~~~ Bienvenido al sistema blockchain ~~~~
            1. Crear cuenta.
            2. Consultar el saldo de una cuenta.
            3. Escribir una transacción.
            4. Verificar la firma de una transacción.
            5. Ver todas las transacciones de una dirección.
            6. Búsqueda de bloque específico usando ID.
            7. DEBUG Crear un bloque vacío (Limit 3 Default).
            8. DEBUG Consultar los bloques totales.
            9. DEBUG Ver todas las cuentas existentes.
            ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
            `)
			fmt.Print("\nOpción: ")

			var opcion string
			fmt.Scanln(&opcion)

			switch opcion {

			case "1":
				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Crear una cuenta ~~~~
				  ~~~~~~~~~~~~~~~~~~~~~~
                `)

				logic.CreateAccount() // Creamos cuenta
				logic.PressEnter()

			case "2":

				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Consultar el saldo de una cuenta ~~~~~
				         ~~~~~~~~~~~~~~~~~~~~~~
                `)
				fmt.Print("\ningrese la dirección: ")
				var address string
				fmt.Scanln(&address)

				balance, err := logic.GetBalance(address)

				if err != nil {
					fmt.Print("\nError al obtener el balance")
					logic.PressEnter()
					break
				}

				fmt.Printf("\nEl saldo de la dirección %s es: %.5f\n", address, balance)

				logic.PressEnter()

			case "3":

				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Escribir una transacción ~~~~~
				     ~~~~~~~~~~~~~~~~~~~~~~
                `)

				fmt.Print("\ningrese la dirección del sender: ")
				var address string
				fmt.Scanln(&address)

				fmt.Print("\ningrese la dirección del destinatario: ")
				var receiver string
				fmt.Scanln(&receiver)

				fmt.Print("\ningrese el monto: ")
				var amount float64
				fmt.Scanln(&amount)

				var priv string = logic.HidePrivateKey()

				fmt.Print("\n\ningrese el Nonce: ")
				var nonce int
				fmt.Scanln(&nonce)

				newTransaction, err := logic.NewTransaction(address, receiver, amount, priv, db, nonce)

				if newTransaction == nil {
					logic.PressEnter()
					break
				}

				blockIsFull := logic.Limit(db)

				if blockIsFull {
					fmt.Println("\nEl último bloque está lleno, así que se creará un nuevo bloque")

					fmt.Print("\nIndique el límite de transacciones del bloque: ")
					var limit int
					fmt.Scanln(&limit)

					transactions := []logic.Transaction{*newTransaction}
					newblock := logic.GenerateBlock(db, transactions, limit)
					err = logic.SaveBlockToDB(newblock, db)
					if err != nil {
						panic(err)
					}
					err = logic.UpdateBalance(address, -amount)
					if err != nil {
						panic(err)
					}
					err = logic.UpdateBalance(receiver, amount)
					if err != nil {
						panic(err)
					}
					fmt.Println("\nSe ha creado el bloque correctamente")
					logic.Pretty(newblock)

					logic.PressEnter()

				} else {
					lastBlock, err := logic.GetLastBlock(db)
					if err != nil {
						fmt.Println("\nError al obtener el último bloque", err)
						panic(err)
					}

					logic.AddTransaction(&lastBlock, *newTransaction)
					logic.UpdateBlockHash(&lastBlock) // actualizamos el hash con las nuevas transacciones
					err = logic.SaveBlockToDB(lastBlock, db)
					if err != nil {
						panic(err)
					}
					err = logic.UpdateBalance(address, -amount)
					if err != nil {
						panic(err)
					}
					err = logic.UpdateBalance(receiver, amount)
					if err != nil {
						panic(err)
					}

					fmt.Println("\nSe ha escrito la información en el bloque:")

					searchblock, err := logic.GetBlockFromDB(lastBlock.Index, db)
					if err != nil {
						panic(err)
					}

					logic.Pretty(searchblock)

					logic.PressEnter()

				}
			case "4":

				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Verificar firma ~~~~
				    ~~~~~~~~~~~~~~~
                `)
				fmt.Print("\ningrese dirección del sender: ")
				var address string
				fmt.Scanln(&address)

				pub, err := logic.GetPublicKeyForUser(address)
				if err != nil {
					fmt.Println("\nLa cuenta no existe: ")
					logic.PressEnter()
					break
				}

				fmt.Print("\ningrese la dirección del destinatario: ")
				var receiver string
				fmt.Scanln(&receiver)

				fmt.Print("\ningrese el monto: ")
				var mount float64
				fmt.Scanln(&mount)

				fmt.Print("\ningrese la firma: ")
				var signature string
				fmt.Scanln(&signature)

				logic.VerifyTransaction(address, receiver, mount, signature, pub)

				logic.PressEnter()

			case "5":
				fmt.Print("\033[H\033[2J")
				fmt.Println(`
                 ~~~~ Obtener todas las transacciones de una cuenta ~~~~
                     ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
				`)
				fmt.Print("\nIngrese la dirección: ")
				var address string
				fmt.Scanln(&address)
				logic.DisplayTransactions(address, db)

				logic.PressEnter()

			case "6":
				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Leer un Bloque ~~~~
				    ~~~~~~~~~~~~~~~
                `)
				fmt.Print("\ningrese el ID del bloque: ")
				var ID int
				fmt.Scanln(&ID)

				searchblock, err := logic.GetBlockFromDB(ID, db)
				if err != nil {
					fmt.Print("\nEl bloque no existe.")
					logic.PressEnter()
				} else {
					fmt.Print("\nEl bloque con el ID correspondiente es el siguiente:\n\n")
					logic.Pretty(searchblock)

					logic.PressEnter()
				}

			case "7":

				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Agregando nuevo bloque ~~~~~~
				     ~~~~~~~~~~~~~~~~~~~~~~
                `)

				newblock := logic.GenerateBlock(db, nil, 3)

				logic.Pretty(newblock)

				err = logic.SaveBlockToDB(newblock, db)
				if err != nil {
					panic(err)
				}

				fmt.Println("\n\nRegistro completado.")
				logic.PressEnter()

			case "8":

				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Consultar bloques totales ~~~~~
				    ~~~~~~~~~~~~~~~~~~~~~~~~~~
                `)

				totalBlocks, err := logic.CountTotalBlocks(db)
				if err != nil {
					fmt.Println("Error al contar los bloques:", err)
					return
				}

				fmt.Printf("Cantidad total de bloques: %d\n", totalBlocks)

				logic.PressEnter()

			case "9":
				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Consultar cuentas totales ~~~~~
				    ~~~~~~~~~~~~~~~~~~~~~~~~~~
                `)
				logic.ShowAllAccounts()

				logic.PressEnter()

			default:
				fmt.Println("Opción no válida. Intente de nuevo.")
				time.Sleep(1 * time.Second)
			}
		}
	}

}
