package comands

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	structures "github.com/emilio-riv20/proyecto1/Codigo/Structures"
	utils "github.com/emilio-riv20/proyecto1/Codigo/utils"
)

type MOUNT struct {
	path string
	name string
}

func Mount(tokens []string) (string, error) {
	cmd := &MOUNT{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+|-name="[^"]+"|-name=[^\s]+`)
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
				return "", errors.New("el path debe tener un argumento")
			}
			cmd.path = value
		case "-name":
			if value == "" {
				return "", errors.New("debe tener un nombre")
			}
			cmd.name = value
		default:
			return "", fmt.Errorf("parametro desconocido: %s", key)
		}
	}

	if cmd.path == "" || cmd.name == "" {
		return "", errors.New("faltan parametros requeridos")
	}

	err := cMount(cmd)
	if err != nil {
		return "", fmt.Errorf("error al montar la partición: %v", err)
	}

	result := fmt.Sprintf("Comando mount ejecutado con éxito.- Ruta: %s, Nombre: %s", cmd.path, cmd.name)
	return result, nil
}

func GenerarId(mount *MOUNT, indexPartition int) (string, error) {
	// Asignar letra
	letter, err := utils.Letra(mount.path)
	if err != nil {
		fmt.Println("Error obteniendo la letra:", err)
		return "", err
	}

	// Crear id
	idPartition := fmt.Sprintf("%s%d%s", utils.Carnet, indexPartition+1, letter)

	return idPartition, nil
}

func cMount(mount *MOUNT) error {

	var mbr structures.MBR
	err := mbr.DeserializeMBR(mount.path)

	if err != nil {
		return err
	}

	particion, indexP := mbr.ParticionPorNombre(mount.name)
	if particion == nil {
		fmt.Println("No se encontro la particion")
		return errors.New("no se encontro la particion")
	}

	idP, err := GenerarId(mount, indexP)
	if err != nil {
		fmt.Println("Error generando el id:", err)
		return err
	}

	utils.MountedPartitions[idP] = mount.path
	particion.MontarP(indexP, idP)
	mbr.Mbr_partitions[indexP] = *particion

	err = mbr.Serializar(mount.path)
	if err != nil {
		fmt.Println("Error serializando el mbr:", err)
		return err
	}

	return nil
}
