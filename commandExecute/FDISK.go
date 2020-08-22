package commandExecute

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"unsafe"

	STRUCTURES "../structures"
)

func FormatDisk(path string, partitionSize int64, partitionName string, partitionType byte, partitionFit byte) {

	//BUSCAMOS EL ARCHIVO
	file, err := os.Open(path)
	defer file.Close()
	if err != nil { //validar que no sea nulo.
		log.Fatal(err)

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
		}

		//Se imprimen los valores guardados en el struct
		fmt.Println("-----MBR DATA-----")
		fmt.Printf("SIZE: %s\nFECHA: %s\nSIGNATURE: %s\n", strconv.Itoa(int(m.Mbr_size)), m.Mbr_creation_date, strconv.Itoa(int(m.Mbr_disk_signature)))
		fmt.Println("----------")
		fmt.Println("-----PARTITION-----")

		s := reflect.ValueOf(&m.Mbr_partition_1).Elem()
		typeOfT := s.Type()

		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			fmt.Printf("%s = %v\n",
				typeOfT.Field(i).Name, f.Interface())
		}
		//Nos situamos en el inicio del archivo
		//file.Seek(0, 0)

		s1 := &m
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, s1)
		escribirBytes2(file, binario3.Bytes())

	}

}

//Método para escribir en un archivo.
func escribirBytes2(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}
