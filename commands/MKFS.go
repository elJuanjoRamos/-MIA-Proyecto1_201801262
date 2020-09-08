package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	STRUCTURES "../structures"
)

/*Variables Globales*/
/*var idUser = 1
var idGroup = 1*/

//funcion que se encarga de formatear la particion
func MKFSFormatPartition(id string, types string) {
	//Se crea el directorio que contiene los txt de usuarios y grupos
	//FUNCTION.CreateADirectory(FUNCTION.RootDir()+"/reports/userfiles", 0777)

	//primero voy a buscar la particion dentro de las particiones montadas
	if SearchPartitionById(id) { //esta funcion se encuenetra en commands/Moun_Umount.go
		//se obtiene la particion montada
		var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go
		//Se abre el archivo para irlo a limpiar
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {

			//Se envia a formatear la particion
			FastAndFullPartition(partition.Mount_part, file, strings.ToLower(types), filepath.Base(partition.Mount_path))

		}

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}

func FastAndFullPartition(partition STRUCTURES.PARTITION, file *os.File, types string, fileName string) {
	//Formatea la particion segun lo que venga

	//Si es full, se llena de ceros toda la particion
	if types == "full" {
		for i := partition.Part_start; i < partition.Part_end; i++ {
			var init int8 = 'p'
			o := &init
			file.Seek(i, 0)
			var binarioTemp bytes.Buffer
			binary.Write(&binarioTemp, binary.BigEndian, o)
			escribirBytes(file, binarioTemp.Bytes())
		}
	}

	var sizeArbolDir = int64(unsafe.Sizeof(STRUCTURES.ARBOLVIRTUALDIR{}))
	var sizeDetalleDir = int64(unsafe.Sizeof(STRUCTURES.DIRECTORYDETAIL{}))
	var sizeInodo = int64(unsafe.Sizeof(STRUCTURES.TABLEINODE{}))
	var sizeBloque = int64(unsafe.Sizeof(STRUCTURES.DATABLOCK{}))
	var sizeLog = int64(unsafe.Sizeof(STRUCTURES.LOG{}))
	var sizeSuperB = int64(unsafe.Sizeof(STRUCTURES.SUPERBOOT{}))

	var nEstructuras = (partition.Part_size - (2 * int64(sizeSuperB))) / (27 + int64(sizeArbolDir) + int64(sizeDetalleDir) + (5*int64(sizeInodo) + (20 * int64(sizeBloque)) + int64(sizeLog)))
	var cantidadAVD = nEstructuras
	var cantidadDD = nEstructuras
	var cantidadInodos = 5 * nEstructuras
	var cantidadBloques = 4 * cantidadInodos //20*nEstructuras
	var cantidadBitacoras = nEstructuras

	Inicio_bitmapAVD := partition.Part_start + int64(sizeSuperB)
	Inicio_AVD := Inicio_bitmapAVD + cantidadAVD
	Inicio_bitmapDD := Inicio_AVD + (int64(sizeArbolDir) * int64(cantidadAVD))
	Inicio_DD := Inicio_bitmapDD + int64(cantidadDD)
	Inicio_bitmapInodo := Inicio_DD + (int64(sizeDetalleDir) * int64(cantidadDD))
	Inicio_Inodos := Inicio_bitmapInodo + int64(cantidadInodos)
	Inicio_bitmapBloque := Inicio_Inodos + (int64(sizeInodo) * int64(cantidadInodos))
	Inicio_Bloque := Inicio_bitmapBloque + int64(cantidadBloques)
	Inicio_Bitacora := Inicio_Bloque + (int64(sizeBloque) * int64(cantidadBloques))
	Inicio_SBRespaldo := Inicio_Bitacora + (int64(sizeLog) * int64(cantidadBitacoras))

	//ESCRIBIMOS EL superboot AL INICIO
	EscribirSuperBlock(partition.Part_size, partition.Part_start, file, sizeSuperB, cantidadAVD, sizeArbolDir, cantidadDD,
		sizeDetalleDir, cantidadInodos, sizeInodo, cantidadBloques, sizeBloque, fileName)
	//ESCRIBIMOS EL BITMAP DEL ARBOL DE DIRECTORIO
	Escribir(Inicio_bitmapAVD, Inicio_AVD, file, '0')
	//ESCRIBIMOS EL ARBOL DE DIRECTORIO
	Escribir(Inicio_AVD, Inicio_bitmapDD, file, 'a')
	//ESCRIBIMOS EL BITMAP DE DETALLE DE DIRECTORIO
	Escribir(Inicio_bitmapDD, Inicio_DD, file, '0')
	//ESCRIBIMOS EL  DETALLE DE DIRECTORIO
	Escribir(Inicio_DD, Inicio_bitmapInodo, file, 'b')
	//ESCRIBIMOS EL BITMAP DEL INODO
	Escribir(Inicio_bitmapInodo, Inicio_Inodos, file, '0')
	//ESCRIBIMOS EL  INODO
	Escribir(Inicio_Inodos, Inicio_bitmapBloque, file, 'd')
	//ESCRIBIMOS EL BITMAP DEL BLOQUE
	Escribir(Inicio_bitmapBloque, Inicio_Bloque, file, '0')
	//ESCRIBIMOS EL  BLOQUE
	Escribir(Inicio_Bloque, Inicio_Bitacora, file, 'e')
	//ESCRIBIMOS LA BITACORA
	Escribir(Inicio_Bitacora, Inicio_SBRespaldo, file, 'f')
	//ESCRIBIMOS EL superboot AL FINAL
	EscribirSuperBlock(partition.Part_size, Inicio_SBRespaldo, file, sizeSuperB, cantidadAVD, sizeArbolDir,
		cantidadDD, sizeDetalleDir, cantidadInodos, sizeInodo, cantidadBloques, sizeBloque, fileName)

}

func Escribir(inicio int64, fin int64, file *os.File, data int8) {
	for i := inicio; i < fin; i++ {
		var init int8 = data
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escribirBytes(file, binarioTemp.Bytes())
	}

}

func EscribirSuperBlock(partitionSize int64, inicio int64, file *os.File, sizeSuperB, cantidadAVD, sizeArbolDir, cantidadDD, sizeDetalleDir, cantidadInodos, sizeInodo, cantidadBloques,
	sizeBloque int64, filename string) {
	//Se escribe el superboot
	file.Seek(inicio, 0)
	var superboot = STRUCTURES.SUPERBOOT{
		SB_magic_num:              201801262,
		SB_ap_bitmap_tree_dir:     inicio + sizeSuperB,
		SB_AVD_count:              cantidadAVD * int64(sizeArbolDir),
		SB_AVD_details_count:      cantidadAVD * int64(sizeDetalleDir),
		SB_Inodes_count:           5 * int64(sizeInodo) * cantidadAVD,
		SB_blocks_count:           20 * int64(sizeBloque) * cantidadAVD,
		SB_AVD_free:               cantidadAVD * int64(sizeArbolDir),
		SB_Inodes_free:            5 * int64(sizeInodo) * cantidadAVD,
		SB_blocks_free:            20 * int64(sizeBloque) * cantidadAVD,
		SB_mount_count:            1,
		SB_size_struct_tree_dir:   sizeArbolDir,
		SB_size_struct_detail_dir: sizeDetalleDir,
		SB_size_struct_inodo:      sizeInodo,
		SB_size_struct_block:      sizeBloque,
	}

	/*fmt.Println("Comprovacion")
	fmt.Println("super block 		 ", sizeSuperB)
	fmt.Println("bitmap arbol vitual 	 ", cantidadAVD)
	fmt.Println("arbol vitual 		 ", cantidadAVD*int64(sizeArbolDir))
	fmt.Println("bitmap detalle dir  	 ", cantidadAVD)
	fmt.Println("detalle dir 		 ", cantidadAVD*int64(sizeDetalleDir))
	fmt.Println("bitmap inodo 		 ", 5*cantidadAVD)
	fmt.Println("inodo 			 ", 5*int64(sizeInodo)*cantidadAVD)
	fmt.Println("bitmap bloque 		 ", 20*cantidadAVD)
	fmt.Println("bloque 			 ", 20*int64(sizeBloque)*cantidadAVD)
	fmt.Println("super block 		 ", sizeSuperB)*/

	//LE COLOCAMOS EL NOMBRE DEL DISCO Y LAS FECHAS DE MONTAJE Y CREACION
	copy(superboot.SB_hd_name[:], filename)
	var time = time.Now()
	copy(superboot.SB_date_lstmount[:], time.Format("2006-01-02 15:04:05"))
	copy(superboot.SB_date[:], time.Format("2006-01-02 15:04:05"))

	superboot.SB_ap_tree_dir = superboot.SB_ap_bitmap_tree_dir + cantidadAVD
	superboot.SB_ap_bitmap_detail_dir = superboot.SB_ap_tree_dir + (sizeArbolDir * cantidadAVD)
	superboot.SB_ap_detail_dir = superboot.SB_ap_bitmap_detail_dir + cantidadDD
	superboot.SB_ap_bitmap_table_inode = superboot.SB_ap_detail_dir + (sizeDetalleDir * cantidadDD)
	superboot.SB_ap_table_inode = superboot.SB_ap_bitmap_table_inode + cantidadInodos
	superboot.SB_ap_bitmap_blocks = superboot.SB_ap_table_inode + (sizeInodo * cantidadInodos)
	superboot.SB_ap_blocks = superboot.SB_ap_bitmap_blocks + cantidadBloques
	superboot.SB_ap_log = superboot.SB_ap_blocks + (sizeBloque * cantidadBloques)

	superboot.SB_first_free_bit_tree_dir = superboot.SB_ap_bitmap_detail_dir
	superboot.SB_first_free_bit_detail_dir = superboot.SB_ap_bitmap_detail_dir
	superboot.SB_first_free_bit_table_dir = superboot.SB_ap_bitmap_table_inode
	superboot.SB_first_free_bit_block = superboot.SB_ap_bitmap_blocks
	superboot.SB_free_space = partitionSize - (sizeSuperB + cantidadAVD + cantidadAVD*int64(sizeArbolDir) + cantidadAVD + cantidadAVD*int64(sizeDetalleDir) + 5*cantidadAVD + 5*int64(sizeInodo)*cantidadAVD + 20*cantidadAVD + 20*int64(sizeBloque)*cantidadAVD + sizeSuperB)

	var p = superboot
	a := &p
	var binarioP bytes.Buffer
	binary.Write(&binarioP, binary.BigEndian, a)
	file.Write(binarioP.Bytes())
}
