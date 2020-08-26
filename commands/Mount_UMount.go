package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"unsafe"

	FUNCTION "../functions"
	STRUCTURES "../structures"
)

var MountList [100][26]STRUCTURES.MOUNT

func MountPrint() {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if (MountList[i][j] != STRUCTURES.MOUNT{}) {
				fmt.Print(MountList[i][j])
				fmt.Println("")
			}
		}
		//fmt.Println("")
	}
}

func Mount(path string, name string) {

	//Se verifica si existe la ruta o archivo
	if FUNCTION.IfExistDirectoryOrPath(path) {

		//Abrimos un archivo.
		file, err := os.Open(path)
		defer file.Close()
		if err != nil { //validar que no sea nulo.
			log.Fatal(err)
		}

		//Declaramos variable de tipo mbr
		m := STRUCTURES.MBR{}
		//Obtenemos el tamanio del mbr
		var size int = int(unsafe.Sizeof(m))

		//Lee la cantidad de <size> bytes del archivo
		data := leerBytes(file, size)
		//Convierte la data en un buffer,necesario para
		//decodificar binario
		buffer := bytes.NewBuffer(data)

		//Decodificamos y guardamos en la variable m
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		//Se busca la particion con nombre especificado
		if m.Mbr_count != 4 {
			//Se verifica si alguna particion cumple con el nombre que viene
			if GetPartitionName(m.Mbr_partition_1.Part_name) == name ||
				GetPartitionName(m.Mbr_partition_2.Part_name) == name ||
				GetPartitionName(m.Mbr_partition_3.Part_name) == name ||
				GetPartitionName(m.Mbr_partition_4.Part_name) == name {

				for i := 0; i < 100; i++ {
					for j := 0; j < 26; j++ {
						if (MountList[i][j] == STRUCTURES.MOUNT{}) {
							asciiNum := 97 + j // Uppercase A
							character := string(asciiNum)

							s := strconv.Itoa(i + 1)
							mountPartition := STRUCTURES.MOUNT{}
							mountPartition.Mount_id = "vd" + character + s
							mountPartition.Mount_path = path
							mountPartition.Mount_particion = name
							mountPartition.Mount_estado = true
							MountList[i][j] = mountPartition
							return
						}
					}
				}
			} else {
				fmt.Println("No existe una particion con el nombre deseado")
			}
		} else {
			fmt.Println("No existen particiones dentro del disco")
		}
	} else {
		fmt.Println("El disco o ruta no existe")
	}

}

func UMount(id string) {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if MountList[i][j].Mount_id == id {
				MountList[i][j] = STRUCTURES.MOUNT{}
				return
			} else {
				fmt.Println("No existe una particion con ese id")
			}
		}
	}
}

func GetPartitionName(name [16]byte) string {
	var s string = ""
	for _, v := range name {
		if v != 0 {
			s = s + string(v)
		}
	}
	return s
}
