package comands

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type MKGRP struct {
	name string
}

func LeerMKGRP(tokens []string) (string, error) {
	cmd := &MKGRP{}
	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-name="[^"]+"|-name=[^\s]+`)
	matches := re.FindAllString(args, -1)

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
		case "-name":
			cmd.name = value
		}
	}

	if cmd.name == "" {
		return "", errors.New("el nombre del grupo es obligatorio")
	}

	return cmd.name, nil
}

func GrupoExiste(usersData string, groupName string) bool {
	Lineas := strings.Split(usersData, "\n")
	for _, linea := range Lineas {
		parts := strings.Split(linea, ",")
		if len(parts) >= 3 && parts[1] == "G" && parts[2] == groupName {
			return true
		}
	}
	return false
}

func ObtenerNuevoID(usersData string) int {
	Lineas := strings.Split(usersData, "\n")
	maxID := 0
	for _, linea := range Lineas {
		parts := strings.Split(linea, ",")
		if len(parts) > 0 {
			id, err := strconv.Atoi(parts[0])
			if err == nil && id > maxID {
				maxID = id
			}
		}
	}
	return maxID + 1
}

func EscribirLineaArchivo(partitionID string, nuevaLinea string) error {
	partitionPath := fmt.Sprintf("/mnt/%s/users.txt", partitionID)

	file, err := os.OpenFile(partitionPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo users.txt: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(nuevaLinea)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo users.txt: %v", err)
	}

	return nil
}

func EjecutarMKGRP(partitionID string, groupName string, isRoot bool) error {
	if !isRoot {
		return errors.New("solo el usuario root puede crear grupos")
	}

	usersData, err := LeerArchivos(partitionID)
	if err != nil {
		return fmt.Errorf("error al leer el archivo users.txt: %v", err)
	}

	// Verificar si el grupo ya existe
	if GrupoExiste(usersData, groupName) {
		return fmt.Errorf("el grupo %s ya existe", groupName)
	}

	newID := ObtenerNuevoID(usersData)

	nuevaLinea := fmt.Sprintf("%d, G, %s\n", newID, groupName)

	err = EscribirLineaArchivo(partitionID, nuevaLinea)
	if err != nil {
		return fmt.Errorf("error al escribir el nuevo grupo: %v", err)
	}

	fmt.Println("Grupo creado exitosamente:", groupName)
	return nil
}
