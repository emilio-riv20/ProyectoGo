package comands

import (
	"fmt"
)

func Clogout() (string, error) {
	//Verificar si hay usuario logueado
	if logueado.Username != "" {
		logueado = Login{}
		fmt.Println("Usuario deslogueado")
	} else {
		fmt.Println("No hay usuario logueado")
	}

	logueado = Login{}
	return "sesion cerrada", nil
}
