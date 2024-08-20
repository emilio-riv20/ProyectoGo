package Analyzer

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Analyzer(input string) (interface{}, error) {
	tokens := strings.Fields(input)

	if len(tokens) == 0 {
		return nil, errors.New("no se proporcionó ningún comando")
	}

	switch tokens[0] {
	case "mkdisk":
		fmt.Println("Funciona 1")
	case "rmdisk":
		fmt.Println("Funciona 2")
	case "fdisk":
		// Implementa la lógica aquí
	case "mount":
		// Implementa la lógica aquí
	case "rep":
		// Implementa la lógica aquí
	case "clear":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return nil, errors.New("no se pudo limpiar la terminal")
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("comando desconocido: %s", tokens[0])
	}
	return nil, nil
}
