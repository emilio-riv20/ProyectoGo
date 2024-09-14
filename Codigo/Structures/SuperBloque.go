package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type SuperBloque struct {
	FileSystemType  int32
	InodeCount      int32
	BlockCount      int32
	FreeInodes      int32
	FreeBlocks      int32
	MountTime       float32
	LastUnmountTime float32
	MountCount      int32
	MagicNumber     int32
	InodeSize       int32
	BlockSize       int32
	FirstInode      int32
	FirstBlock      int32
	BmInodeStart    int32
	BmBlockStart    int32
	InodeTableStart int32
	BlockTableStart int32
}

func (sb *SuperBloque) Serialize(path string, offset int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Seek(offset, 0); err != nil {
		return err
	}

	if err := binary.Write(file, binary.LittleEndian, sb); err != nil {
		return err
	}

	return nil
}

func (sb *SuperBloque) Deserialize(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Seek(offset, 0); err != nil {
		return err
	}

	bufferSize := binary.Size(sb)
	if bufferSize <= 0 {
		return fmt.Errorf("invalid SuperBloque size: %d", bufferSize)
	}

	buffer := make([]byte, bufferSize)
	if _, err := file.Read(buffer); err != nil {
		return err
	}

	reader := bytes.NewReader(buffer)
	if err := binary.Read(reader, binary.LittleEndian, sb); err != nil {
		return err
	}

	return nil
}

func (sb *SuperBloque) CreateUserFile(path string) error {
	// Step 1: Create root inode
	rootInode := sb.createInode(1, 1, sb.BlockCount, '0')
	if err := rootInode.Serializar(path, int64(sb.FirstInode)); err != nil {
		return err
	}

	sb.updateInodeStats()

	// Step 2: Create root directory block
	rootBlock := sb.createRootFolderBlock()
	if err := rootBlock.Serializar(path, int64(sb.FirstBlock)); err != nil {
		return err
	}

	sb.updateBlockStats()

	// Step 3: Update root inode with "users.txt"
	if err := sb.updateRootInodeWithUserFile(path); err != nil {
		return err
	}

	// Step 4: Create inode for "users.txt"
	usersText := "1,G,root\n1,U,root,123\n"
	usersInode := sb.createInode(1, 1, int32(len(usersText)), '1')
	if err := usersInode.Serializar(path, int64(sb.FirstInode)); err != nil {
		return err
	}

	sb.updateInodeStats()

	// Step 5: Create block for "users.txt"
	usersBlock := sb.createFileBlock(usersText)
	if err := usersBlock.Serializar(path, int64(sb.FirstBlock)); err != nil {
		return err
	}

	sb.updateBlockStats()

	return nil
}

func (sb *SuperBloque) createInode(uid, gid, blockCount int32, fileType byte) *INODO {
	return &INODO{
		Iuid:   uid,
		Igid:   gid,
		Isize:  0,
		Itime:  float32(time.Now().Unix()),
		Ictime: float32(time.Now().Unix()),
		Imtime: float32(time.Now().Unix()),
		Iblock: [15]int32{blockCount, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		Itype:  [1]byte{fileType},
		Iperm:  [3]byte{'7', '7', '7'},
	}
}

func (sb *SuperBloque) createRootFolderBlock() *BloqueFolder {
	return &BloqueFolder{
		B_content: [4]BloqueFolderContent{
			{B_name: [12]byte{'.'}, B_inodo: 0},
			{B_name: [12]byte{'.', '.'}, B_inodo: 0},
			{B_name: [12]byte{'-'}, B_inodo: -1},
			{B_name: [12]byte{'-'}, B_inodo: -1},
		},
	}
}

func (sb *SuperBloque) updateRootInodeWithUserFile(path string) error {
	rootInode := &INODO{}
	if err := rootInode.Deserializar(path, int64(sb.InodeTableStart)); err != nil {
		return err
	}

	rootInode.Itime = float32(time.Now().Unix())
	if err := rootInode.Serializar(path, int64(sb.InodeTableStart)); err != nil {
		return err
	}

	rootBlock := &BloqueFolder{}
	if err := rootBlock.Deserializar(path, int64(sb.BlockTableStart)); err != nil {
		return err
	}

	rootBlock.B_content[2] = BloqueFolderContent{
		B_name:  [12]byte{'u', 's', 'e', 'r', 's', '.', 't', 'x', 't'},
		B_inodo: sb.InodeCount,
	}

	if err := rootBlock.Serializar(path, int64(sb.BlockTableStart)); err != nil {
		return err
	}

	return nil
}

func (sb *SuperBloque) createFileBlock(content string) *ARCHIVOBLOQUE {
	fileBlock := &ARCHIVOBLOQUE{}
	copy(fileBlock.B_content[:], content)
	return fileBlock
}

func (sb *SuperBloque) updateInodeStats() {
	sb.InodeCount++
	sb.FreeInodes--
	sb.FirstInode += sb.InodeSize
}

func (sb *SuperBloque) updateBlockStats() {
	sb.BlockCount++
	sb.FreeBlocks--
	sb.FirstBlock += sb.BlockSize
}
