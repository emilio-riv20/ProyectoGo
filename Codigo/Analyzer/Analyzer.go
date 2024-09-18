package Analyzer

import (
	"errors"
	"fmt"
	"strings"

	Comands "github.com/emilio-riv20/proyecto1/Codigo/Comands"
)

func Analyzer(input string) (interface{}, error) {
	lines := strings.Split(input, "\n")
	var resultados []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			resultados = append(resultados, line)
		} else {
			tokens := strings.Fields(line)
			if len(tokens) == 0 {
				continue
			}

			switch tokens[0] {
			case "mkdisk":
				result, err := Comands.Cmkdisk(tokens[1:])
				if err != nil {
					resultados = append(resultados, fmt.Sprintf("Error en el comando mkdisk: %s", err))
				} else {
					resultados = append(resultados, fmt.Sprintf("%v", result))
				}
			case "rmdisk":
				result, err := Comands.CommandRmdisk(tokens[1:])
				if err != nil {
					resultados = append(resultados, fmt.Sprintf("Error en el comando rmdisk: %s", err))
				} else {
					resultados = append(resultados, result)
				}
			case "fdisk":
				result, err := Comands.Leerfdisk(tokens[1:])
				if err != nil {
					resultados = append(resultados, fmt.Sprintf("Error en el comando fdisk: %s", err))
				} else {
					resultados = append(resultados, result)
				}
			case "mount":
				result, err := Comands.LeerMount(tokens[1:])
				if err != nil {
					resultados = append(resultados, fmt.Sprintf("Error en el comando mount: %s", err))
				} else {
					resultados = append(resultados, result)
				}

			case "login":
				result, err := Comands.CLogin(tokens[1:])
				if err != nil {
					resultados = append(resultados, fmt.Sprintf("Error en el comando login: %s", err))
				} else {
					resultados = append(resultados, result)
				}

			case "logout":
				result, err := Comands.Clogout()
				if err != nil {
					resultados = append(resultados, fmt.Sprintf("Error en el comando logout: %s", err))
				} else {
					resultados = append(resultados, result)
				}

			case "mkfs":

			default:
				resultados = append(resultados, fmt.Sprintf("Comando desconocido: %s", tokens[0]))
			}

		}
	}

	if len(resultados) == 0 {
		return nil, errors.New("sin comando o comentario")
	}

	return resultados, nil
}
