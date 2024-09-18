package utils

import (
	"errors"

	structures "github.com/emilio-riv20/proyecto1/Codigo/Structures"
)

const Carnet string = "12"

var (
	MountedPartitions map[string]string = make(map[string]string)
)

func ParticionMontada(id string) (*structures.PARTITION, string, error) {
	path := MountedPartitions[id]
	if path == "" {
		return nil, "", errors.New("la partición no está montada")
	}

	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(path)
	if err != nil {
		return nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, _ := mbr.ParticionPorId(id)
	if partition == nil {
		return nil, "", errors.New("partición no encontrada")
	}

	return partition, path, nil
}

func GetParticionMontada(id string) (*structures.MBR, *structures.SuperBloque, string, error) {
	path := MountedPartitions[id]
	if path == "" {
		return nil, nil, "", errors.New("la partición no está montada")
	}
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(path)
	if err != nil {
		return nil, nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, partitionIndex := mbr.ParticionPorId(id)
	if partition == nil {
		return nil, nil, "", errors.New("partición no encontrada")
	}

	var sb structures.SuperBloque

	err = sb.Deserialize(path, int64(mbr.Mbr_partitions[partitionIndex].PartStart))
	if err != nil {
		return nil, nil, "", err
	}

	return &mbr, &sb, path, nil
}
