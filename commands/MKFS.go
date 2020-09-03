package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"unsafe"

	FUNCTION "../functions"
	STRUCTURES "../structures"
)

/*Variables Globales*/
/*var idUser = 1
var idGroup = 1*/

//funcion que se encarga de formatear la particion
func MKFSFormatPartition(id string, types string) {
	//Se crea el directorio que contiene los txt de usuarios y grupos
	FUNCTION.CreateADirectory(FUNCTION.RootDir()+"/reports/userfiles", 0777)

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

			var pathFile = FUNCTION.RootDir() + "/reports/userfiles/" + "users_" + id + ".txt" //La funcion RootDir retorna la ruta del proyecto, se encuentra en special_functions.go
			// Mando a actualizar la particion montada, dando la ruta del archivo de usuarios, esto para no tener que leer en el disco la ruta
			UpdateUserTxt(id, pathFile)
			//Se crea el archivo txt con los usuario iniciales, cambiar carnet
			FUNCTION.CreateAFile(pathFile, "1,G,root\n1,U,root,root,201801262,root1,123")

			//Se envia a formatear la particion
			FastAndFullPartition(partition.Mount_part, file, strings.ToLower(types), pathFile)

		}

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}

func FastAndFullPartition(partition STRUCTURES.PARTITION, file *os.File, types string, pathFile string) {
	//Formatea la particion segun lo que venga
	for i := partition.Part_start; i < partition.Part_end+1; i++ {
		var init int8 = '0'
		if types == "full" {
			init = 0
		}
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escribirBytes(file, binarioTemp.Bytes())
	}

	//Nos posicionamos al inicio de la particion
	/*file.Seek(partition.Part_start, 0)
	//Escribe en el inicio de la particion, la ruta del path del archivo de usuarios
	var path = []byte(pathFile)
	o := &path
	var binarioTemp bytes.Buffer
	binary.Write(&binarioTemp, binary.BigEndian, o)
	escribirBytes(file, binarioTemp.Bytes())*/

	var inicio_particion int64 = partition.Part_start
	var size_particion int64 = partition.Part_size
	var cero byte = '0'

	var sizeAVD = unsafe.Sizeof(STRUCTURES.ARBOLVIRTUALDIR{})
	var sizeDD = unsafe.Sizeof(STRUCTURES.DIRECTORYDETAIL{})
	var sizeInodo = unsafe.Sizeof(STRUCTURES.TABLEINODE{})
	var sizeBloque = unsafe.Sizeof(STRUCTURES.DATABLOCK{})
	var sizeBitacora = unsafe.Sizeof(STRUCTURES.LOG{})
	var sizeSB = unsafe.Sizeof(STRUCTURES.SUPERBOOT{})

	var nEstructuras = (size_particion - (2 * int64(sizeSB))) / (27 + int64(sizeAVD) + int64(sizeDD) + (5*int64(sizeInodo) + (20 * int64(sizeBloque)) + int64(sizeBitacora)))

	var cantidadAVD = nEstructuras
	//var cantidadDD = nEstructuras
	//var cantidadInodos = 5 * nEstructuras
	//var cantidadBloques = 4 * cantidadInodos //20*nEstructuras
	//var cantidadBitacoras = nEstructuras

	var superbloque = STRUCTURES.SUPERBOOT{
		SB_magic_num: 5,
		//SB_ap_bitmap_tree_dir: inicio_particion + int64(sizeSB),
	}

	/*superbloque.SB_ap_tree_dir = superbloque.SB_ap_bitmap_tree_dir + int64(cantidadAVD)
	superbloque.SB_ap_bitmap_detail_dir = superbloque.SB_ap_tree_dir + (int64(sizeAVD) * int64(cantidadAVD))
	superbloque.SB_ap_detail_dir = superbloque.SB_ap_bitmap_detail_dir + int64(cantidadDD)
	superbloque.SB_ap_bitmap_table_inode = superbloque.SB_ap_detail_dir + (int64(sizeDD) * int64(cantidadDD))
	superbloque.SB_ap_table_inode = superbloque.SB_ap_bitmap_table_inode + int64(cantidadInodos)
	superbloque.SB_ap_bitmap_blocks = superbloque.SB_ap_table_inode + (int64(sizeInodo) * int64(cantidadInodos))
	superbloque.SB_ap_blocks = superbloque.SB_ap_bitmap_blocks + int64(cantidadBloques)
	superbloque.SB_ap_log = superbloque.SB_ap_blocks + (int64(sizeBloque) * int64(cantidadBloques))*/

	Inicio_bitmapAVD := inicio_particion + int64(sizeSB)
	//Inicio_AVD := Inicio_bitmapAVD + cantidadAVD
	//Inicio_bitmapDD := Inicio_AVD + (int64(sizeAVD) * int64(cantidadAVD))
	//Inicio_DD := Inicio_bitmapDD + int64(cantidadDD)
	//Inicio_bitmapInodo := Inicio_DD + (int64(sizeDD) * int64(cantidadDD))
	//Inicio_Inodos := Inicio_bitmapInodo + int64(cantidadInodos)
	//Inicio_bitmapBloque := Inicio_Inodos + (int64(sizeInodo) * int64(cantidadInodos))
	//Inicio_Bloque := Inicio_bitmapBloque + int64(cantidadBloques)
	//Inicio_Bitacora := Inicio_Bloque + (int64(sizeBloque) * int64(cantidadBloques))

	file.Seek(inicio_particion, 0)

	//file.Write(superbloque);
	var p = superbloque
	a := &p
	var binarioP bytes.Buffer
	binary.Write(&binarioP, binary.BigEndian, a)
	file.Write(binarioP.Bytes())

	file.Seek(Inicio_bitmapAVD, 0)
	// Bitmap de AVDs
	for i := 0; i < int(cantidadAVD); i++ {
		o := &cero
		file.Seek(int64(i), 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escribirBytes(file, binarioTemp.Bytes())
	}

}
