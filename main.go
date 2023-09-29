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
	db, err := leveldb.OpenFile("./level.db", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	logic.CreateGenesis(db) // Evaluamos si es necesario crear el bloque genesis

	for { // despliegue de las opciones
		{
			fmt.Print("\033[H\033[2J")
			fmt.Println(`
            ~~~~ Bienvenido al sistema blockchain ~~~~
            1. Crear cuenta.
            2. Consultar el saldo de una cuenta.
            3. Escribir una transacción.
            4. Verificar la firma de una transacción.
            5. Búsqueda de bloque específico usando ID.
            6. Búsqueda de bloque usando transacción.
            7. DEBUG Crear un bloque vacío (Limit 3 Default).
            8. DEBUG Consultar los bloques totales.
            9. DEBUG Ver todas las cuentas existentes.
            10. DEBUG Volcar todos los bloques.
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

				logic.CreateAccount()

				fmt.Print("\n\nPress Enter... ")
				var wait int
				fmt.Scanln(&wait)
			case "2":

			case "3":
				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Escribir una transacción ~~~~~
				     ~~~~~~~~~~~~~~~~~~~~~~
                `)

				fmt.Print("\ningrese su dirección: ")
				var address string
				fmt.Scanln(&address)
				fmt.Print("\ningrese la dirección del destinatario: ")
				var receiver string
				fmt.Scanln(&receiver)
				fmt.Print("\ningrese el monto: ")
				var mount float64
				fmt.Scanln(&mount)
				fmt.Print("\ningrese su llave privada: ")
				var priv string
				fmt.Scanln(&priv)

				newTransaction, err := logic.NewTransaction(address, receiver, mount, priv)

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
					fmt.Println("\nSe ha creado el bloque correctamente")
					logic.Pretty(newblock)

					fmt.Print("\n\nPress Enter... ")
					var wait int
					fmt.Scanln(&wait)

				} else {
					lastBlock, err := logic.GetLastBlock(db)
					if err != nil {
						fmt.Println("\nError al obtener el último bloque", err)
						panic(err)
					}

					logic.AddTransaction(&lastBlock, *newTransaction)

					err = logic.SaveBlockToDB(lastBlock, db)
					if err != nil {
						panic(err)
					}

					fmt.Println("\nSe ha escrito la información en el bloque:")

					searchblock, err := logic.GetBlockFromDB(lastBlock.Index, db)
					if err != nil {
						panic(err)
					}

					logic.Pretty(searchblock)

					fmt.Print("\n\nPress Enter... ")
					var wait int
					fmt.Scanln(&wait)

				}
			case "4":

				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Verificar firma ~~~~
				    ~~~~~~~~~~~~~~~
                `)

				fmt.Print("\ningrese su dirección: ")
				var address string
				fmt.Scanln(&address)
				fmt.Print("\ningrese la dirección del destinatario: ")
				var receiver string
				fmt.Scanln(&receiver)
				fmt.Print("\ningrese el monto: ")
				var mount float64
				fmt.Scanln(&mount)
				fmt.Print("\ningrese la firma: ")
				var signature string
				fmt.Scanln(&signature)
				fmt.Print("\ningrese su llave pública: ")
				var pub string
				fmt.Scanln(&pub)

				verify := logic.VerifyTransaction(address, receiver, mount, signature, pub)
				if verify {
					fmt.Print("\nLa firma es válida ")
					fmt.Print("\n\nPress Enter... ")
					var wait int
					fmt.Scanln(&wait)
				} else {

					fmt.Print("\nLa firma no es válida ")
					fmt.Print("\n\nPress Enter... ")
					var wait int
					fmt.Scanln(&wait)
				}

			case "5":
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
					fmt.Print("\n\nPress Enter... ")
					var wait int
					fmt.Scanln(&wait)
				} else {
					fmt.Print("\nEl bloque con el ID correspondiente es el siguiente:\n\n")
					logic.Pretty(searchblock)

					fmt.Print("\n\nPress Enter... ")
					var wait int
					fmt.Scanln(&wait)
				}

			case "6":

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

				time.Sleep(1 * time.Second)
				fmt.Println("\n\nRegistro completado.")
				fmt.Print("\n\nPress Enter... ")
				var wait int
				fmt.Scanln(&wait)

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

				fmt.Print("\n\nPress Enter... ")
				var wait int
				fmt.Scanln(&wait)

			case "9":
				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Consultar cuentas totales ~~~~~
				    ~~~~~~~~~~~~~~~~~~~~~~~~~~
                `)
				logic.ShowAllAccounts()

				fmt.Print("\n\nPress Enter... ")
				var wait int
				fmt.Scanln(&wait)

			case "10":

			default:
				fmt.Println("Opción no válida. Intente de nuevo.")
				time.Sleep(1 * time.Second)
			}
		}
	}

}
