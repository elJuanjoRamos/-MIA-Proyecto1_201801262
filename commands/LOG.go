package commands

import (
	"time"
)

type BITACORA struct {
	operacion int
	tipo      string
	nombre    string
	contenido string
	fecha     string
}

var arreglo []BITACORA

func AddAccion(op int, t, n, c string) {
	var time = time.Now()

	var bit = BITACORA{
		operacion: op,
		tipo:      t,
		nombre:    n,
		contenido: c,
		fecha:     time.Format("2006-01-02 15:04:05"),
	}
	arreglo = append(arreglo, bit)

}

func GetBitacora() []BITACORA {
	return arreglo
}
