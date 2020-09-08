package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"

	CONTROLLER "../controllers"
	STRUCTURES "../structures"
)

func CreateGroup(name string, id string) {
	if SearchPartitionById(id) { //VOY A BUSCAR LA PARTICION MONTADA, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
		if CONTROLLER.IsRootLogged() { //VOY AL CONTROLADOR A VER SI HAY UN SUSIARIO LOGUEADO

			var partition = GetPartitionById(id) //Obtengo la particion montada, ESTE METODO ESTA EN MOUNT_UMOUNT.GO

			//SE ABRE EL ARCHIVO
			file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
			defer file.Close()
			if err != nil {
				fmt.Println("Hay un error, no se pudo abrir el disco duro")
			}

			//OBTENGO EL SUPER BOOT
			var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO

			//NOS SITUAMOS AL INICIO DEL BLOQUE INODOS, EL PRIMER INODO ES EL DE INICIO, DONDE SE ENCUENTRA EL ARCHIVO USERS.TXT DE LOS USUAIROS
			//Y PASSWORD GUARDADAS
			//file.Seek(sb.SB_ap_table_inode, 0)

			//NOS SITUAMOS AL INICIO DEL DE BLOQQUE PARA LEER EL TEXTO DENTRO DE ELLOS
			//file.Seek(sb.SB_ap_blocks, 0)

			var str = GetContentInINodes(file, sb.SB_ap_table_inode)

			var users []string = strings.Split(str, "\n")
			var flag bool = false
			var lstId int = 0
			for i := 0; i < len(users)-1; i++ {
				var userParts = strings.Split(strings.Trim(string(users[i]), " "), ",")
				if strings.Trim(userParts[1], " ") == "G" { //Se hace un for solo en los usuarios
					ar, er := strconv.Atoi(userParts[0])
					lstId = ar
					if er != nil {
						fmt.Println(er)
					}
					if userParts[2] == name { //Significa que el grupo existe
						flag = true
						break
					}
				}
			}
			if !flag { //Signiica que el grupo no exsite
				var iduser, err = strconv.Atoi(CONTROLLER.GetLogedUser().User_id)
				if err != nil {

				}
				var temp = strconv.Itoa(lstId+1) + ",G," + name + "\n"
				CONTROLLER.BlockController_InsertText(sb, sb.SB_ap_table_inode, temp, file, int64(iduser))
			}

		} else {
			fmt.Println("Este comando solo puede ser ejecutado por un usuario root")
		}
	} else {
		fmt.Println("No hay particiones con el nombre " + id)
	}
}

func GetContentInINodes(file *os.File, inodo int64) string {
	if inodo != -1 {
		//NOS SITUAMOS EN LA POSICION DEL INODO
		file.Seek(inodo, 0)
		//LEEMOS LA ESTRUCTURA TABLAINODO
		m := STRUCTURES.TABLEINODE{}
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
		var str = ""
		//ITERAMOS EN LOS BLOQUES DEL INIDO
		for i := 0; i < len(m.I_array_bloques); i++ {

			//VERIFICAMOS QUE LA POSICION DEL ARREGLO NO ESTE VACIA
			if m.I_array_bloques[i] != 0 {
				file.Seek(m.I_array_bloques[i], 0)
				//LEEMOS LA ESTRUCTURA BLOQUE
				m1 := STRUCTURES.DATABLOCK{}
				//Obtenemos el tamanio del BLOQUE
				var size1 int = int(unsafe.Sizeof(m1))
				//Lee la cantidad de <size> bytes del archivo
				data1 := leerBytes(file, size1)
				//Convierte la data en un buffer,necesario para
				//decodificar binario
				buffer1 := bytes.NewBuffer(data1)

				//Decodificamos y guardamos en la variable m
				err1 := binary.Read(buffer1, binary.BigEndian, &m1)
				if err1 != nil {

					log.Fatal("binary.Read failed", err1)
				}
				//CONCATENAMOS EL CONTENIDO DE LOS BLOQUES DE ESE INODO
				str = str + GetStringByBytes(m1.DB_data)
			}

		}
		return str + GetContentInINodes(file, m.I_ap_indirecto)
	}
	return ""
}

