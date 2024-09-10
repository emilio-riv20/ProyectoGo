package comands

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// RMDISK estructura que representa el comando rmdisk con su parámetro
type RMDISK struct {
	path string // Ruta del archivo del disco
}

// ParserRmdisk parsea el comando rmdisk y devuelve un mensaje detallado y un error
func CommandRmdisk(tokens []string) (string, error) {
	cmd := &RMDISK{} // Crea una nueva instancia de RMDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar el parámetro del comando rmdisk
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Elimina las comillas del valor si están presentes
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar el parámetro -path
		switch key {
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -path haya sido proporcionado
	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}

	// Lógica para ejecutar el comando RMDISK
	err := cRm(cmd)
	if err != nil {
		return "", fmt.Errorf("error al eliminar el disco: %v", err)
	}

	// Construye un mensaje detallado con las especificaciones del comando ejecutado
	result := fmt.Sprintf("Comando rmdisk ejecutado con éxito.- Ruta: %s", cmd.path)

	return result, nil // Devuelve el mensaje detallado
}

func cRm(rmdisk *RMDISK) error {
	// Intenta eliminar el archivo en la ruta especificada
	print()
	err := os.Remove(rmdisk.path)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el archivo: %v", err)
	}

	return nil
}
