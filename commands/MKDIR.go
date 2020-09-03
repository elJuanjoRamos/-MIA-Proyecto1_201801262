package commands

import (
	"fmt"

	CONTROLLER "../controllers"
	FUNCTIONS "../functions"
)

func MakeADir(path string, id string, p string) {
	//primero voy a buscar la particion dentro de las particiones montadas
	if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
		//se obtiene la particion montada
		//var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go

		if CONTROLLER.IsLogged() {
			if !FUNCTIONS.IfExistDirectoryOrPath(path) {
				if p != "" { //verifico si no viene vacia, si no viene, se manda a crear
					FUNCTIONS.CreateADirectory(path, 0777) // Se manda a crear
				} else {
					fmt.Println("No cuenta con el parametro que permita crear los directorios")
				}
			} else {
				fmt.Println("La carpeta ya existe")
			}
		} else {
			fmt.Println("No hay ningun usuario logueado")
		}

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}
