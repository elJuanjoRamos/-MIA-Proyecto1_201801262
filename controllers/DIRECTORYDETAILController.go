package controllers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	STRUCTURES "../structures"
)

func MakeAVD(sb STRUCTURES.SUPERBOOT, path string, filename string, file *os.File, id string, cont string) {
	//OBTENGO EL ARBOL DE ROOT

	//SE PARTE EL PATH
	var pathlist []string = strings.Split(path, "/")
	var stringTemporal = ""
	for i := 1; i < len(pathlist); i++ {
		stringTemporal = stringTemporal + "/" + pathlist[i]
	}
	SearchPath(stringTemporal, filename, id, cont, sb.SB_ap_tree_dir, sb, file)
	SearchDetailDirectory(stringTemporal, filename, id, cont, sb.SB_ap_tree_dir, sb, file)

}

func SearchPath(path, namefile, id, cont string, inicioArbol int64, sb STRUCTURES.SUPERBOOT, file *os.File) {
	//PATH -> LA PATH QUE VOY A BUSCAR
	//DIR -> LA CARPETA QUE VOY A CREAR

	//OBTENGO EL ARBOL
	var arbol = GetTreeAVD(inicioArbol, file)

	if len(path) != 0 {

		//PARTO EL PATH QUE VIENE
		var pathlist []string = strings.Split(path, "/")
		var banderaBuscado = false

		//BUCAR EL PATH EN EL ARBOL
		//OSEA
		/*
			 path = 'root'/home/user/

			 root apunta a
				 -home
				 	- user

		*/

		var stringTemporal = ""

		if len(pathlist) != 0 {
			var bit int64 = 0
			for i := 0; i < len(arbol.Avd_ap_array_subdirectorios); i++ {
				if arbol.Avd_ap_array_subdirectorios[i] != -1 {
					var inicio = arbol.Avd_ap_array_subdirectorios[i]

					var arbolSiguiente = GetTreeAVD(inicio, file)
					if GetStringByBytes(arbolSiguiente.Avd_nombre_directorio) == CorregirName(pathlist[1]) {
						bit = arbol.Avd_ap_array_subdirectorios[i]
						banderaBuscado = true
						break
					}
				}
			}

			if banderaBuscado { // SI EXISTE LA PATH EN ESE MOMENTO, ENVIO LA
				//EN LOS SUBDIRECTORIOS
				for i := 2; i < len(pathlist); i++ {
					stringTemporal = stringTemporal + "/" + pathlist[i]
				}
				path = stringTemporal

				SearchPath(path, namefile, id, cont, bit, sb, file)

				//SI LA BANDERA ES FALSA, SIGNIFICA QUE NO EXISTE EL DIRECTORIO

			} else { //SI BANDERA SIGUE FALSO, SIGNIFICA QUE VOY A CREAR UN NUEVO ARBOL Y AL ARBOL ORIGINAL, LE VOY

				fmt.Println("==================================")
				fmt.Println("	ALERTA		   ")
				fmt.Println("   LA PATH:" + path)
				fmt.Println(" no existe y no tiene permisos para crearla")
				fmt.Println("==================================")

			}
		}
	} else { //SIGNIFICA QUE YA LLEGO AL FINAL DE LA PATH Y HAY QYE CREAR EL DIRECTORIO
		//A INSERTAR EL NUEVO ARBOL CREADO

		//LLEGA AL FINAL DEL PATH Y SE CREA UN DETALLE DE DIRECTORIO

		arbol.Avd_ap_detalle_directorio = AddDetailDirectory(namefile, sb, file, id, cont)
		//arbol.Avd_ap_detalle_directorio
	}
	file.Seek(inicioArbol, 0)
	//PEGAMOS OTRA VEZ EL ARBOL, ESTA VEZ, CON LA CANTIDAD DE SUBDIRECTORIOS NUEVOS
	s1 := &arbol
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

}

func GetTreeAVD(inodoInicio int64, file *os.File) STRUCTURES.ARBOLVIRTUALDIR {
	//NOS SITUAMOS AL INICIO DEL INODO
	file.Seek(inodoInicio, 0)

	//LEEMOS LA ESTRUCTURA TABLAINODO
	m := STRUCTURES.ARBOLVIRTUALDIR{}
	//Obtenemos el tamanio del TABLE INODO
	var size int = int(unsafe.Sizeof(m))
	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	return m
}

