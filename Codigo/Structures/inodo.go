package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type INODO struct {
	Iuid   int32
	Igid   int32
	Isize  int32
	Itime  float32
	Ictime float32
	Imtime float32
	Iblock [15]int32
	Itype  [1]byte
	Iperm  [3]byte
}

func (inode *INODO) Serializar(path string, offset int64) error {
	archivo, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer archivo.Close()
	_, err = archivo.Seek(offset, 0)
	if err != nil {
		return err
	}

	err = binary.Write(archivo, binary.LittleEndian, inode)
	if err != nil {
		return err
	}

	return nil
}

func (inode *INODO) Deserializar(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	InodoSize := binary.Size(inode)
	if InodoSize <= 0 {
		return fmt.Errorf("invalid Inode size: %d", InodoSize)
	}

	buffer := make([]byte, InodoSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, inode)
	if err != nil {
		return err
	}

	return nil
}

func NewInode() *INODO {
	return &INODO{
		Iuid:   1,
		Igid:   1,
		Isize:  0,
		Itime:  float32(time.Now().Unix()),
		Ictime: float32(time.Now().Unix()),
		Imtime: float32(time.Now().Unix()),
		Iblock: [15]int32{},
		Itype:  [1]byte{'0'},
		Iperm:  [3]byte{'7', '7', '5'},
	}
}

func (inode *INODO) Print() {
	atime := time.Unix(int64(inode.Itime), 0)
	ctime := time.Unix(int64(inode.Ictime), 0)
	mtime := time.Unix(int64(inode.Imtime), 0)

	fmt.Printf("I_uid: %d\n", inode.Iuid)
	fmt.Printf("I_gid: %d\n", inode.Igid)
	fmt.Printf("I_s: %d\n", inode.Isize)
	fmt.Printf("I_atime: %s\n", atime.Format(time.RFC3339))
	fmt.Printf("I_ctime: %s\n", ctime.Format(time.RFC3339))
	fmt.Printf("I_mtime: %s\n", mtime.Format(time.RFC3339))
	fmt.Printf("I_block: %v\n", inode.Iblock)
	fmt.Printf("I_type: %s\n", string(inode.Itype[:]))
	fmt.Printf("I_perm: %s\n", string(inode.Iperm[:]))
}
