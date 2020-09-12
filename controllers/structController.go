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

			//Verifico si es primera creacion de logica, eso quiere ecir que el arreglo de ebr va a tener 1 y el arreglo de particiones ninguno
			if len(particionExtendidaAuxiliar.Part_ebr) == 1 && len(particionExtendidaAuxiliar.Part_partition) == 0 {
				particionExtendidaAuxiliar.Part_ebr[0].Part_fit = partitionFit[0]
				copy(particionExtendidaAuxiliar.Part_ebr[0].Part_name[:], "EBR"+partitionName)

				//////////////////////////////////////// NUEVO 	////////////////////////////////////////////////////
				var espacioLibre = particionExtendidaAuxiliar.Part_size - particionExtendidaAuxiliar.Part_ebr[0].Part_size

				if partitionSize <= espacioLibre { //Verifico que la particion quepa
					var particionAux = CreatePartition(partitionType, partitionFit[0],
						particionExtendidaAuxiliar.Part_ebr[0].Part_end, partitionSize, partitionName)

					array[i].extendedPartition.Part_partition = append(array[i].extendedPartition.Part_partition, particionAux)

				} else {
					fmt.Println("=======================================================================================")
					fmt.Println("La particion logica que esta intentando crear, no cabe dentro la particion extendida")
					fmt.Println("Espacio Disponible: ", espacioLibre, " Tam de Particion: ", partitionSize)
					fmt.Println("=======================================================================================")
				}

				//////////////////////////////////////// TERMINA NUEVO 	////////////////////////////////////////////////////

				//Sinifica que el arreglo de ebr y particiones ya son iguales y
			} else {

				//SE MANIPULA EL EBR ANTERIOR
				//Se obtiene la longitud del ebr
				var lenEbr int = len(array[i].extendedPartition.Part_ebr)
				//Se obtiene la longitud de la paticion logica asociada al ebr
				var lenPart int = len(array[i].extendedPartition.Part_partition)

				//////////////////  NUEVO ///////////////
				var espacioLibre = particionExtendidaAuxiliar.Part_size //OBTENGO EL SIZE DE LA PARTICION

				for j := 0; j < lenEbr; j++ {
					espacioLibre = espacioLibre - array[i].extendedPartition.Part_ebr[j].Part_size
				}
				for j := 0; j < lenPart; j++ {
					espacioLibre = espacioLibre - array[i].extendedPartition.Part_partition[j].Part_size
				}
				///////////////	TERMINA NUEVO /////////////////////
				//Se le asigna al ebr anterior el puntero del siguiente ebr que sera, el bit siguiente a donde termina la particion logica
				//asociada a ese ebr
				array[i].extendedPartition.Part_ebr[lenEbr-1].Part_next = array[i].extendedPartition.Part_partition[lenPart-1].Part_end

				if (int64(unsafe.Sizeof(STRUCTURES.EBR{})) + partitionSize) <= espacioLibre { //VERIFICO QUE QUEPA EL EBR Y LA PARTICION
					var ebrPartition = CreateEBR(1, partitionFit[0], array[i].extendedPartition.Part_partition[lenPart-1].Part_end,
						int64(unsafe.Sizeof(STRUCTURES.EBR{})), -1, "EBR"+partitionName)
					//Guardamos el ebr que creamos
					array[i].extendedPartition.Part_ebr = append(array[i].extendedPartition.Part_ebr, ebrPartition)

					//volvemos a obtener el len del arreglo de ebr
					lenEbr = len(array[i].extendedPartition.Part_ebr)

					var primPartition = CreatePartition("L", partitionFit[0], array[i].extendedPartition.Part_ebr[lenEbr-1].Part_end, partitionSize, partitionName)
					//Guaramos la particion
					array[i].extendedPartition.Part_partition = append(array[i].extendedPartition.Part_partition, primPartition)
				} else {
					fmt.Println("=======================================================================================")
					fmt.Println("La particion logica que esta intentando crear, no cabe dentro la particion extendida")
					fmt.Println("Espacio Disponible: ", espacioLibre, " Tam de Particion: ", partitionSize)
					fmt.Println("=======================================================================================")
				}
			}

			break
		}
	}
}

//Crea un ebr con los datos necesarios y los retorna
func CreateEBR(status int8, fit byte, start int64, size int64, next int64, name string) STRUCTURES.EBR {
	var ebr = STRUCTURES.EBR{
		Part_status: status,
		Part_fit:    fit,
		Part_start:  start,
		Part_size:   size,
		Part_next:   next,
	}
	ebr.Part_end = ebr.Part_start + ebr.Part_size
	copy(ebr.Part_name[:], CorregirNombre(name))
	return ebr
}

//Arma una particion con los datos necesarios y la retorna
func CreatePartition(types string, fit byte, end int64, size int64, name string) STRUCTURES.PARTITION {

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

	return part

}

//Busca una particion dentro del arreglo de particiones
func SearchPartition(diskName string, partitionName string) bool {

	for i := 0; i < len(array); i++ {
		var aux = array[i]
		if aux.diskName == diskName {

			for j := 0; j < len(aux.extendedPartition.Part_partition); j++ {
				var partition = aux.extendedPartition.Part_partition[j]
				if AssemblePartName(partition.Part_name) == partitionName {
					return true
				}
			}
			break
		}
	}
	return false
}

func GetLogicPartition(diskName string, partitionName string) STRUCTURES.PARTITION {

	var partition = STRUCTURES.PARTITION{}
	for i := 0; i < len(array); i++ {
		var aux = array[i]
		if aux.diskName == diskName {

			for j := 0; j < len(aux.extendedPartition.Part_partition); j++ {
				partition = aux.extendedPartition.Part_partition[j]
				if AssemblePartName(partition.Part_name) == partitionName {
					return partition
				}
			}
			break
		}
	}
	return partition
}

//Retorna un string de nombre que antes eran bytes
func AssemblePartName(name [16]byte) string {
	var s string = ""
	for _, v := range name {
		if v != 0 {
			s = s + string(v)
		}
	}
	return s
}

func GetExtendedPartition(diskName string) (STRUCTURES.EXTENDED, bool) {
	for i := 0; i < len(array); i++ {
		if array[i].diskName == diskName {
			return array[i].extendedPartition, true
		}
	}
	return STRUCTURES.EXTENDED{}, false
}

func CorregirNombre(name string) string {

	if len(name) < 16 {
		var cont = 16 - len(name)
		for i := 0; i < cont; i++ {
			name = name + "+"
		}
	}
	return name
}

func RemoveExtendida(diskName string) {
	for i := 0; i < len(array); i++ {
		if array[i].diskName == diskName {
			array[i].extendedPartition = STRUCTURES.EXTENDED{}
		}
	}
}

func RemoveLogica(diskName string, partName string) (STRUCTURES.PARTITION, bool) {
	for i := 0; i < len(array); i++ {
		if array[i].diskName == diskName {
			var particionExtendida = array[i].extendedPartition

			//RECORRO LAS PARTICIONSE LOGICA PARA BUSCAR LA PRTICION
			for j := 0; j < len(particionExtendida.Part_partition); j++ {
				var particion = particionExtendida.Part_partition[j]

				if GetString(particion.Part_name) == partName {

					array[i].extendedPartition.Part_partition[j] = STRUCTURES.PARTITION{Part_size: particion.Part_size}
					array[i].extendedPartition.Part_ebr[j] = STRUCTURES.EBR{Part_size: int64(unsafe.Sizeof(STRUCTURES.EBR{}))}

					return particion, true
				}

			}
		}
	}
	return STRUCTURES.PARTITION{}, false
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