//AGREGA UN DETALLE DE DIRECTORIO
func AddDetailDirectory(path string, sb STRUCTURES.SUPERBOOT, file *os.File, id string, cont string) int64 {

	//NOS SITUAMOS AL INICIO DEL BITMAP DE DETALLE DIRECTORIO
	file.Seek(sb.SB_ap_bitmap_detail_dir, 0)

	//ITERO EN EL BITMAP DE BLOQUES PARA SABER CUANTOS HAY Y EN QUE POSICION DEBO CREAR EL NUEVO BLOQUE
	b1 := make([]byte, (sb.SB_ap_detail_dir - sb.SB_ap_bitmap_detail_dir))
	n1, err := file.Read(b1)
	if err != nil {

	}
	var contador int64 = 0
	for i := 0; i < len(string(b1[:n1])); i++ {
		if string(b1[:n1][i]) == "1" {
			contador = contador + 1
		}
	}

	//NOS SITUAMOS EN EL BITMAP DE DETALLE correspondiente Y ESCRIBIMOS UNO
	file.Seek(sb.SB_ap_bitmap_detail_dir+contador, 0)
	var unit int8 = '1'
	s1 := &unit
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

	//NOS SITUAMOS EN EL BIT CORRESPONDIENTE PARA CREAR EL NUEVO BLOQUE
	file.Seek(sb.SB_ap_detail_dir+(contador*int64(unsafe.Sizeof(STRUCTURES.DIRECTORYDETAIL{}))), 0)

	var time = time.Now()
	var detailDirectory = STRUCTURES.DIRECTORYDETAIL{
		DD_ap_detalle_directorio: -1,
		DD_num:                   contador + 1,
	}
	copy(detailDirectory.DD_file_nombre[:], CorregirName(path))
	copy(detailDirectory.DD_file_date_creacion[:], time.Format("2006-01-02 15:04:05"))
	copy(detailDirectory.DD_file_date_modificacion[:], time.Format("2006-01-02 15:04:05"))
	s := id
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
	}
	detailDirectory.DD_file_permiso = n
	detailDirectory.DD_file_lleno = true

	if cont == "" {
		detailDirectory.DD_file_ap_inodo = -1
	} else {

		//SETEAMOS EL APUNTADOR A INODO
		detailDirectory.DD_file_ap_inodo = INODOCONTROLLER_CreateINODO(sb, file, int64(len(cont)), n)

		BlockController_InsertText(sb, detailDirectory.DD_file_ap_inodo, cont, file, n)

	}

	file.Seek(sb.SB_ap_detail_dir+(contador*int64(unsafe.Sizeof(STRUCTURES.DIRECTORYDETAIL{}))), 0)
	dirInit1 := &detailDirectory
	var binario5 bytes.Buffer
	binary.Write(&binario5, binary.BigEndian, dirInit1)
	escribirBytes(file, binario5.Bytes())

	/*//ITERO EN EL BITMAP DE DETALLE DE DIRECTORIO PARA SABER CUANTOS HAY Y EN QUE POSICION DEBO CREAR EL NUEVO ARBOL DIRECTORIO
	file.Seek(sb.SB_ap_bitmap_detail_dir, 0)
	b1 := make([]byte, (sb.SB_ap_detail_dir - sb.SB_ap_bitmap_detail_dir))
	n1, err := file.Read(b1)
	if err != nil {

	}
	///CONTAMOS LA CANTIDAD DE DETALLE DE DIRECTORIO CREADOS
	var contador int64 = 0
	for i := 0; i < len(string(b1[:n1])); i++ {
		if string(b1[:n1][i]) == "1" {
			contador = contador + 1
		}
	}
	//NOS SITUAMOS EN EL BITMAP DE DETALLE DE DIRECTORIO correspondiente Y ESCRIBIMOS UNO
	file.Seek(sb.SB_ap_bitmap_detail_dir+contador, 0)
	var unit int8 = '1'
	s1 := &unit
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

	//NOS SITUAMOS EN EL BIT CORRESPONDIENTE PARA CREAR EL NUEVO BLOQUE
	file.Seek(sb.SB_ap_detail_dir+(contador*int64(unsafe.Sizeof(STRUCTURES.DIRECTORYDETAIL{}))), 0)

	var time = time.Now()
	var detailDirectory = STRUCTURES.DIRECTORYDETAIL{
		DD_ap_detalle_directorio: -1,
	}
	copy(detailDirectory.DD_file_nombre[:], CorregirName(path))
	copy(detailDirectory.DD_file_date_creacion[:], time.Format("2006-01-02 15:04:05"))
	copy(detailDirectory.DD_file_date_modificacion[:], time.Format("2006-01-02 15:04:05"))
	s := id
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
	}
	detailDirectory.DD_file_permiso = n
	detailDirectory.DD_file_lleno = true
	//ACA INGRESAR EL TEXTO AL INODO
	if cont == "" {
		detailDirectory.DD_file_ap_inodo = -1
	} else {
		fmt.Println("ENTRA A INODO CONTROLLER")

		//SETEAMOS EL APUNTADOR A INODO
		detailDirectory.DD_file_ap_inodo = INODOCONTROLLER_CreateINODO(sb, file, int64(len(cont)), n)

		/*var n2 int64 = detailDirectory.DD_file_ap_inodo
		s2 := strconv.FormatInt(n2, 10)
		fmt.Println(s2)*/

	/*inodo := STRUCTURES.TABLEINODE{}

	inodo = GETInode(detailDirectory.DD_file_ap_inodo, file)
	fmt.Println("ENTRA A INODO CONTROLLER")
	fmt.Println(inodo)*/

	/*	fmt.Println("ENTRA A BLOQUE CONTROLLER Y LO CREA")
		BlockController_InsertText(sb, detailDirectory.DD_file_ap_inodo, cont, file, n)

	}*/

	/*detailDirectory.Avd_ap_array_subdirectorios[0] = -1
	detailDirectory.Avd_ap_array_subdirectorios[1] = -1
	detailDirectory.Avd_ap_array_subdirectorios[2] = -1
	detailDirectory.Avd_ap_array_subdirectorios[3] = -1
	detailDirectory.Avd_ap_array_subdirectorios[4] = -1
	detailDirectory.Avd_ap_array_subdirectorios[5] = -1
	copy(detailDirectory.Avd_fecha_creacion[:], time.Format("2006-01-02 15:04:05"))
	copy(detailDirectory.Avd_nombre_directorio[:], CorregirName(path))
	copy(detailDirectory.Avd_proper[:], CorregirName(id))*/

	/*block11 := &detailDirectory
	var binario7 bytes.Buffer
	binary.Write(&binario7, binary.BigEndian, block11)
	escribirBytes(file, binario7.Bytes())*/
	return sb.SB_ap_detail_dir + (contador * int64(unsafe.Sizeof(STRUCTURES.DIRECTORYDETAIL{})))
}

