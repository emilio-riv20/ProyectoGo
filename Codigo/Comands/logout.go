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
		return "No hay usuario logueado", nil
	}

	logueado = Login{}
	return "sesion cerrada", nil
}
