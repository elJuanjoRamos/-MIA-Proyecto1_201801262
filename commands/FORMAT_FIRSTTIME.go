package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	STRUCTURES "../structures"
)

func MakeADirFirsTime(path string, id string) {
	//Voy a buscar la particion montada
	if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
		//se obtiene la particion montada
		var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go
		//Se abre el archivo para modificarlo
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			//OBTENGO EL SUPER BOOT
			var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA AQUI MISMO

			//Nos situamos al inicio del bitmap de arbol
			file.Seek(sb.SB_ap_bitmap_tree_dir, 0)

			//Escribimos un 1 en el bitmap
			var unit int8 = '1'
			s1 := &unit
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, s1)
			escribirBytes(file, binario3.Bytes())

			//NOS SITUAMOS AL INICIO DEL ARBOL DE DIRECOTORIO
			file.Seek(sb.SB_ap_tree_dir, 0)

			//CREAMOS EL ARBOL DE DIRECTORIO INICIAL "/"
			var tree = STRUCTURES.ARBOLVIRTUALDIR{
				Avd_ap_detalle_directorio:       sb.SB_ap_detail_dir,
				Avd_ap_arbol_virtual_directorio: -1,
			}
			var time = time.Now()
			copy(tree.Avd_fecha_creacion[:], time.Format("2006-01-02 15:04:05"))
			copy(tree.Avd_nombre_directorio[:], CorregirName(path))
			copy(tree.Avd_proper[:], CorregirName("root"))
			tree.Avd_ap_array_subdirectorios[0] = -1
			tree.Avd_ap_array_subdirectorios[1] = -1
			tree.Avd_ap_array_subdirectorios[2] = -1
			tree.Avd_ap_array_subdirectorios[3] = -1
			tree.Avd_ap_array_subdirectorios[4] = -1
			tree.Avd_ap_array_subdirectorios[5] = -1
			//ESCRIBIMOS EL ARBOL AL INICIO
			tree1 := &tree
			var binario4 bytes.Buffer
			binary.Write(&binario4, binary.BigEndian, tree1)
			escribirBytes(file, binario4.Bytes())

			//NOS SITUAMOS AL INICIO DEL BITMAP DE DETALLE DIRECTORIO
			file.Seek(sb.SB_ap_bitmap_detail_dir, 0)
			//Escribimos un 1 en el bitmap
			escribirBytes(file, binario3.Bytes())

			//NOS SITUAMOS AL INICIO DEL DETALLE DIRECTORIO
			file.Seek(sb.SB_ap_detail_dir, 0)

			fmt.Println()

			//CREAMOS EL PRIMER DETALLE Y LLENAMOS EL PRIMER FILE
			var dirInit = STRUCTURES.DIRECTORYDETAIL{
				DD_file_ap_inodo:         sb.SB_ap_bitmap_table_inode, //Como es la primera vez que lo creamos, el apuntador inodo de ese detalle, el el inicio del inodo
				DD_file_lleno:            true,
				DD_ap_detalle_directorio: 0,
			}
			copy(dirInit.DD_file_date_creacion[:], time.Format("2006-01-02 15:04:05"))
			copy(dirInit.DD_file_date_modificacion[:], time.Format("2006-01-02 15:04:05"))
			copy(dirInit.DD_file_nombre[:], "users.txt")

			//ESCRIBIMOS EL PRIMER DETALLE
			dirInit1 := &dirInit
			var binario5 bytes.Buffer
			binary.Write(&binario5, binary.BigEndian, dirInit1)
			escribirBytes(file, binario5.Bytes())

			//NOS SITUAMOS EN EL BITMAP DE TABLA INODO

			file.Seek(sb.SB_ap_bitmap_table_inode, 0)
			//Escribimos un 1 en el bitmap
			escribirBytes(file, binario3.Bytes())

			//NOS SITUAMOS AL INICIO DEL BLOQUE DE INODOS
			file.Seek(sb.SB_ap_table_inode, 0)

			var inodo = STRUCTURES.TABLEINODE{
				I_count_inodo:             1,
				I_size_archivo:            33,
				I_count_bloques_asignados: 2,
				I_ap_indirecto:            -1,
				I_id_proper:               1,
			}
			inodo.I_array_bloques[0] = sb.SB_ap_blocks                                                //bit de inicio del bloque de bloques
			inodo.I_array_bloques[1] = sb.SB_ap_blocks + int64(unsafe.Sizeof(STRUCTURES.DATABLOCK{})) //bit de inicio del bloque de bloques + size del bloque
			inodo.I_array_bloques[2] = 0
			inodo.I_array_bloques[3] = 0

			//ESCRIBIMOS EL PRIMER INODO
			inodo1 := &inodo
			var binarioInodo bytes.Buffer
			binary.Write(&binarioInodo, binary.BigEndian, inodo1)
			escribirBytes(file, binarioInodo.Bytes())

			//NOS SITUAMOS AL INICIO DEL BITMAP DE BLOQUES
			//ESCRIBIMOS UN UNO
			file.Seek(sb.SB_ap_bitmap_blocks, 0)
			escribirBytes(file, binario3.Bytes())

			file.Seek(sb.SB_ap_bitmap_blocks+1, 0)
			escribirBytes(file, binario3.Bytes())

			//CREAMOS LOS BLOQUES
			var block = STRUCTURES.DATABLOCK{}
			copy(block.DB_data[:], "1,G,root\n1,U,root,root,20")
			var block1 = STRUCTURES.DATABLOCK{}
			copy(block1.DB_data[:], "1801262\n")

			//NOS SITUAMOS EN EL INICIO DEL BLOQUE

			file.Seek(sb.SB_ap_blocks, 0)

			//ESCRIBIMOS LOS BLOQUES
			block11 := &block
			var binario7 bytes.Buffer
			binary.Write(&binario7, binary.BigEndian, block11)
			escribirBytes(file, binario7.Bytes())

			file.Seek(sb.SB_ap_blocks+int64(unsafe.Sizeof(STRUCTURES.DATABLOCK{})), 0)

			block12 := &block1
			var binario8 bytes.Buffer
			binary.Write(&binario8, binary.BigEndian, block12)
			escribirBytes(file, binario8.Bytes())

		}

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}

func GetSuperBoot(startPartition int64, file *os.File) STRUCTURES.SUPERBOOT {
	//Declaramos variable de tipo SUPERBOOT
	sb := STRUCTURES.SUPERBOOT{}
	//Obtenemos el tamanio del Super boot
	var sbSize int = int(unsafe.Sizeof(sb))

	//Nos situamos al inicio de la particion
	file.Seek(startPartition, 0)
	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, sbSize)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err := binary.Read(buffer, binary.BigEndian, &sb)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return sb
}

func CorregirName(texto string) string {
	if len(texto) <= 20 {
		var temp = len(texto)
		for i := 0; i < 20-temp; i++ {
			texto = texto + "+"
		}
	}
	return texto
}
