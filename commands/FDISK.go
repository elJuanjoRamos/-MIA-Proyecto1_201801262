package commands

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	CONTROLLER "../controllers"
	FUNCTION "../functions"
	STRUCTURES "../structures"
)

func FormatDisk(path string, partitionSize int64, partitionName string, partitionType string, partitionFit string) {
	diskName := filepath.Base(path)
	//BUSCAMOS EL ARCHIVO
	if FUNCTION.IfExistFile(path) {

		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {

			//Obtengo el size del archivo
			//var fileSize = FUNCTION.FileSize(path)

			//Declaramos variable de tipo mbr
			m := STRUCTURES.MBR{}
			//Obtenemos el tamanio del mbr
			var mbrSize int = int(unsafe.Sizeof(m))

			//Lee la cantidad de <size> bytes del archivo
			data := leerBytes(file, mbrSize)
			//Convierte la data en un buffer,necesario para
			//decodificar binario
			buffer := bytes.NewBuffer(data)

			//Decodificamos y guardamos en la variable m
			err = binary.Read(buffer, binary.BigEndian, &m)
			if err != nil {
				log.Fatal("binary.Read failed", err)
			}

			//Se revisan la cantidad de particiones disponibles en el disco
			//if m.Mbr_count > 0 {
			//REVISO EL TIPO DE PARTICION

			//Si la particion es primaria, simplemente se manda a crear
			if partitionType == "P" && m.Mbr_count > 0 {
				m = CreatePartition(m, partitionType, partitionFit, mbrSize, partitionSize, partitionName, "")
				m.Mbr_count = m.Mbr_count - 1
			} else if partitionType == "E" {
				//Si el disco aun no tiene particiones extendidas, mando a crear la particion extendida y
				//Cambio la bandera a 1, eso quiere decir que ya tiene particiones extendidas dentro
				if m.Mbr_Ext == 0 {

					filename := filepath.Base(path)
					m = CreatePartition(m, partitionType, partitionFit, mbrSize, partitionSize, partitionName, filename)
					m.Mbr_Ext = 1

				} else {
					fmt.Println("No se pueden crear mas particiones extendidas en el disco")
				}

			} else if partitionType == "L" {
				//VALIDO QUE EL DISCO TENGA UNA PARTICION EXTENDIDA
				if m.Mbr_Ext == 1 {
					CONTROLLER.AddLogicPartition(partitionType, partitionFit, partitionSize, partitionName, filepath.Base(path))

				} else {
					fmt.Println("No se ha creado una particion extendida para el disco")
				}
			} else {
				fmt.Println("No se pueden crear mas particiones primarias en el disco")
			}

			//Se situa en la posicion 0,0 del archivo
			file.Seek(0, 0)
			//Escribe el mbr con particiones en el archivo
			s1 := &m
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, s1)
			escribirBytes(file, binario3.Bytes())
			CreateMBRReport(m)
			CreateDiskReport(m, diskName)

			/*} else {
				fmt.Println("No se puede escribir mas particiones en el disco")
			}*/

		}
	} else {
		fmt.Println("El DISCO con nombre '" + diskName + "' no  existe")
	}

}

func SendToFull(path string) {
	filename := filepath.Base(path)
	CONTROLLER.FullEBR(filename, path)

}

func CreatePartition(m STRUCTURES.MBR, partitionType string, partitionFit string, mbrSize int, partitionSize int64, partitionName string, filename string) STRUCTURES.MBR {
	//Verifica si la paticion 1 esta vacia
	if m.Mbr_partition_1.Part_isEmpty == 0 {
		m.Mbr_partition_1 = AssemblePartition(partitionType, partitionFit[0], int64(mbrSize), partitionSize, partitionName, filename)
		//Verifica si la paticion 2 esta vacia
	} else if m.Mbr_partition_2.Part_isEmpty == 0 {
		m.Mbr_partition_2 = AssemblePartition(partitionType, partitionFit[0], m.Mbr_partition_1.Part_end, partitionSize, partitionName, filename)
		//Verifica si la paticion 3 esta vacia
	} else if m.Mbr_partition_3.Part_isEmpty == 0 {
		m.Mbr_partition_3 = AssemblePartition(partitionType, partitionFit[0], m.Mbr_partition_2.Part_end, partitionSize, partitionName, filename)
		//Verifica si la paticion 4 esta vacia
	} else if m.Mbr_partition_4.Part_isEmpty == 0 {
		m.Mbr_partition_4 = AssemblePartition(partitionType, partitionFit[0], m.Mbr_partition_3.Part_end, partitionSize, partitionName, filename)
	}
	return m
}

//Arma la particion con la data necesaria
func AssemblePartition(types string, fit byte, end int64, size int64, name string, filename string) STRUCTURES.PARTITION {

	var part = STRUCTURES.PARTITION{
		Part_status:  1,
		Part_type:    types[0],
		Part_fit:     fit,
		Part_start:   end,
		Part_size:    size,
		Part_isEmpty: 1,
	}
	copy(part.Part_name[:], CorregirNombre(name))
	part.Part_end = part.Part_start + size

	if strings.ToLower(types) == "e" {
		CONTROLLER.AddExtended(1, fit, part.Part_start, part.Part_end, size, name, filename)
	}
	return part

}

