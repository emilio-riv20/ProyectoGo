package comands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	structures "github.com/emilio-riv20/proyecto1/Codigo/Structures"
	utils "github.com/emilio-riv20/proyecto1/Codigo/utils"
)

type FDISK struct {
	size int
	unit string
	path string
	fit  string
	typ  string
	name string
}

func Command_fdisk(tokens []string) (string, error) {
	cmd := &FDISK{}
	args := strings.Join(tokens, "")
	re := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-path="[^"]+"|-path=[^\s]+|-fit=[bBfFwW]{2}|-type=[pPeE]|-name="[^"]+"|-name=[^\s]+`)
	var result string

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
		case "-size":
			// Convierte el valor del tamaño a un entero
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return "", errors.New("el tamaño debe ser un número entero positivo")
			}
			cmd.size = size
		case "-unit":
			// Verifica que la unidad sea "K" o "M"
			if value != "K" && value != "M" {
				return "", errors.New("la unidad debe ser K o M")
			}
			cmd.unit = strings.ToUpper(value)
		case "-fit":
			// Verifica que el ajuste sea "BF", "FF" o "WF"
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return "", errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.fit = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-type":
			// Verifica que el tipo sea "P", "E" o "L"
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return "", errors.New("el tipo debe ser P, E o L")
			}
			cmd.typ = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return "", errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	if cmd.size == 0 {
		return "", errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return "", errors.New("faltan parámetros requeridos: -name")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if cmd.unit == "" {
		cmd.unit = "M"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if cmd.fit == "" {
		cmd.fit = "WF"
	}

	// Si no se proporcionó el tipo, se establece por defecto a "P"
	if cmd.typ == "" {
		cmd.typ = "P"
	}

	// Crear la partición con los parámetros proporcionados
	err := fdis(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// Construye un mensaje detallado con las especificaciones del comando ejecutado
	result = fmt.Sprintf("Comando fdisk ejecutado con éxito.- Tamaño: %d, Unidad: %s, Ajuste: %s, Ruta: %s, Tipo: %s, Nombre: %s", cmd.size, cmd.unit, cmd.fit, cmd.path, cmd.typ, cmd.name)

	return result, nil
}

func fdis(fdisk *FDISK) error {
	sizeBytes, err := utils.ConvertToBytes(fdisk.size, fdisk.unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

	if fdisk.typ == "P" {
		// Crear partición primaria
		err = crearParticionP(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición primaria:", err)
			return err
		}
	} else if fdisk.typ == "E" {
		fmt.Println("Creando partición extendida...") // Les toca a ustedes implementar la partición extendida
	} else if fdisk.typ == "L" {
		fmt.Println("Creando partición lógica...") // Les toca a ustedes implementar la partición lógica
	}

	return nil
}

func crearParticionP(fdisk *FDISK, sizeBytes int) error {
	var mbr structures.MBR
	err := mbr.DeserializeMBR(fdisk.path)

	if err != nil {
		fmt.Println("Error deserializando MBR:", err)
		return err
	}

	particionD, startP, indexP := mbr.GetFirstAvailablePartition()
	if particionD == nil {
		fmt.Println("No hay particiones disponibles.")
	}

	// Crear la partición con los parámetros proporcionados
	particionD.CrearP(startP, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

	// Colocar la partición en el MBR
	if particionD != nil {
		mbr.Mbr_partitions[indexP] = *particionD
	}

	// Serializar el MBR en el archivo binario
	err = mbr.Serializar(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}
