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

func Leerfdisk(tokens []string) (string, error) {
	cmd := &FDISK{}
	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-size=\d+|-unit=[kKmM]|-fit=[bBfF]{2}|-path="[^"]+"|-path=[^\s]+|-type=[pPeElL]|-name="[^"]+"|-name=[^\s]+`)
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
	if cmd.unit == "" {
		cmd.unit = "M"
	}
	if cmd.fit == "" {
		cmd.fit = "WF"
	}

	if cmd.typ == "" {
		cmd.typ = "P"
	}

	err := Cfdisk(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}
	result = fmt.Sprintf("Comando fdisk ejecutado con éxito.- Tamaño: %d, Unidad: %s, Ajuste: %s, Ruta: %s, Tipo: %s, Nombre: %s", cmd.size, cmd.unit, cmd.fit, cmd.path, cmd.typ, cmd.name)

	return result, nil
}

func Cfdisk(fdisk *FDISK) error {
	sizeBytes, err := utils.ConvertToBytes(fdisk.size, fdisk.unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

	if fdisk.typ == "P" {
		err = crearParticionP(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición primaria:", err)
			return err
		}
	} else if fdisk.typ == "E" {
		err = crearParticionExtendida(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición extendida:", err)
			return err
		}
	} else if fdisk.typ == "L" {
		err = crearParticionL(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición lógica:", err)
			return err
		}
	} else {
		return errors.New("tipo de partición no válido")
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
	err = mbr.SerializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

func crearParticionExtendida(fdisk *FDISK, sizeBytes int) error {
	var mbr structures.MBR
	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando MBR:", err)
		return err
	}

	// Verifica si ya existe una partición extendida
	if mbr.HasExtendedPartition() {
		return errors.New("ya existe una partición extendida en el disco")
	}

	particionD, startP, indexP := mbr.GetFirstAvailablePartition()
	if particionD == nil {
		return errors.New("no hay particiones disponibles")
	}

	// Crear la partición extendida
	particionD.CrearP(startP, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

	// Guardar la partición en el MBR
	mbr.Mbr_partitions[indexP] = *particionD

	// Serializar el MBR
	err = mbr.SerializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error serializando MBR:", err)
	}
	return nil
}

// CrearParticionL crea una partición lógica en el disco
func crearParticionL(fdisk *FDISK, sizeBytes int) error {
	var mbr structures.MBR
	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando MBR:", err)
		return err
	}

	// Buscar la partición extendida
	extendida, _ := mbr.ParticionPorNombre("E")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Verificar si la partición extendida tiene espacio disponible
	if int(extendida.PartSize) < sizeBytes {
		return errors.New("no hay suficiente espacio en la partición extendida")
	}

	// Crear la partición lógica
	particionD, startP, indexP := mbr.GetFirstAvailablePartition()
	if particionD == nil {
		return errors.New("no hay particiones disponibles")
	}

	// Crear la partición lógica
	particionD.CrearP(startP, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

	// Guardar la partición en el MBR
	mbr.Mbr_partitions[indexP] = *particionD

	// Serializar el MBR
	err = mbr.SerializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error serializando MBR:", err)
	}
	return nil
}

func PrintPartitions(mbr *structures.MBR) {
	for i, partition := range mbr.Mbr_partitions {
		fmt.Printf("Partición %d:\n", i+1)
		fmt.Printf("  Nombre: %s\n", strings.Trim(string(partition.PartName[:]), "\x00 "))
		fmt.Printf("  Tipo: %s\n", strings.Trim(string(partition.PartType[:]), "\x00 "))
		fmt.Printf("  Ajuste: %s\n", strings.Trim(string(partition.PartFit[:]), "\x00 "))
		fmt.Printf("  Inicio: %d\n", partition.PartStart)
		fmt.Printf("  Tamaño: %d\n", partition.PartSize)
	}
}
