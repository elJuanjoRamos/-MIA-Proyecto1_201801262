package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"unsafe"

	CONTROLLER "../controllers"
	FUNCTION "../functions"
	STRUCTURES "../structures"
)

var MountList [100][26]STRUCTURES.MOUNT

func MountPrint() {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if (MountList[i][j] != STRUCTURES.MOUNT{}) {
				fmt.Println(MountList[i][j])
				fmt.Println("")
			}
		}
		//fmt.Println("")
	}
}

func Mount(path string, name string) {
	//Se verifica si existe la ruta o archivo
	if FUNCTION.IfExistFile(path) {

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

			//BUSCA SI ES UNA PARTICION PRIMARIA
			if GetPartitionName(m.Mbr_partition_1.Part_name) == name &&
				GetPartitionName([16]byte{m.Mbr_partition_1.Part_type}) == "P" {
				InsertPart(m.Mbr_partition_1, path, name)
				fmt.Println("Particion Montada correctamente")
			} else if GetPartitionName(m.Mbr_partition_2.Part_name) == name &&
				GetPartitionName([16]byte{m.Mbr_partition_2.Part_type}) == "P" {
				InsertPart(m.Mbr_partition_2, path, name)
				fmt.Println("Particion Montada correctamente")

			} else if GetPartitionName(m.Mbr_partition_3.Part_name) == name &&
				GetPartitionName([16]byte{m.Mbr_partition_1.Part_type}) == "P" {
				InsertPart(m.Mbr_partition_3, path, name)
				fmt.Println("Particion Montada correctamente")

			} else if GetPartitionName(m.Mbr_partition_4.Part_name) == name {
				InsertPart(m.Mbr_partition_4, path, name)
				fmt.Println("Particion Montada correctamente")
			} else {
				filename := filepath.Base(path)
				//No se encuentra se va a buscar a las logicas
				var partExist = CONTROLLER.SearchPartition(filename, name)

				if partExist {
					InsertPart(CONTROLLER.GetLogicPartition(filename, name), path, name)
					fmt.Println("Particion Montada correctamente")

					//Si tampoco se encuentra en las logicas, se muestra un mensaje de error
				} else {
					fmt.Println("No existe la particion con nombre " + name + " en el disco " + filename)
				}

			}

		} else {
			fmt.Println("No existen particiones dentro del disco")
		}
	} else {
		fmt.Println("El disco o ruta no existe")
	}

}

func InsertPart(part STRUCTURES.PARTITION, path string, name string) {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if (MountList[i][j] == STRUCTURES.MOUNT{}) {
				asciiNum := 97 + j // Uppercase A
				character := string(asciiNum)
				s := strconv.Itoa(i + 1)
				MountList[i][j] = STRUCTURES.MOUNT{
					Mount_id:        "vd" + character + s,
					Mount_path:      path,
					Mount_particion: name,
					Mount_estado:    true,
					Mount_part:      part,
				}
				return
			}
		}
	}
}

func UMount(id string) {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if MountList[i][j].Mount_id == id {
				MountList[i][j] = STRUCTURES.MOUNT{}
				return
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

//BUsca una particion por id, retorna true si existe, false si no
func SearchPartitionById(id string) bool {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if MountList[i][j].Mount_id == id {
				return true
			}
		}
	}
	return false
}

//Busca una particion por id, retorna la particion montada
func GetPartitionById(id string) STRUCTURES.MOUNT {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if MountList[i][j].Mount_id == id {
				return MountList[i][j]
			}
		}
	}
	return STRUCTURES.MOUNT{}
}

//update user txt, le da la ruta del archivo txt de usuarios asignado a esa particion
func UpdateUserTxt(id string, path string) {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if MountList[i][j].Mount_id == id {
				MountList[i][j].Mount_usrtxt = path
				break
			}
		}
	}
}
