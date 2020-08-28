package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"

	CONTROLLER "../controllers"
	FUNCTION "../functions"
	STRUCTURES "../structures"
)

func FormatDisk(path string, partitionSize int64, partitionName string, partitionType string, partitionFit string) {

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

			//variable que va contar el espacio libre

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
			if m.Mbr_count > 0 {
				//REVISO EL TIPO DE PARTICION

				//Si la particion es primaria, simplemente se manda a crear
				if partitionType == "P" {
					m = CreatePartition(m, partitionType, partitionFit, mbrSize, partitionSize, partitionName)
					FullPartition(m.Mbr_partition_1, file, 1)
					FullPartition(m.Mbr_partition_2, file, 2)
					//FullPartition(m.Mbr_partition_3, file, 3)
					FullPartition(m.Mbr_partition_4, file, 4)
					//Si la paticion es extendida
				} else if partitionType == "E" {
					//Si el disco aun no tiene particiones extendidas, mando a crear la particion extendida y
					//Cambio la bandera a 1, eso quiere decir que ya tiene particiones extendidas dentro
					if m.Mbr_Ext == 0 {
						m = CreatePartition(m, partitionType, partitionFit, mbrSize, partitionSize, partitionName)
						FullPartition(m.Mbr_partition_1, file, 1)
						FullPartition(m.Mbr_partition_2, file, 2)
						//FullPartition(m.Mbr_partition_3, file, 3)
						FullPartition(m.Mbr_partition_4, file, 4)
						m.Mbr_Ext = 1

					} else {
						fmt.Println("No se pueden crear mas particiones extendidas en el disco")
					}

				} else if partitionType == "L" {
					//VALIDO QUE EL DISCO TENGA UNA PARTICION EXTENDIDA
					if m.Mbr_Ext == 1 {
						CONTROLLER.AddLogicPartition(partitionType, partitionFit, partitionSize, partitionName)

					} else {
						fmt.Println("No se ha creado una particion extendida para el disco")
					}
				}

				//Se situa en la posicion 0,0 del archivo
				file.Seek(0, 0)
				//Escribe el mbr con particiones en el archivo
				s1 := &m
				var binario3 bytes.Buffer
				binary.Write(&binario3, binary.BigEndian, s1)
				escribirBytes(file, binario3.Bytes())

				//REPORTS.CreateMBRReport(m)
			} else {
				fmt.Println("No se puede escribir mas particiones en el disco")
			}

		}
	}
	if partitionType == "E" {
		CONTROLLER.FullEBR(path, partitionSize, partitionName)
	}
}

func CreatePartition(m STRUCTURES.MBR, partitionType string, partitionFit string, mbrSize int, partitionSize int64, partitionName string) STRUCTURES.MBR {
	//Verifica si la paticion 1 esta vacia
	if m.Mbr_partition_1.Part_isEmpty == 0 {
		m.Mbr_partition_1 = AssemblePartition(partitionType, partitionFit[0], int64(mbrSize), partitionSize, partitionName)
		m.Mbr_count = 3
		//Verifica si la paticion 2 esta vacia
	} else if m.Mbr_partition_2.Part_isEmpty == 0 {
		m.Mbr_partition_2 = AssemblePartition(partitionType, partitionFit[0], m.Mbr_partition_1.Part_end, partitionSize, partitionName)
		m.Mbr_count = 2
		//Verifica si la paticion 3 esta vacia
	} else if m.Mbr_partition_3.Part_isEmpty == 0 {
		m.Mbr_partition_3 = AssemblePartition(partitionType, partitionFit[0], m.Mbr_partition_2.Part_end, partitionSize, partitionName)
		m.Mbr_count = 1
		//Verifica si la paticion 4 esta vacia
	} else if m.Mbr_partition_4.Part_isEmpty == 0 {
		m.Mbr_partition_4 = AssemblePartition(partitionType, partitionFit[0], m.Mbr_partition_3.Part_end, partitionSize, partitionName)
		m.Mbr_count = 0
	}
	return m
}

//Arma la particion con la data necesaria
func AssemblePartition(types string, fit byte, end int64, size int64, name string) STRUCTURES.PARTITION {

	var part = STRUCTURES.PARTITION{
		Part_status:  1,
		Part_type:    types[0],
		Part_fit:     fit,
		Part_start:   end + 1,
		Part_size:    size,
		Part_isEmpty: 1,
	}
	copy(part.Part_name[:], name)
	part.Part_end = part.Part_start + size

	if strings.ToLower(types) == "e" {
		CONTROLLER.AddExtended(1, fit, part.Part_start, part.Part_end, size, name)
	}
	return part

}

func FullPartition(partition STRUCTURES.PARTITION, file *os.File, number int8) {
	for i := partition.Part_start; i < partition.Part_end+1; i++ {
		var init int8 = 'P' + number
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escribirBytes(file, binarioTemp.Bytes())
	}
}
