package main

import (
	"Blockchain-proyect/logic"
	"fmt"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	db, err := leveldb.OpenFile("./level.db", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	// Verifica si la base de datos está vacía
	if !iter.First() { // Si está vacía agregamos el bloque genesis
		fmt.Println("\nCreando bloque Genesis.")
		genesis := logic.Genesis()
		err = logic.SaveBlockToDB(genesis, db)
		if err != nil {
			panic(err)
		}
		logic.Pretty(genesis)
		fmt.Print("\n\nEnter... ")
		var wait int
		fmt.Scanln(&wait)
	}

	for { // despliegue de las opciones
		{
			fmt.Print("\033[H\033[2J")
			fmt.Println(`
            ~~~~ Bienvenido al sistema blockchain ~~~~
            1. Crear cuenta.
            2. Ver todas las cuentas existentes.
            3. Escribir una transacción.
            4. Búsqueda de bloque específico usando ID.
            5. Búsqueda de bloque usando transacción.
            6. DEBUG Crear un bloque vacío.
            7. DEBUG Consultar los bloques totales.
            8. DEBUG Volcar todos los bloques.
            ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
            `)
			fmt.Print("\nOpción: ")

			var opcion string
			fmt.Scanln(&opcion)

			switch opcion {

			case "1":

			case "2":

			case "3":
				fmt.Print("\033[H\033[2J")
				fmt.Println(`
				~~~~ Escribir una transacción ~~~~~
				     ~~~~~~~~~~~~~~~~~~~~~~
                `)

				fmt.Print("\ningrese su dirección: ")
				var user string
				fmt.Scanln(&user)
				fmt.Print("\ningrese la dirección del destinatario: ")
				var receiver string
				fmt.Scanln(&receiver)
				fmt.Print("\ningrese el monto: ")
				var mount float64
				fmt.Scanln(&mount)

				newTransaction := logic.NewTransaction(user, receiver, mount)

				blockIsFull := logic.Limit(db)

				if blockIsFull {
					fmt.Println("\nEl último bloque está lleno, así que se creará un nuevo bloque")

					fmt.Print("\nIndique el límite de transacciones del bloque: ")
					var limit int
					fmt.Scanln(&limit)

					transactions := []logic.Transaction{newTransaction}
					newblock := logic.GenerateBlock(db, transactions, limit)
					err = logic.SaveBlockToDB(newblock, db)
					if err != nil {
						panic(err)
					}
					fmt.Println("\nSe ha creado el bloque correctamente")
					logic.Pretty(newblock)

					fmt.Print("\n\nEnter... ")
					var wait int
					fmt.Scanln(&wait)

				} else {
					lastBlock, err := logic.GetLastBlock(db)
					if err != nil {
						fmt.Println("\nError al obtener el último bloque", err)
						panic(err)
					}

					logic.AddTransaction(&lastBlock, newTransaction)

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

					fmt.Print("\n\nEnter... ")
					var wait int
					fmt.Scanln(&wait)

				}
			case "4":
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
					panic(err)
				}

				fmt.Print("\nEl bloque con el ID correspondiente es el siguiente:\n\n")
				logic.Pretty(searchblock)

				fmt.Print("\n\nPress Enter... ")
				var wait int
				fmt.Scanln(&wait)
			case "5":

			case "6":

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

			case "7":

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

			case "8":

			default:
				fmt.Println("Opción no válida. Intente de nuevo.")
				time.Sleep(1 * time.Second)
			}
		}
	}

}
