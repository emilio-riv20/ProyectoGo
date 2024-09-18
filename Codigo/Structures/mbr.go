package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

type MBR struct {
	Mbr_size           int32
	Mbr_creation_date  int64
	Mbr_disk_signature int32
	Mbr_disk_fit       [1]byte
	Mbr_partitions     [4]PARTITION
}

func (mbr *MBR) SerializeMBR(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

func (mbr *MBR) DeserializeMBR(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	mbrSize := binary.Size(mbr)
	if mbrSize <= 0 {
		return fmt.Errorf("invalid MBR size: %d", mbrSize)
	}
	buffer := make([]byte, mbrSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

func (mbr *MBR) Print() {
	creationTime := time.Unix(int64(mbr.Mbr_creation_date), 0)
	diskFit := rune(mbr.Mbr_disk_fit[0])

	fmt.Printf("MBR Size: %d\n", mbr.Mbr_size)
	fmt.Printf("Creation Date: %s\n", creationTime.Format(time.RFC3339))
	fmt.Printf("Disk Signature: %d\n", mbr.Mbr_disk_signature)
	fmt.Printf("Disk Fit: %c\n", diskFit)
}

func (mbr *MBR) PrintPartitions() {
	for i, partition := range mbr.Mbr_partitions {
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(partition.PartStatus[0])
		partType := rune(partition.PartType[0])
		partFit := rune(partition.PartFit[0])

		// Convertir Part_name a string
		partName := string(partition.PartName[:])

		fmt.Printf("Partition %d:\n", i+1)
		fmt.Printf("  Status: %c\n", partStatus)
		fmt.Printf("  Type: %c\n", partType)
		fmt.Printf("  Fit: %c\n", partFit)
		fmt.Printf("  Start: %d\n", partition.PartStart)
		fmt.Printf("  Size: %d\n", partition.PartSize)
		fmt.Printf("  Name: %s\n", partName)
		fmt.Printf("  Correlative: %d\n", partition.PartCorrelative)
		fmt.Printf("  ID: %d\n", partition.PartId)
	}
}

func (mbr *MBR) ParticionPorNombre(name string) (*PARTITION, int) {
	for i, partition := range mbr.Mbr_partitions {
		partitionName := strings.Trim(string(partition.PartName[:]), "\x00 ")
		inputName := strings.Trim(name, "\x00 ")

		if strings.EqualFold(partitionName, inputName) {
			return &partition, i
		}
	}
	return nil, -1
}

func (mbr *MBR) ParticionPorId(id string) (*PARTITION, int) {
	for i, partition := range mbr.Mbr_partitions {

		partitionID := strings.Trim(string(partition.PartId[:]), "\x00 ")
		inputID := strings.Trim(id, "\x00 ")
		if strings.EqualFold(partitionID, inputID) {
			return &partition, i
		}
	}
	return nil, -1
}

func (mbr *MBR) GetFirstAvailablePartition() (*PARTITION, int, int) {
	offset := binary.Size(mbr)
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		if mbr.Mbr_partitions[i].PartStart == -1 {
			return &mbr.Mbr_partitions[i], offset, i
		} else {
			offset += int(mbr.Mbr_partitions[i].PartSize)
		}
	}
	return nil, -1, -1
}

func (mbr *MBR) HasExtendedPartition() bool {
	for _, partition := range mbr.Mbr_partitions {
		if partition.PartType[0] == 'E' {
			return true
		}
	}
	return false
}
