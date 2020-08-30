package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	STRUCTURES "../structures"
)

//funcion que se encarga de formatear la particion
func MKFSFormatPartition(id string, types string) {
	//primero voy a buscar la particion dentro de las particiones montadas

	if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
		//se obtiene la particion montada
		var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go

		//Se abre el archivo para irlo a llenar
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			FastAndFullPartition(partition.Mount_part, file, strings.ToLower(types))

		}

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}

func FastAndFullPartition(partition STRUCTURES.PARTITION, file *os.File, types string) {
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
}

func FullPartition(partition STRUCTURES.PARTITION, file *os.File) {
	for i := partition.Part_start; i < partition.Part_end+1; i++ {
		var init int8 = 0
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escribirBytes(file, binarioTemp.Bytes())
	}
}