func SearchDetailDirectory(path, namefile, id, cont string, inicioArbol int64, sb STRUCTURES.SUPERBOOT, file *os.File) {
	//PATH -> LA PATH QUE VOY A BUSCAR
	//DIR -> LA CARPETA QUE VOY A CREAR

	//OBTENGO EL ARBOL
	var arbol = GetTreeAVD(inicioArbol, file)

	if len(path) != 0 {

		//PARTO EL PATH QUE VIENE
		var pathlist []string = strings.Split(path, "/")
		var banderaBuscado = false

		//BUCAR EL PATH EN EL ARBOL
		//OSEA
		/*
			 path = 'root'/home/user/

			 root apunta a
				 -home
				 	- user

		*/

		var stringTemporal = ""

		if len(pathlist) != 0 {
			var bit int64 = 0
			for i := 0; i < len(arbol.Avd_ap_array_subdirectorios); i++ {
				if arbol.Avd_ap_array_subdirectorios[i] != -1 {
					var inicio = arbol.Avd_ap_array_subdirectorios[i]

					var arbolSiguiente = GetTreeAVD(inicio, file)
					if GetStringByBytes(arbolSiguiente.Avd_nombre_directorio) == CorregirName(pathlist[1]) {
						bit = arbol.Avd_ap_array_subdirectorios[i]
						banderaBuscado = true
						break
					}
				}
			}

			if banderaBuscado { // SI EXISTE LA PATH EN ESE MOMENTO, ENVIO LA
				//EN LOS SUBDIRECTORIOS
				for i := 2; i < len(pathlist); i++ {
					stringTemporal = stringTemporal + "/" + pathlist[i]
				}
				path = stringTemporal

				SearchDetailDirectory(path, namefile, id, cont, bit, sb, file)

				//SI LA BANDERA ES FALSA, SIGNIFICA QUE NO EXISTE EL DIRECTORIO

			} else { //SI BANDERA SIGUE FALSO, SIGNIFICA QUE VOY A CREAR UN NUEVO ARBOL Y AL ARBOL ORIGINAL, LE VOY

				fmt.Println("==================================")
				fmt.Println("	ALERTA		   ")
				fmt.Println("   LA PATH:" + path)
				fmt.Println(" no existe y no tiene permisos para crearla")
				fmt.Println("==================================")

			}
		}
	}
	file.Seek(inicioArbol, 0)
	//PEGAMOS OTRA VEZ EL ARBOL, ESTA VEZ, CON LA CANTIDAD DE SUBDIRECTORIOS NUEVOS
	s1 := &arbol
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

}

func GETDetails(inodoInicio int64, file *os.File) STRUCTURES.DIRECTORYDETAIL {

	//NOS SITUAMOS AL INICIO DEL INODO
	file.Seek(inodoInicio, 0)

	//LEEMOS LA ESTRUCTURA TABLAINODO
	m := STRUCTURES.DIRECTORYDETAIL{}
	//Obtenemos el tamanio del TABLE INODO
	var size int = int(unsafe.Sizeof(m))
	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	return m
}