func GetStringByBytes(name [25]byte) string {
	var s string = ""
	for _, v := range name {
		if v != 0 {
			s = s + string(v)
		}
	}
	return s
}

//Se encarga de buscar la psosicion en bytes de los bloques disponibles
func ObtenerBytesBloquesSiguientes(file *os.File, inodo int64, str string, sb STRUCTURES.SUPERBOOT) {
	if inodo != -1 {
		//NOS SITUAMOS EN LA POSICION DEL INODO
		file.Seek(inodo, 0)
		//LEEMOS LA ESTRUCTURA TABLAINODO
		m := STRUCTURES.TABLEINODE{}
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
		//ITERAMOS EN LOS BLOQUES DEL INIDO
		var iterInString = 0

		for i := 0; i < len(m.I_array_bloques); i++ {
			var contador = 0
			var strTemp = ""

			for j := iterInString; j < len(str); j++ {
				strTemp = strTemp + string(str[j])
				contador = contador + 1
				if contador == 25 || ((contador + iterInString) == int(len(str))) {
					iterInString = iterInString + contador
					break
				}
			}

			if m.I_array_bloques[i] != 0 {
				file.Seek(m.I_array_bloques[i], 0)
				var block = STRUCTURES.DATABLOCK{}

				copy(block.DB_data[:], strTemp)
				//MANDO A ESCRIBIR EL BLOQUE
				block11 := &block
				var binario7 bytes.Buffer
				binary.Write(&binario7, binary.BigEndian, block11)
				escribirBytes(file, binario7.Bytes())
				strTemp = ""
			} else { //Si es cero, signfica que debo crear mas bloques

				if strTemp != "" {
					//ITERO EN EL BITMAP DE BLOQUES PARA SABER CUANTOS HAY Y EN QUE POSICION DEBO CREAR EL NUEVO BLOQUE
					file.Seek(sb.SB_ap_bitmap_blocks, 0)
					b1 := make([]byte, (sb.SB_ap_blocks - sb.SB_ap_bitmap_blocks))
					n1, err := file.Read(b1)
					if err != nil {

					}
					var contador int64 = 0
					for i := 0; i < len(string(b1[:n1])); i++ {
						if string(b1[:n1][i]) == "1" {
							contador = contador + 1
						}
					}
					//NOS SITUAMOS EN EL BITMAP DE BLOQUE correspondiente Y ESCRIBIMOS UNO
					file.Seek(sb.SB_ap_bitmap_blocks+contador, 0)
					var unit int8 = '1'
					s1 := &unit
					var binario3 bytes.Buffer
					binary.Write(&binario3, binary.BigEndian, s1)
					escribirBytes(file, binario3.Bytes())

					//NOS SITUAMOS EN EL BIT CORRESPONDIENTE PARA CREAR EL NUEVO BLOQUE
					file.Seek(sb.SB_ap_blocks+(contador*int64(unsafe.Sizeof(STRUCTURES.DATABLOCK{}))), 0)

					//CREAMOS EL BLOQUE
					var block = STRUCTURES.DATABLOCK{}

					copy(block.DB_data[:], strTemp)
					//MANDO A ESCRIBIR EL BLOQUE
					block11 := &block
					var binario7 bytes.Buffer
					binary.Write(&binario7, binary.BigEndian, block11)
					escribirBytes(file, binario7.Bytes())

					//MANDAMOS A ESCRIBIR EL BITMAP CORRESPONDIENTE AL ARRAY
					m.I_array_bloques[i] = sb.SB_ap_blocks + (contador * int64(unsafe.Sizeof(STRUCTURES.DATABLOCK{})))
					m.I_count_bloques_asignados = m.I_count_bloques_asignados + 1
				}
			}
		}
		//PREGAMOS LOS CAMBIOS HECHOS A LA TABLA INODO
		//NOS SITUAMOS AL INICIO DEL INODO
		file.Seek(inodo, 0)
		//PEGAMOS OTRA VEZ EL INODO, ESTA VEZ, CON LA CANTIDAD DE BLOQUES NUEVOS
		s1 := &m
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, s1)
		escribirBytes(file, binario3.Bytes())

		//VERIFICO SI TODAVIA EXISTE CONTENIDO DENTRO DEL Texto de entrada

	}

}
