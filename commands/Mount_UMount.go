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

	fmt.Println("")
	fmt.Println("=====================================")
	fmt.Println("> Listado de Particiones montadas")
	fmt.Println("=====================================")
	fmt.Println("")

	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if (MountList[i][j] != STRUCTURES.MOUNT{}) {
				var mount = MountList[i][j]
				fmt.Println("Partition ID: " + mount.Mount_id)
				fmt.Println("Partition Path: " + mount.Mount_path)
				fmt.Println("Partition Name: " + mount.Mount_particion)
				fmt.Println("Partition State: ", mount.Mount_estado)
				fmt.Println("Partition Size: ", mount.Mount_part.Part_size)
				fmt.Println("Partition Type: ", mount.Mount_part.Part_type)
				fmt.Println("Partition Fit: ", mount.Mount_part.Part_fit)
				fmt.Println("Partition Start: ", mount.Mount_part.Part_start)
				fmt.Println("-----------")

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

		var PartitionArray = [4]STRUCTURES.PARTITION{}
		PartitionArray[0] = m.Mbr_partition_1
		PartitionArray[1] = m.Mbr_partition_2
		PartitionArray[2] = m.Mbr_partition_3
		PartitionArray[3] = m.Mbr_partition_4

		var bandera = false
		for i := 0; i < 4; i++ { //BUsco en las particiones primarias
			var partition = PartitionArray[i]
			if partition.Part_isEmpty == 1 {
				if GetPartitionName(partition.Part_name) == name && string(partition.Part_type) == "P" {
					InsertPart(partition, path, name)
					bandera = true
					break
				}

			}
		}
		if !bandera { // si no esta en las primarias, la busco en las logicas
			for i := 0; i < 4; i++ {
				var partition = PartitionArray[i]
				if partition.Part_isEmpty == 1 {
					if string(partition.Part_type) == "E" {

						filename := filepath.Base(path)
						//No se encuentra se va a buscar a las logicas
						var partExist = CONTROLLER.SearchPartition(filename, name)

						if partExist {
							InsertPart(CONTROLLER.GetLogicPartition(filename, name), path, name)
							//Si tampoco se encuentra en las logicas, se muestra un mensaje de error
						} else {
							fmt.Println("")
							fmt.Println("============================================================================")
							fmt.Println("> No existe la particion con nombre " + name + " en el disco " + filename)
							fmt.Println("============================================================================")
							fmt.Println("")
						}
						break
					}

				}
			}
		}

	} else {
		fmt.Println("El disco o ruta no existe")
	}

}

func InsertPart(part STRUCTURES.PARTITION, path string, name string) {
	var encontrado = false
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if (MountList[i][j] != STRUCTURES.MOUNT{}) {
				if MountList[i][j].Mount_particion == name {
					encontrado = true
					break
				}
			}
		}
	}

	if encontrado {
		fmt.Println("")
		fmt.Println("============================================")
		fmt.Println("> La particion " + name + " ya esta montada")
		fmt.Println("============================================")
		fmt.Println("")
	} else {
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
					fmt.Println("")
					fmt.Println("=====================================")
					fmt.Println("> Particion Montada correctamente")
					fmt.Println("=====================================")
					fmt.Println("")

					return
				}
			}
		}
	}

}

func UMount(ids []string) {

	for k := 0; k < len(ids); k++ {
		var bandera = false
		for i := 0; i < 100; i++ {
			for j := 0; j < 26; j++ {
				if MountList[i][j].Mount_id == ids[k] {
					bandera = true
					MountList[i][j] = STRUCTURES.MOUNT{}
					break
				}
			}
		}
		if bandera {
			fmt.Println("")
			fmt.Println("==================================================")
			fmt.Println("> Particion " + ids[k] + " desmontada correctamente")
			fmt.Println("==================================================")
			fmt.Println("")
		} else {
			fmt.Println("")
			fmt.Println("==================================================")
			fmt.Println("> Particion " + ids[k] + " no se encuentra montada")
			fmt.Println("==================================================")
			fmt.Println("")
		}
	}
}

func GetPartitionName(name [16]byte) string {
	var s string = ""
	for _, v := range name {
		if v != 0 {
			if string(v) == "+" {
				break
			} else {
				s = s + string(v)
			}
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
