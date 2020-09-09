package commands

import (
	"fmt"
	"os"

	CONTROLLER "../controllers"
)

func MakeADir(path string, id string, p string) {

	if CONTROLLER.IsLogged() {

		if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
			//se obtiene la particion montada
			var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go
			//Se abre el archivo para modificarlo
			file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
			defer file.Close()
			if err != nil {
				fmt.Println(err)
			} else {

				//OBTENGO EL SUPER BOOT
				var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO

				CONTROLLER.MakeAnDirectoryInDisk(sb, path, p, file, (CONTROLLER.GetLogedUser()).User_username)

			}

		} else {
			fmt.Println("No existe una particion con el id " + id)
		}

	} else {
		fmt.Println("No hay ningun usuario logueado")
	}
}
