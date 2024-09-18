package comands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Login struct {
	Username string
	Password string
	Id       string
	Uid      int
	Gid      int
}

var logueado Login

func ValidarDatos(contexto []string) bool {
	id := ""
	usuario := ""
	contra := ""

	for _, token := range contexto {
		token = strings.TrimPrefix(token, "-")

		tk := strings.Split(token, "=")
		if len(tk) == 2 {
			switch tk[0] {
			case "id":
				id = tk[1]
			case "user":
				usuario = tk[1]
			case "pass":
				contra = tk[1]
			}
		}
	}

	return id != "" && usuario != "" && contra != ""
}

func LeerArchivos(partitionID string) (string, error) {
	partitionPath := fmt.Sprintf("/mnt/%s/users.txt", partitionID)

	// Verifica si el archivo existe
	if _, err := os.Stat(partitionPath); os.IsNotExist(err) {
		// Si no existe, crea el archivo con el grupo y usuario root
		fmt.Println("Archivo users.txt no encontrado. Generando archivo con usuario root...")

		err := crearArchivo(partitionPath)
		if err != nil {
			return "", fmt.Errorf("error al crear el archivo users.txt: %v", err)
		}
	}

	// Abre el archivo para leerlo
	file, err := os.Open(partitionPath)
	if err != nil {
		return "", fmt.Errorf("error al abrir el archivo users.txt: %v", err)
	}
	defer file.Close()

	// Lee el contenido del archivo
	var contenido strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contenido.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error al leer el archivo users.txt: %v", err)
	}

	return contenido.String(), nil
}

func crearArchivo(partitionPath string) error {
	file, err := os.Create(partitionPath)
	if err != nil {
		return fmt.Errorf("error al crear users.txt: %v", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString("1, G, root\n1, U, root, root, 123\n")
	if err != nil {
		return fmt.Errorf("error al escribir en users.txt: %v", err)
	}

	// Asegura que los datos se guarden en el archivo
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error al guardar en users.txt: %v", err)
	}

	fmt.Println("Archivo users.txt creado exitosamente con el usuario root.")
	return nil
}

func ValidarUsuario(usersData string, user string, pass string) bool {
	Lineas := strings.Split(usersData, "\n")
	for _, linea := range Lineas {
		parts := strings.Split(linea, ",")
		if len(parts) == 2 {
			if parts[0] == user && parts[1] == pass {
				return true
			}
		}
	}
	return false
}

func CLogin(tokens []string) (string, error) {

	// Verifica que se hayan proporcionado todos los datos necesarios
	if !ValidarDatos(tokens) {
		return "", errors.New("error en login: faltan datos")
	}

	// Verifica si ya hay un usuario logueado
	if logueado.Username != "" {
		return "", fmt.Errorf("error en login: ya hay un usuario logueado")
	}

	// Procesa los tokens para extraer los valores de login
	for _, token := range tokens {
		token = strings.TrimPrefix(token, "-")

		tk := strings.Split(token, "=")
		if len(tk) == 2 {
			switch tk[0] {
			case "id":
				logueado.Id = tk[1]
			case "user":
				logueado.Username = tk[1]
			case "pass":
				logueado.Password = tk[1]
			}
		}
	}

	// Leer el archivo users.txt de la partición montada
	DatosUsuario, err := LeerArchivos(logueado.Id)
	if err != nil {
		return "", fmt.Errorf("error al leer los archivos: %v", err)
	}

	// Verificar si el usuario y la contraseña son válidos
	if ValidarUsuario(DatosUsuario, logueado.Username, logueado.Password) {
		logueado.Uid = 1 // Aquí podrías cambiar el UID basado en el archivo
		logueado.Gid = 1 // También puedes ajustar el GID según los datos del usuario
		return fmt.Sprintf("Usuario %s logueado correctamente", logueado.Username), nil
	}

	return "", fmt.Errorf("usuario o contraseña incorrectos para %s", logueado.Username)
}
