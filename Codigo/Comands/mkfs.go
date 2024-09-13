package comands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	structures "github.com/emilio-riv20/proyecto1/Codigo/Structures"
)

var particionMontada map[string]*structures.PARTITION

func obtenerParticion(id string) (*structures.PARTITION, error) {
	if p, existe := particionMontada[id]; existe {
		return p, nil
	}
	return nil, errors.New("no se ha montado la particion")
}

func CrearArchivo() {
	_, err := os.Create("archivo.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func estructuraExt2(partition *structures.PARTITION) {
	//CrearArchivo()
	//fmt.Println("Creando estructura EXT2")
}

func ComandoMkfs(id string, format string) error {

	partition, err := obtenerParticion(id)
	if err != nil {
		return err
	}
	if strings.ToLower(format) == "full" || format == "" {
		fmt.Printf("Realizando formateo completo en la partición %s...\n", partition.PartName)
		// Aquí se simula el formateo completo de la partición
	} else {
		return errors.New("error: tipo de formateo no soportado")
	}
	estructuraExt2(partition)
	CrearArchivo()

	return nil
}
