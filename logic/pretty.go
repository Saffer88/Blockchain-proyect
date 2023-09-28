package logic

import (
	"encoding/json"
	"fmt"
)

func Pretty(block Block) {

	blockJson, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error al convertir el bloque a JSON:", err)
		return
	}
	fmt.Println(string(blockJson))

}