func CorregirNombre(nombre string) string {
	if len(nombre) < 16 {
		var newLen = 16 - len(nombre)
		for i := 0; i < newLen; i++ {
			nombre = nombre + "+"
		}
	}
	return nombre
}

func GetString(str [16]byte) string {
	var cadena = ""
	for i := 0; i < len(str); i++ {
		if string(str[i]) == "+" {
			break
		}
		cadena = cadena + string(str[i])
	}
	return cadena
}

////////////ELIMINAR PARTICION

func DeletePartition(path string, partitionName string, delete string) {
	diskName := filepath.Base(path)

	//BUSCAMOS EL ARCHIVO
	if FUNCTION.IfExistFile(path) {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			//Obtengo el size del archivo
			//var fileSize = FUNCTION.FileSize(path)

			//Declaramos variable de tipo mbr
			m := STRUCTURES.MBR{}
			//Obtenemos el tamanio del mbr
			var mbrSize int = int(unsafe.Sizeof(m))

			//Lee la cantidad de <size> bytes del archivo
			data := leerBytes(file, mbrSize)
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

			//BUSCO EN LAS PARTICIONES PRIMARIAS LA PARTICION A ELIMINAR
			var flag = false
			for i := 0; i < 4; i++ {
				var particion = PartitionArray[i]

				if string(particion.Part_type) == "P" {
					if GetString(particion.Part_name) == partitionName {
						if ConfirmacionELiminacion(partitionName) {
							Format(particion, delete, file)
							var sizeTemp = particion.Part_size
							m.Mbr_count = m.Mbr_count + 1
							PartitionArray[i] = STRUCTURES.PARTITION{Part_size: sizeTemp}
							flag = true
						}
						break
					}
				}
			}

			var encontrado = false
			if !flag { //SI ENTRA A BANDERA SIGNIFICA QUE PUEDE SER LA EXTENDIDA LA QUE VOY A ELMINAR
				for i := 0; i < 4; i++ {
					var particion = PartitionArray[i]
					if string(particion.Part_type) == "E" {
						archivo := filepath.Base(path)

						if GetString(particion.Part_name) == partitionName { //SI ES LA PARTICION EXTENDIDA, LA ELIMINO Y ELIMINO SUS LOGIACAS

							if ConfirmacionELiminacion(partitionName) {
								m.Mbr_Ext = 0
								CONTROLLER.RemoveExtendida(archivo)
								var sizeTemp = particion.Part_size
								PartitionArray[i] = STRUCTURES.PARTITION{Part_size: sizeTemp}
								encontrado = true
							}

						} else { // SI NO ES LA EXTENDIDA, PUEDE QUE SEA UNA LOGICA

							var particio, enc = CONTROLLER.RemoveLogica(archivo, partitionName)

							if enc {
								if ConfirmacionELiminacion(partitionName) {
									encontrado = true
									Format(particio, delete, file)
								}
							}
						}

						break
					}
				}
			}

			if encontrado {
				//SACO  LAS PARTICIONES DEL ARREGLO
				m.Mbr_partition_1 = PartitionArray[0]
				m.Mbr_partition_2 = PartitionArray[1]
				m.Mbr_partition_3 = PartitionArray[2]
				m.Mbr_partition_4 = PartitionArray[3]

				file.Seek(0, 0)
				//Escribe el mbr con particiones en el archivo
				s1 := &m
				var binario3 bytes.Buffer
				binary.Write(&binario3, binary.BigEndian, s1)
				escribirBytes(file, binario3.Bytes())
				CreateMBRReport(m)
				CreateDiskReport(m, diskName)

			} else {
				fmt.Println("La particion con nombre '" + partitionName + "' no exite en el disco")
			}

		}
	} else {
		fmt.Println("El DISCO con nombre '" + diskName + "' no  existe")
	}

}

func ConfirmacionELiminacion(partitionName string) bool {
	fmt.Println("=====================================================")
	fmt.Println("  Are you sure you want to delete " + partitionName + "?")
	fmt.Println("=====================================================")
	fmt.Println("")
	fmt.Print("Press Y/N: ")

	reader := bufio.NewReader(os.Stdin)
	comando, _ := reader.ReadString('\n')
	input := ""
	if runtime.GOOS == "windows" {
		input = strings.TrimRight(comando, "\r\n")
	} else {
		input = strings.TrimRight(comando, "\n")
	}
	if strings.TrimRight(input, "\n") == "Y" || strings.TrimRight(input, "\n") == "y" {

		fmt.Println("Partition '" + partitionName + "' successfully deleted")
		return true
	} else {
		fmt.Println("The partition '" + partitionName + "'  was not erased")
		return false
	}
}

func Format(partition STRUCTURES.PARTITION, types string, file *os.File) {
	//Si es full, se llena de ceros toda la particion
	var init int8 = 0
	if strings.ToLower(types) == "full" {
		init = '0'
	}
	for i := partition.Part_start; i < partition.Part_end; i++ {
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escribirBytes(file, binarioTemp.Bytes())
	}
}
