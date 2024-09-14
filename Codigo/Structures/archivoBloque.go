package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type ARCHIVOBLOQUE struct {
	B_content [64]byte
}

// Serialize escribe la estructura FileBlock en un archivo binario en la posición especificada
func (fb *ARCHIVOBLOQUE) Serializar(path string, offset int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Serializar la estructura FileBlock directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, fb)
	if err != nil {
		return err
	}

	return nil
}

func (fb *ARCHIVOBLOQUE) Deserializar(path string, offset int64) error {
	archivo, err := os.Open(path)
	if err != nil {
		return err
	}
	defer archivo.Close()

	_, err = archivo.Seek(offset, 0)
	if err != nil {
		return err
	}

	fbSize := binary.Size(fb)
	if fbSize <= 0 {
		return fmt.Errorf("tamano invalido: %d", fbSize)
	}

	buffer := make([]byte, fbSize)
	_, err = archivo.Read(buffer)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, fb)
	if err != nil {
		return err
	}

	return nil
}
