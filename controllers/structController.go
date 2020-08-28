package controllers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"

	FUNCTION "../functions"
	STRUCTURES "../structures"
)

var extendedArray = []STRUCTURES.EXTENDED{}

func AddExtended(status int8, fit byte, start int64, end int64, size int64, name string) {
	//Se crea un EBR vacio

	var ebr = CreateEBR(status, 'W', start, int64(unsafe.Sizeof(STRUCTURES.EBR{})), -1, "EBR")

	//SE CREA UNA PARTICION TIPO EXTENDIDA
	var ext = STRUCTURES.EXTENDED{
		Part_status: status,
		Part_type:   'E',
		Part_fit:    fit,
		Part_start:  start,
		Part_end:    end,
		Part_size:   size,
	}
	copy(ext.Part_name[:], name)
	//Se guarda el ebr vacio en la particion
	ext.Part_ebr = append(ext.Part_ebr, ebr)

	//Se guarda la particion dentro del arreglo de particiones
	extendedArray = append(extendedArray, ext)
}
func FullEBR(path string, partitionSize int64, partitionName string) {
	if FUNCTION.IfExistFile(path) {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {

			for i := 0; i < len(extendedArray); i++ {
				var aux = extendedArray[i]

				var s string
				for _, v := range aux.Part_name {
					if v != 0 {
						s = s + string(v)
					}
				}
				//BUSCA LA PARTICION EXTENDIDA DEL DISCO EN EL ARREGLO DE PARTICIONES
				if aux.Part_size == partitionSize && s == partitionName {
					fmt.Println()
					file.Seek(aux.Part_start, 0)
					s1 := &aux.Part_ebr[0]
					var binario3 bytes.Buffer
					binary.Write(&binario3, binary.BigEndian, s1)
					escribirBytes(file, binario3.Bytes())
					break
				} else {
					fmt.Println("no la encontro")
				}
			}

			//Se situa en la posicion 0,0 del archivo

		}
	}

}
func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func AddLogicPartition(partitionType string, partitionFit string, partitionSize int64, partitionName string) {

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
