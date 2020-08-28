package controllers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"

	STRUCTURES "../structures"
)

type extendedDisk struct {
	diskName          string
	extendedPartition STRUCTURES.EXTENDED
}

var array = []extendedDisk{}

func AddNewExtendedDisk(name string) {
	var disk = extendedDisk{
		diskName: name,
	}

	array = append(array, disk)
}

func AddExtended(status int8, fit byte, start int64, end int64, size int64, name string, filename string) {

	for i := 0; i < len(array); i++ {
		var aux = array[i]
		if aux.diskName == filename {

			//Se crea un EBR vacio
			var ebr = CreateEBR(status, 'W', start, int64(unsafe.Sizeof(STRUCTURES.EBR{})), -1, "EBR")

			//SE CREA UNA PARTICION TIPO EXTENDIDA
			array[i].extendedPartition = STRUCTURES.EXTENDED{
				Part_status: status,
				Part_type:   'E',
				Part_fit:    fit,
				Part_start:  start,
				Part_end:    end,
				Part_size:   size,
			}
			copy(array[i].extendedPartition.Part_name[:], name)
			//Se guarda el ebr vacio en la particion
			array[i].extendedPartition.Part_ebr = append(array[i].extendedPartition.Part_ebr, ebr)
			break
		}
	}

	//Se guarda la particion dentro del arreglo de particiones
	//extendedArray = append(extendedArray, ext)
}
func FullEBR(filename string, path string) {

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	} else {

		for i := 0; i < len(array); i++ {
			if array[i].diskName == filename {
				var particionExtendidaAuxiliar = array[i].extendedPartition
				for i := 0; i < len(particionExtendidaAuxiliar.Part_ebr); i++ {
					var ebrAuxiliar = particionExtendidaAuxiliar.Part_ebr[i]
					file.Seek(ebrAuxiliar.Part_start, 0)
					s1 := &ebrAuxiliar
					var binario3 bytes.Buffer
					binary.Write(&binario3, binary.BigEndian, s1)
					escribirBytes(file, binario3.Bytes())
					break
				}
				if len(particionExtendidaAuxiliar.Part_partition) != 0 {
					for i := 0; i < len(particionExtendidaAuxiliar.Part_partition); i++ {
						var ebrAuxiliar = particionExtendidaAuxiliar.Part_partition[i]
						file.Seek(ebrAuxiliar.Part_start, 0)
						s1 := &ebrAuxiliar
						var binario3 bytes.Buffer
						binary.Write(&binario3, binary.BigEndian, s1)
						escribirBytes(file, binario3.Bytes())
						break
					}
				}
				break
			}
		}
	}

}
func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func AddLogicPartition(partitionType string, partitionFit string, partitionSize int64, partitionName string, filename string) {
	for i := 0; i < len(array); i++ {
		if array[i].diskName == filename {
			var particionExtendidaAuxiliar = array[i].extendedPartition

			if len(particionExtendidaAuxiliar.Part_ebr) == 1 && len(particionExtendidaAuxiliar.Part_partition) == 0 {
				particionExtendidaAuxiliar.Part_ebr[0].Part_fit = partitionFit[0]
				copy(particionExtendidaAuxiliar.Part_ebr[0].Part_name[:], "EBR"+partitionName)

				var particionAux = CreatePartition(partitionType, partitionFit[0],
					particionExtendidaAuxiliar.Part_ebr[0].Part_end, partitionSize, partitionName)

				array[i].extendedPartition.Part_partition = append(array[i].extendedPartition.Part_partition, particionAux)
			} else {

			}

			break
		}
	}
}

func CreateEBR(status int8, fit byte, start int64, size int64, next int64, name string) STRUCTURES.EBR {
	var ebr = STRUCTURES.EBR{
		Part_status: status,
		Part_fit:    fit,
		Part_start:  start,
		Part_size:   size,
		Part_next:   next,
	}
	ebr.Part_end = ebr.Part_start + ebr.Part_size
	copy(ebr.Part_name[:], name)

	return ebr
}

func CreatePartition(types string, fit byte, end int64, size int64, name string) STRUCTURES.PARTITION {

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

	return part

}
