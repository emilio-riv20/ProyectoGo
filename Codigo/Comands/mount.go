package comands

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	structures "github.com/emilio-riv20/proyecto1/Codigo/Structures"
	utils "github.com/emilio-riv20/proyecto1/Codigo/utils"
)

// MOUNT estructura que representa el comando mount con sus parámetros
type MOUNT struct {
	path string
	name string
}

// CommandMount parsea el comando mount y devuelve una instancia de MOUNT
func LeerMount(tokens []string) (string, error) {
	cmd := &MOUNT{} // Crea una nueva instancia de MOUNT

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mount
	re := regexp.MustCompile(`(?i)-path="[^"]+"|-path=[^\s]+|-name="[^"]+"|-name=[^\s]+`)
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

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
		case "-path":
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-name":
			if value == "" {
				return "", errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return "", errors.New("faltan parámetros requeridos: -name")
	}

	// Montamos la partición
	err := cMount(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// Generar el ID de partición para el resultado final
	idPartition, err := GenerateIdPartition(cmd, 0)
	if err != nil {
		return "", err
	}

	// Devuelve el nombre de la partición junto con su ID
	result := fmt.Sprintf("Partición: %s montada correctamente, con ID: %s", cmd.name, idPartition)
	return result, nil
}

func cMount(mount *MOUNT) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(mount.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Buscar la partición con el nombre especificado
	partition, indexPartition := mbr.ParticionPorNombre(mount.name)
	if partition == nil {
		return errors.New("error: la partición no existe")
	}

	// Generar un id único para la partición
	idPartition, err := GenerateIdPartition(mount, indexPartition)
	if err != nil {
		fmt.Println("Error generando el id de partición:", err)
		return err
	}

	// Guardar la partición montada en la lista de montajes utilses
	utils.MountedPartitions[idPartition] = mount.path

	// Modificamos la partición para indicar que está montada
	partition.MontarP(indexPartition, idPartition)

	// Guardar la partición modificada en el MBR
	mbr.Mbr_partitions[indexPartition] = *partition

	// Serializar la estructura MBR en el archivo binario
	err = mbr.SerializeMBR(mount.path)
	if err != nil {
		fmt.Println("Error serializando el MBR:", err)
		return err
	}

	return nil
}

// GenerateIdPartition crea un ID único para una partición montada
func GenerateIdPartition(mount *MOUNT, indexPartition int) (string, error) {
	letter, err := utils.GetLetter(mount.path)
	if err != nil {
		fmt.Println("Error obteniendo la letra:", err)
		return "", err
	}

	existingPartitionCount := 0
	for path := range utils.MountedPartitions {
		if path == mount.path {
			existingPartitionCount++
		}
	}

	idPartition := fmt.Sprintf("%s%d%s", utils.Carnet, existingPartitionCount+1, letter)

	// Verifica si el ID ya existe
	for id := range utils.MountedPartitions {
		if id == idPartition {
			return "", fmt.Errorf("ID de partición ya existe: %s", idPartition)
		}
	}

	return idPartition, nil
}

func CommandListMounts() string {
	if len(utils.MountedPartitions) == 0 {
		return "No hay particiones montadas."
	}

	var result strings.Builder
	result.WriteString("Particiones montadas:\n")

	for id, path := range utils.MountedPartitions {
		result.WriteString(fmt.Sprintf("ID: %s, Ruta: %s\n", id, path))
	}

	return result.String()
}
