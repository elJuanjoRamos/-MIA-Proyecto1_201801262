package commandExecute

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"

	STRUCTURES "../structures"
)

func FormatDisk(path string, partitionSize int64, partitionName string, partitionType byte, partitionFit byte) {

	//BUSCAMOS EL ARCHIVO
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	} else {
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

		//Se revisan la cantidad de particiones disponibles en
		//el disco
		if m.Mbr_count > 0 {

			switch m.Mbr_count {
			case 4: //Se mete en la particion 1
				m.Mbr_partition_1 = STRUCTURES.PARTITION{
					Part_status: 'V',
					Part_type:   partitionType,
					Part_fit:    partitionFit,
					Part_start:  1,
					Part_size:   partitionSize}
				copy(m.Mbr_partition_1.Part_name[:], partitionName)
				//Se disminuye el contador de particiones
				m.Mbr_count = 3
				break
			case 3:
				m.Mbr_partition_2 = STRUCTURES.PARTITION{
					Part_status: 'V',
					Part_type:   partitionType,
					Part_fit:    partitionFit,
					Part_start:  1,
					Part_size:   partitionSize}
				copy(m.Mbr_partition_2.Part_name[:], partitionName)
				//Se disminuye el contador de particiones
				m.Mbr_count = 2
				break
			case 2:
				m.Mbr_partition_3 = STRUCTURES.PARTITION{
					Part_status: 'V',
					Part_type:   partitionType,
					Part_fit:    partitionFit,
					Part_start:  1,
					Part_size:   partitionSize}
				copy(m.Mbr_partition_3.Part_name[:], partitionName)
				//Se disminuye el contador de particiones
				m.Mbr_count = 1
				break
			case 1:
				m.Mbr_partition_4 = STRUCTURES.PARTITION{
					Part_status: 'V',
					Part_type:   partitionType,
					Part_fit:    partitionFit,
					Part_start:  1,
					Part_size:   partitionSize}
				copy(m.Mbr_partition_4.Part_name[:], partitionName)
				//Se disminuye el contador de particiones
				m.Mbr_count = 0
				break
			}

			//Se situa en la posicion 0,0 del archivo
			file.Seek(0, 0)
			//Escribe el mbr con particiones en el archivo
			s1 := &m
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, s1)
			escribirBytes(file, binario3.Bytes())

		} else {
			fmt.Println("No se puede escribir mas particiones en el disco")
		}

	}

}
