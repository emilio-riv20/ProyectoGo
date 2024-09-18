package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func Int32ToBytes(n int32) [4]byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(n))
	return buf
}

func Float64ToBytes(f float64) [4]byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(f))
	return buf
}

func ConvertToBytes(size int, unit string) (int, error) {
	switch unit {
	case "K":
		return size * 1024, nil
	case "M":
		return size * 1024 * 1024, nil
	default:
		return 0, errors.New("invalid unit")
	}
}

var abecedario = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

var path_Letra = make(map[string]string)
var sig = 0

func Letra(path string) (string, error) {
	if _, existe := path_Letra[path]; !existe {
		if sig < len(abecedario) {
			path_Letra[path] = abecedario[sig]
			sig++
		} else {
			fmt.Println("Ya no hay letras disponibles")
			return "", errors.New("ya no hay letras disponibles")
		}
	}
	return path_Letra[path], nil
}

func GetLetter(path string) (string, error) {
	if letra, existe := path_Letra[path]; existe {
		return letra, nil
	}
	return Letra(path)
}
