package controllers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unsafe"

	STRUCTURES "../structures"
)

func MakeAnDirectoryInDisk(sb STRUCTURES.SUPERBOOT, path string, p string, file *os.File, username string) {
	//OBTENGO EL ARBOL DE ROOT

	if p != "" { //SIGNIFICA QUE BUSCO EL DIRECOTIRIO Y SI NO EXISTE, LO CREO
		SearchAndCreateDir(sb.SB_ap_tree_dir, sb, path, "/", file, username)
	} else {
		//PARTO EL PATH
		var pathlist []string = strings.Split(path, "/")
		var stringTemporal = ""
		for i := 1; i < len(pathlist)-1; i++ {
			stringTemporal = stringTemporal + "/" + pathlist[i]
		}

		SearchDir(stringTemporal, pathlist[len(pathlist)-1], username, "/", sb.SB_ap_tree_dir, sb, file)
	}

}

func SearchAndCreateDir(inicioArbol int64, sb STRUCTURES.SUPERBOOT, path string, pathPadre string, file *os.File, username string) {

	if len(path) != 0 {
		//OBTENGO EL ARBOL
		var arbol = GetArbolVirual(inicioArbol, file)
		pathPadre = GetStringByBytes(arbol.Avd_nombre_directorio)
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

					var arbolSiguiente = GetArbolVirual(inicio, file)

					if GetStringByBytes(arbolSiguiente.Avd_nombre_directorio) == CorregirName(pathlist[1]) {
						bit = arbol.Avd_ap_array_subdirectorios[i]
						banderaBuscado = true
						break
					}
				}
			}

			if banderaBuscado { // SI EXISTE LA PATH EN ESE MOMENTO,
				//EN LOS SUBDIRECTORIOS
				for i := 2; i < len(pathlist); i++ {
					stringTemporal = stringTemporal + "/" + pathlist[i]
				}
				path = stringTemporal

				SearchAndCreateDir(bit, sb, path, pathPadre, file, username)
				//SI LA BANDERA ES FALSA, SIGNIFICA QUE NO EXISTE EL

			} else { //SI BANDERA SIGUE FALSO, SIGNIFICA QUE VOY A CREAR UN NUEVO ARBOL Y AL ARBOL ORIGINAL, LE VOY
				var inicioTemp int64 = 0
				//A INSERTAR EL NUEVO ARBOL CREADO
				var bandera = false
				for i := 0; i < len(arbol.Avd_ap_array_subdirectorios); i++ {
					if arbol.Avd_ap_array_subdirectorios[i] == -1 {
						arbol.Avd_ap_array_subdirectorios[i] = CreateArbolVirtual(pathlist[1], sb, file, username)
						inicioTemp = arbol.Avd_ap_array_subdirectorios[i]
						bandera = true
						break
					}
				}

				/*///////////////////////////////////////////////////////////

								ESTO ES NUEVO


				/////////////////////////////////////////////////////////*/

				var inicio = 2
				if !bandera { // SI LA BANDERA SIGUE FALSO, SIGNIFICA QUE YA NO HAY ESPACIO EN EL ARREGLO, HAY QUE CREAR UN
					//APUNTADOR INDIRECTO
					if arbol.Avd_ap_arbol_virtual_directorio == -1 {

						arbol.Avd_ap_arbol_virtual_directorio = CreateArbolVirtual(pathPadre, sb, file, username)
						inicio = 1
					}
					inicioTemp = arbol.Avd_ap_arbol_virtual_directorio

				}

				/*///////////////////////////////////////////////////////////

								TERMINA LO NUEVO


				/////////////////////////////////////////////////////////*/

				for i := inicio; i < len(pathlist); i++ {
					stringTemporal = stringTemporal + "/" + pathlist[i]
				}
				path = stringTemporal
				SearchAndCreateDir(inicioTemp, sb, path, pathPadre, file, username)

			}

			file.Seek(inicioArbol, 0)
			//PEGAMOS OTRA VEZ EL ARBOL, ESTA VEZ, CON LA CANTIDAD DE SUBDIRECTORIOS NUEVOS
			s1 := &arbol
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, s1)
			escribirBytes(file, binario3.Bytes())

		}
	}

}

