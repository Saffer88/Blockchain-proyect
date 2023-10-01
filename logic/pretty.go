package logic

import (
	"encoding/json"
	"fmt"
)

/*
//Colorines
reset := "\033[0m"
red := "\033[31m"
green := "\033[32m"
yellow := "\033[33m"
blue := "\033[34m"
/*USO:
fmt.Println(red + "Este es un mensaje en rojo." + reset)
fmt.Println(green + "Este es un mensaje en verde." + reset)
fmt.Println(yellow + "Este es un mensaje en amarillo." + reset)
fmt.Println(blue + "Este es un mensaje en azul." + reset)
*/

func Pretty(block Block) {

	blockJson, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error al convertir el bloque a JSON:", err)
		return
	}
	fmt.Println(string(blockJson))

}
