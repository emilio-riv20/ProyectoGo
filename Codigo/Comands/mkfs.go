package comands

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	structures "github.com/emilio-riv20/proyecto1/Codigo/Structures"
	utils "github.com/emilio-riv20/proyecto1/Codigo/utils"
)

type MFKS struct {
	id     string
	format string
}

func calcN(partition *structures.PARTITION) int32 {

	numerator := int(partition.PartSize) - binary.Size(structures.SuperBloque{})
	denominator := 4 + binary.Size(structures.INODO{}) + 3*binary.Size(structures.ARCHIVOBLOQUE{})
	n := math.Floor(float64(numerator) / float64(denominator))

	return int32(n)
}

func CrearSupBloque(partition *structures.PARTITION, n int32) *structures.SuperBloque {
	bm_inode_start := partition.PartStart + int32(binary.Size(structures.SuperBloque{}))
	bm_block_start := bm_inode_start + n
	// Inodos
	inode_start := bm_block_start + (3 * n)
	// Bloques
	block_start := inode_start + (int32(binary.Size(structures.INODO{})) * n)

	// Crear superbloque
	superBloque := &structures.SuperBloque{
		FileSystemType:  2,
		InodeCount:      0,
		BlockCount:      0,
		FreeInodes:      int32(n),
		FreeBlocks:      int32(n * 3),
		MountTime:       float32(time.Now().Unix()),
		LastUnmountTime: float32(time.Now().Unix()),
		MountCount:      1,
		MagicNumber:     0xEF53,
		InodeSize:       int32(binary.Size(structures.INODO{})),
		BlockSize:       int32(binary.Size(structures.ARCHIVOBLOQUE{})),
		FirstInode:      inode_start,
		FirstBlock:      block_start,
		BmInodeStart:    bm_inode_start,
		BmBlockStart:    bm_block_start,
		InodeTableStart: inode_start,
		BlockTableStart: block_start,
	}
	return superBloque
}

func Cmkfs(cmd *MFKS) error {
	PartMont, path, err := utils.ParticionMontada(cmd.id)
	if err != nil {
		return err
	}

	n := calcN(PartMont)
	superBlock := CrearSupBloque(PartMont, n)

	// Crear los bitmaps
	err = superBlock.CrearBitmap(path)
	if err != nil {
		return err
	}

	// Crear archivo users.txt
	err = superBlock.CreateUserFile(path)
	if err != nil {
		return err
	}

	// Serializar el superbloque
	err = superBlock.Serialize(path, int64(PartMont.PartStart))
	if err != nil {
		return err
	}

	return nil

}

func LeerMkfs(tokens []string) (string, error) {

	cmd := &MFKS{}
	args := strings.Join(tokens, "")
	er := regexp.MustCompile(`(?i)-id=[^\s]+|-type=[^\s]+`)
	matches := er.FindAllString(args, -1)

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
		case "-id":
			if value == "" {
				return "", errors.New("el id es obligatorio")
			}
			cmd.id = value

		case "-format":
			if value != "full" {
				return "", errors.New("el formato debe ser full")
			}
			cmd.format = value
		default:
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if cmd.id == "" {
		return "", errors.New("el id es obligatorio")
	}

	if cmd.format == "" {
		cmd.format = "full"
	}

	err := Cmkfs(cmd)
	if err != nil {
		return "", err
	}
	res := fmt.Sprintf("Formateo de disco con id %s completado", cmd.id)

	return res, nil
}
