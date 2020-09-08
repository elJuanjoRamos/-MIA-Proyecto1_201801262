package commands

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"

	STRUCTURES "../structures"
)

func MakeADir(path string, id string, p string) {
	//primero voy a buscar la particion dentro de las particiones montadas
	if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
		//se obtiene la particion montada
		//var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go

		//if CONTROLLER.IsLogged() {

		if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
			//se obtiene la particion montada
			var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go
			//Se abre el archivo para modificarlo
			file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
			defer file.Close()
			if err != nil {
				fmt.Println(err)
			} else {

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

				//Empezamos a leer en el bitmap de arbol de directorios
				//desde   el inicio del bitmap	hasta  el inicio del arbol de directorio
				var bitmapArbolSize = int(sb.SB_ap_tree_dir) - int(sb.SB_ap_bitmap_tree_dir)

				r4 := bufio.NewReader(file)
				b4, err := r4.Peek(bitmapArbolSize)
				if err != nil {
					fmt.Println(err)
				}
				for i := 0; i < len(b4); i++ {

					fmt.Println(string(b4[i]))
				}

				//Se envia a formatear la particion
				//FastAndFullPartition(partition.Mount_part, file, strings.ToLower(types), filepath.Base(partition.Mount_path))

			}

		} else {
			fmt.Println("No existe una particion con el id " + id)
		}

		/*} else {
			fmt.Println("No hay ningun usuario logueado")
		}*/

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}
