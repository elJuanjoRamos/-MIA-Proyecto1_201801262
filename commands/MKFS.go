package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	FUNCTION "../functions"
	STRUCTURES "../structures"
)

/*Variables Globales*/
/*var idUser = 1
var idGroup = 1*/

//funcion que se encarga de formatear la particion
func MKFSFormatPartition(id string, types string) {
	//Se crea el directorio que contiene los txt de usuarios y grupos
	FUNCTION.CreateADirectory(FUNCTION.RootDir() + "/reports/userfiles")

	//primero voy a buscar la particion dentro de las particiones montadas
	if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
		//se obtiene la particion montada
		var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go

		//Se abre el archivo para irlo a limpiar
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {

			var pathFile = FUNCTION.RootDir() + "/reports/userfiles/" + "users_" + id + ".txt" //La funcion RootDir retorna la ruta del proyecto, se encuentra en special_functions.go
			// Mando a actualizar la particion montada, dando la ruta del archivo de usuarios, esto para no tener que leer en el disco la ruta
			UpdateUserTxt(id, pathFile)
			//Se crea el archivo txt con los usuario iniciales, cambiar carnet
			FUNCTION.CreateAFile(pathFile, "1, G, root\n1, U, root, root , 201801262")

			//Se envia a formatear la particion
			FastAndFullPartition(partition.Mount_part, file, strings.ToLower(types), pathFile)

		}

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}

func FastAndFullPartition(partition STRUCTURES.PARTITION, file *os.File, types string, pathFile string) {
	//Formatea la particion segun lo que venga
	for i := partition.Part_start; i < partition.Part_end+1; i++ {
		var init int8 = '0'
		if types == "full" {
			init = 0
		}
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escribirBytes(file, binarioTemp.Bytes())
	}

	//Nos posicionamos al inicio de la particion
	file.Seek(partition.Part_start, 0)
	//Escribe en el inicio de la particion, la ruta del path del archivo de usuarios
	var path = []byte(pathFile)
	o := &path
	var binarioTemp bytes.Buffer
	binary.Write(&binarioTemp, binary.BigEndian, o)
	escribirBytes(file, binarioTemp.Bytes())
}
