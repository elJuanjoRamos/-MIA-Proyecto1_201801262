package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"

	STRUCTURES "../structures"
)

func MakeAFileInLogicalDisk(path string, id string, p string, size int64, cont string) {
	//primero voy a buscar la particion dentro de las particiones montadas
	if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
		//se obtiene la particion montada
		var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go
		//Se abre el archivo para modificarlo
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {

			//vuelvo un array de paths para ver la longitud d

			//Declaramos variable de tipo SUPERBOOT
			sb := STRUCTURES.SUPERBOOT{}
			//Obtenemos el tamanio del Super boot
			var sbSize int = int(unsafe.Sizeof(sb))

			//Nos situamos al inicio de la particion
			file.Seek(partition.Mount_part.Part_start, 0)
			//Lee la cantidad de <size> bytes del archivo
			data := leerBytes(file, sbSize)
			//Convierte la data en un buffer,necesario para
			//decodificar binario
			buffer := bytes.NewBuffer(data)

			//Decodificamos y guardamos en la variable m
			err = binary.Read(buffer, binary.BigEndian, &sb)
			if err != nil {
				log.Fatal("binary.Read failed", err)
			}

			fmt.Println("SUPERBOOT")

			//Nos situamos al inicio del bitmap de arbol
			file.Seek(sb.SB_ap_bitmap_tree_dir, 0)

			//Escribimos un 1 en el bitmap
			var unit int8 = '1'
			s1 := &unit
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, s1)
			escribirBytes(file, binario3.Bytes())

			//var bitmapArbolSize = int(sb.SB_ap_tree_dir) - int(sb.SB_ap_bitmap_tree_dir)

			//Empezamos a leer en el bitmap de arbol de directorios
			//desde   el inicio del bitmap	hasta  el inicio del arbol de directorio

			/*r4 := bufio.NewReader(file)
			b4, err := r4.Peek(bitmapArbolSize)
			if err != nil {
				fmt.Println(err)
			}
			for i := 0; i < len(b4); i++ {
				fmt.Println(string(b4[i]))
			}*/

			//Se envia a formatear la particion
			//FastAndFullPartition(partition.Mount_part, file, strings.ToLower(types), filepath.Base(partition.Mount_path))

		}

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}

func MakeAFileInLogicalDiskFirstTime(path string, id string, p string, size int64, cont string) {
	//primero voy a buscar la particion dentro de las particiones montadas
	if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
		//se obtiene la particion montada
		var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go
		//Se abre el archivo para modificarlo
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {

			//OBTENGO EL SUPER BLOQUE

			//Declaramos variable de tipo SUPERBOOT
			sb := GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA EN MKDIR.GO

			//CREAMOS LOS BLOQUES

			if len(cont) <= 25 {

				var block = STRUCTURES.DATABLOCK{}
				copy(block.DB_data[:], cont)
				//NOS SITUAMOS AL INICIO DEL BLOQUE DE DATOS
				file.Seek(sb.SB_ap_blocks, 0)
				block1 := &block
				var binario4 bytes.Buffer
				binary.Write(&binario4, binary.BigEndian, block1)
				escribirBytes(file, binario4.Bytes())
			} else {
				/*var str = ""
				var tope int = 0
				for i := 0; i < 25; i++ {
					str = str + string(cont[i])
					tope = i
				}
				fmt.Println(str)
				if (len(cont) - len(str)) <= 25 {

				}*/
			}

		}

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}

func evaluar() {

}
