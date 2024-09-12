package comands

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	structures "github.com/emilio-riv20/proyecto1/Codigo/Structures"
	utils "github.com/emilio-riv20/proyecto1/Codigo/utils"
)

type MKDISK struct {
	size int
	unit string
	fit  string
	path string
}

func Command_mkdisk(tokens []string) (string, error) {
	cmd := &MKDISK{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-fit=[bBfFwW]{2}|-path="[^"]+"|-path=[^\s]+`)
	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", errors.New("error al parsear los argumentos")
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = value[1 : len(value)-1]

		}
		switch key {
		case "-size":
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return "", errors.New("el tamaño del disco debe ser un número entero positivo")
			}
			cmd.size = size
		case "-unit":
			if value != "K" && value != "M" && value != "m" && value != "k" {
				return "", errors.New("la unidad de medida debe ser K o M")
			}
			cmd.fit = value
		case "-path":
			if value == "" {
				return "", errors.New("el path es obligatorio")
			}
			cmd.path = value
		default:
			return "", errors.New("argumento no reconocido")
		}
	}
	if cmd.size == 0 {
		return "", errors.New("el tamaño del disco es obligatorio")
	}
	if cmd.path == "" {
		return "", errors.New("el path es obligatorio")
	}
	if cmd.unit == "" {
		cmd.unit = "M"
	}
	if cmd.fit == "" {
		cmd.fit = "FF"
	}
	sizeBytes, err := utils.ConvertToBytes(cmd.size, cmd.unit)
	if err != nil {
		return "", fmt.Errorf("error al convertir tamaño a bytes: %v", err)
	}
	err = cmk(cmd)
	if err != nil {
		return "", fmt.Errorf("error al crear el disco: %v", err)
	}

	result := fmt.Sprintf("Comando mkdisk ejecutado con éxito.\nDetalles:\n- Tamaño: %d bytes\n- Ajuste: %s\n- Ruta: %s",
		sizeBytes, cmd.fit, cmd.path)

	return result, nil
}

func cmk(mkdisk *MKDISK) error {
	sizeBytes, err := utils.ConvertToBytes(mkdisk.size, mkdisk.unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

	err = createDisk(mkdisk, sizeBytes)
	if err != nil {
		fmt.Println("Error creating disk:", err)
		return err
	}

	err = crearMBR(mkdisk, sizeBytes)
	if err != nil {
		fmt.Println("Error creating MBR:", err)
		return err
	}

	return nil
}

func createDisk(mkdisk *MKDISK, sizeBytes int) error {
	err := os.MkdirAll(filepath.Dir(mkdisk.path), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return err
	}

	file, err := os.Create(mkdisk.path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	buffer := make([]byte, 1024*1024) // Crea un buffer de 1 MB
	for sizeBytes > 0 {
		writeSize := len(buffer)
		if sizeBytes < writeSize {
			writeSize = sizeBytes
		}
		if _, err := file.Write(buffer[:writeSize]); err != nil {
			return err
		}
		sizeBytes -= writeSize
	}
	return nil
}

func crearMBR(mkdisk *MKDISK, sizeBytes int) error {
	mbr := &structures.MBR{
		Mbr_size:           int32(sizeBytes),
		Mbr_creation_date:  float32(time.Now().Unix()),
		Mbr_disk_signature: rand.Int31(),
		Mbr_disk_fit:       [1]byte{mkdisk.fit[0]},
		Mbr_partitions: [4]structures.PARTITION{
			{PartStatus: [1]byte{'N'}, PartType: [1]byte{'N'}, PartFit: [1]byte{'N'}, PartStart: -1, PartSize: -1, PartName: [16]byte{'P'}, PartCorrelative: 1, PartId: [4]byte{'0'}},
			{PartStatus: [1]byte{'N'}, PartType: [1]byte{'N'}, PartFit: [1]byte{'N'}, PartStart: -1, PartSize: -1, PartName: [16]byte{'P'}, PartCorrelative: 2, PartId: [4]byte{'0'}},
			{PartStatus: [1]byte{'N'}, PartType: [1]byte{'N'}, PartFit: [1]byte{'N'}, PartStart: -1, PartSize: -1, PartName: [16]byte{'P'}, PartCorrelative: 3, PartId: [4]byte{'0'}},
			{PartStatus: [1]byte{'N'}, PartType: [1]byte{'N'}, PartFit: [1]byte{'N'}, PartStart: -1, PartSize: -1, PartName: [16]byte{'N'}, PartCorrelative: 4, PartId: [4]byte{'0'}},
		},
	}
	err := mbr.SerializeMBR(mkdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}