func SearchDir(path, dir, username, pathpadre string, inicioArbol int64, sb STRUCTURES.SUPERBOOT, file *os.File) {
	//PATH -> LA PATH QUE VOY A BUSCAR
	//DIR -> LA CARPETA QUE VOY A CREAR

	//home / hola
	//OBTENGO EL ARBOL
	var arbol = GetArbolVirual(inicioArbol, file)
	pathpadre = GetStringByBytes(arbol.Avd_nombre_directorio)
	if len(path) != 0 {

		//PARTO EL PATH QUE VIENE
		var pathlist []string = strings.Split(path, "/")
		var banderaBuscado = false

		var stringTemporal = ""

		if len(pathlist) != 0 {
			var bit int64 = 0
			for i := 0; i < len(arbol.Avd_ap_array_subdirectorios); i++ {
				if arbol.Avd_ap_array_subdirectorios[i] != -1 {
					var inicio = arbol.Avd_ap_array_subdirectorios[i]

					var arbolSiguiente = GetArbolVirual(inicio, file)
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

				SearchDir(path, dir, username, pathpadre, bit, sb, file)

				//SI LA BANDERA ES FALSA, SIGNIFICA QUE NO EXISTE EL DIRECTORIO

			} else {

				if arbol.Avd_ap_arbol_virtual_directorio != -1 {

					SearchDir(path, dir, username, pathpadre, arbol.Avd_ap_arbol_virtual_directorio, sb, file)

				} else {
					fmt.Println("==================================")
					fmt.Println("	ALERTA		   ")
					fmt.Println("   LA PATH:" + path)
					fmt.Println(" no existe y no tiene permisos para crearla")
					fmt.Println("==================================")
				}

				/*///////////////////////////////////////////////////////////

								ESTO ES NUEVO


				/////////////////////////////////////////////////////////*/

			}
		}
	} else { //SIGNIFICA QUE YA LLEGO AL FINAL DE LA PATH Y HAY QYE CREAR EL DIRECTORIO
		//A INSERTAR EL NUEVO ARBOL CREADO
		var bandera = false
		for i := 0; i < len(arbol.Avd_ap_array_subdirectorios); i++ {
			if arbol.Avd_ap_array_subdirectorios[i] == -1 {
				arbol.Avd_ap_array_subdirectorios[i] = CreateArbolVirtual(dir, sb, file, username)
				bandera = true
				break
			}
		}

		if !bandera {
			if arbol.Avd_ap_arbol_virtual_directorio == -1 {

				arbol.Avd_ap_arbol_virtual_directorio = CreateArbolVirtual(pathpadre, sb, file, username)
				SearchDir(path, dir, username, pathpadre, arbol.Avd_ap_arbol_virtual_directorio, sb, file)
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

/////	CRUD

func CreateArbolVirtual(path string, sb STRUCTURES.SUPERBOOT, file *os.File, username string) int64 {

	//ITERO EN EL BITMAP DE DIRECTORIOS PARA SABER CUANTOS HAY Y EN QUE POSICION DEBO CREAR EL NUEVO ARBOL DIRECTORIO
	file.Seek(sb.SB_ap_bitmap_tree_dir, 0)
	b1 := make([]byte, (sb.SB_ap_tree_dir - sb.SB_ap_bitmap_tree_dir))
	n1, err := file.Read(b1)
	if err != nil {

	}
	///CONTAMOS LA CANTIDAD DE ARBOLES CREADOS
	var contador int64 = 0
	for i := 0; i < len(string(b1[:n1])); i++ {
		if string(b1[:n1][i]) == "1" {
			contador = contador + 1
		}
	}
	//NOS SITUAMOS EN EL BITMAP DE ARBOLES correspondiente Y ESCRIBIMOS UNO
	file.Seek(sb.SB_ap_bitmap_tree_dir+contador, 0)
	var unit int8 = '1'
	s1 := &unit
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

	//NOS SITUAMOS EN EL BIT CORRESPONDIENTE PARA CREAR EL NUEVO BLOQUE
	file.Seek(sb.SB_ap_tree_dir+(contador*int64(unsafe.Sizeof(STRUCTURES.ARBOLVIRTUALDIR{}))), 0)

	var time = time.Now()
	var tree = STRUCTURES.ARBOLVIRTUALDIR{
		Avd_ap_detalle_directorio:       -1,
		Avd_ap_arbol_virtual_directorio: -1,
		Avd_num:                         contador + 1,
	}
	tree.Avd_ap_array_subdirectorios[0] = -1
	tree.Avd_ap_array_subdirectorios[1] = -1
	tree.Avd_ap_array_subdirectorios[2] = -1
	tree.Avd_ap_array_subdirectorios[3] = -1
	tree.Avd_ap_array_subdirectorios[4] = -1
	tree.Avd_ap_array_subdirectorios[5] = -1
	copy(tree.Avd_fecha_creacion[:], time.Format("2006-01-02 15:04:05"))
	copy(tree.Avd_nombre_directorio[:], CorregirName(path))
	copy(tree.Avd_proper[:], CorregirName(username))

	block11 := &tree
	var binario7 bytes.Buffer
	binary.Write(&binario7, binary.BigEndian, block11)
	escribirBytes(file, binario7.Bytes())

	return sb.SB_ap_tree_dir + (contador * int64(unsafe.Sizeof(STRUCTURES.ARBOLVIRTUALDIR{})))
}

func GetArbolVirual(inodoInicio int64, file *os.File) STRUCTURES.ARBOLVIRTUALDIR {
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

func GetStringByBytes(name [20]byte) string {
	var s string = ""
	for _, v := range name {
		if v != 0 {
			s = s + string(v)
		}
	}
	return s
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
