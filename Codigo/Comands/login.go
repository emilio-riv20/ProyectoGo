package comands

import (
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
	partitionPath := "/path/to/mounted/partition"

	file, err := os.Open(partitionPath)
	if err != nil {
		return "", fmt.Errorf("error al abrir la partición: %v", err)
	}
	defer file.Close()
	ArchivosUser := "root,123\nadmin,admin123\n"

	return ArchivosUser, nil
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
