package structures

import (
	"bytes"
	"encoding/binary"
	"os"
)

type BloqueFolder struct {
	B_content [4]BloqueFolderContent
}

type BloqueFolderContent struct {
	B_name  [12]byte
	B_inodo int32
}

func (bf *BloqueFolder) Serializar(path string, offset int64) error {
	archivo, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer archivo.Close()

	_, err = archivo.Seek(offset, 0)
	if err != nil {
		return err
	}

	err = binary.Write(archivo, binary.LittleEndian, bf)
	if err != nil {
		return err
	}

	return nil

}

func (bf *BloqueFolder) Deserializar(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	BloqueFolderSize := binary.Size(bf)
	if BloqueFolderSize <= 0 {
		return err
	}

	buffer := make([]byte, BloqueFolderSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, bf)
	if err != nil {
		return err
	}

	return nil
}
