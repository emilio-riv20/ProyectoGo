package comands

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type RMDISK struct {
	path string
}

func CommandRmdisk(tokens []string) (string, error) {
	cmd := &RMDISK{}
	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+`)
	matches := re.FindAllString(args, -1)
	for _, match := range matches {
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
		case "-path":
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		default:

			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}
	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}

	err := cRm(cmd)
	if err != nil {
		return "", fmt.Errorf("error al eliminar el disco: %v", err)
	}

	result := fmt.Sprintf("Comando rmdisk ejecutado con éxito.- Ruta: %s", cmd.path)

	return result, nil // Devuelve el mensaje detallado
}

func cRm(rmdisk *RMDISK) error {
	print()
	err := os.Remove(rmdisk.path)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el archivo: %v", err)
	}

	return nil
}
