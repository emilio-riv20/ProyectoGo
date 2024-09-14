package structures

import (
	"encoding/binary"
	"os"
)

func (sb *SuperBloque) CrearBitmap(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(int64(sb.BmInodeStart), 0)
	if err != nil {
		return err
	}

	buffer := make([]byte, sb.FreeInodes)
	for i := range buffer {
		buffer[i] = '0'
	}

	err = binary.Write(file, binary.LittleEndian, buffer)
	if err != nil {
		return err
	}

	_, err = file.Seek(int64(sb.BmBlockStart), 0)
	if err != nil {
		return err
	}

	buffer = make([]byte, sb.FreeBlocks)
	for i := range buffer {
		buffer[i] = 'O'
	}

	err = binary.Write(file, binary.LittleEndian, buffer)
	if err != nil {
		return err
	}

	return nil
}

func (sb *SuperBloque) ActualizarInodo(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(int64(sb.BmInodeStart)+int64(sb.InodeCount), 0)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte{'1'})
	if err != nil {
		return err
	}

	return nil
}

func (sb *SuperBloque) ActualizarBLoque(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(int64(sb.BmBlockStart)+int64(sb.BlockCount), 0)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte{'X'})
	if err != nil {
		return err
	}

	return nil
}
