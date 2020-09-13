package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	CONTROLLER "../controllers"
	STRUCTURES "../structures"
)

func MakeAFileInLogicalDiskFirstTime(path string, id string, p string, size int64, cont string) {
	//primero voy a buscar la particion dentro de las particiones montadas
	if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
		//se obtiene la particion montada
		var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go
		//Se abre el archivo para modificarlo
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {

			//OBTENGO EL SUPER BLOQUE

			//Declaramos variable de tipo SUPERBOOT
			sb := GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA EN MKDIR.GO

			//CREAMOS LOS BLOQUES

			if len(cont) <= 25 {

				var block = STRUCTURES.DATABLOCK{}
				copy(block.DB_data[:], cont)
				//NOS SITUAMOS AL INICIO DEL BLOQUE DE DATOS
				file.Seek(sb.SB_ap_blocks, 0)
				block1 := &block
				var binario4 bytes.Buffer
				binary.Write(&binario4, binary.BigEndian, block1)
				escribirBytes(file, binario4.Bytes())
			} else {
				/*var str = ""
				var tope int = 0
				for i := 0; i < 25; i++ {
					str = str + string(cont[i])
					tope = i
				}
				fmt.Println(str)
				if (len(cont) - len(str)) <= 25 {

				}*/
			}

		}

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}
func MakeAFileInLogicalDisk(path string, id string, p string, size int64, cont string) {

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
				arregloPath := strings.Split(path, "/")
				strPath := ""
				for i := 1; i < len(arregloPath)-1; i++ {
					strPath += "/" + arregloPath[i]
				}
				nameFile := arregloPath[len(arregloPath)-1]

				//CREA LAS CARPETAS SI NO EXISTIERAN
				var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO
				CONTROLLER.MakeAnDirectoryInDisk(sb, strPath, p, file, (CONTROLLER.GetLogedUser()).User_username)

				CONTROLLER.MakeAVD(sb, strPath, nameFile, file, (CONTROLLER.GetLogedUser()).User_id, cont)

				//CONTROLLER.SearchDetailDirectory(sb, strPath, nameFile, file, (CONTROLLER.GetLogedUser()).User_id, cont);
				//Avd_ap_detalle_directorio =
			}

		} else {
			fmt.Println("No existe una particion con el id " + id)
		}

	} else {
		fmt.Println("No hay ningun usuario logueado")
	}
}
