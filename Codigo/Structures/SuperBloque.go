package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type SuperBloque struct {
	S_filesystem_type int32
	S_inodes_count    int32
	S_blocks_count    int32
	S_free_blocks     int32
	S_free_inodes     int32
	S_mtime           int32
	S_umtime          int32
	S_mnt_count       int16
	S_magic           int16
	S_inode_size      int16
	S_block_size      int16
	S_firs_ino        int32
	S_first_blo       int32
	S_bm_inode_start  int32
	S_bm_block_start  int32
	S_inode_start     int32
	S_block_start     int32
}

func (sb *SuperBloque) Serializar(path string, offset int64) error {
	archivo, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer archivo.Close()

	_, err = archivo.Seek(offset, 0)
	if err != nil {
		return err
	}

	err = binary.Write(archivo, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}

func (sb *SuperBloque) Deserializar(path string, offset int64) error {
	archivo, err := os.Open(path)
	if err != nil {
		return err
	}
	defer archivo.Close()

	_, err = archivo.Seek(offset, 0)
	if err != nil {
		return err
	}

	SuperSize := binary.Size(sb)
	if SuperSize <= 0 {
		return fmt.Errorf("tamaño de SuperBloque inválido: %d", SuperSize)
	}

	buffer := make([]byte, SuperSize)
	_, err = archivo.Read(buffer)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}
