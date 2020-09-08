package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"

	CONTROLLER "../controllers"
	STRUCTURES "../structures"
)

func Login(usr string, password string, idPartition string) {

	//Busco particiones montadas
	if SearchPartitionById(idPartition) { //Esta funcion esta en mount_umount go
		var partition = GetPartitionById(idPartition) // Obtego la particion montada este comando esta en mount_umount.go

		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error, no se pudo abrir el disco duro")
		}

		//OBTENGO EL SUPER BOOT
		var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO

		//NOS SITUAMOS AL INICIO DEL BLOQUE INODOS, EL PRIMER INODO ES EL DE INICIO, DONDE SE ENCUENTRA EL ARCHIVO USERS.TXT DE LOS USUAIROS
		//Y PASSWORD GUARDADAS
		file.Seek(sb.SB_ap_table_inode, 0)

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
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		//NOS SITUAMOS AL INICIO DEL DE BLOQQUE PARA LEER EL TEXTO DENTRO DE ELLOS
		file.Seek(sb.SB_ap_blocks, 0)
		//LEEMOS EL ARCHIVO

		var cadena string = ""
		for i := 0; i < len(m.I_array_bloques); i++ {
			if m.I_array_bloques[i] != 0 {
				file.Seek(m.I_array_bloques[i], 0)
				b1 := make([]byte, (int64(unsafe.Sizeof(STRUCTURES.DATABLOCK{}))))
				n1, err := file.Read(b1)
				if err != nil {

				}
				cadena = cadena + string(b1[:n1])
			}
		}

		var users []string = strings.Split(cadena, "\n")

		//bandera que servira para buscar si existe el usiario
		var flag = false

		for i := 0; i < len(users); i++ {
			var userParts = strings.Split(strings.Trim(string(users[i]), " "), ",")

			if strings.Trim(userParts[1], " ") == "U" { //Se hace un for solo en los usuarios

				for j := 3; j < len(userParts)-1; j = j + 2 {

					if userParts[j+1] == password {
						fmt.Println("Usuario y password insertados correctamente")
						CONTROLLER.AddLogedUser(TRIM(userParts[0]), TRIM(userParts[1]), TRIM(userParts[2]), TRIM(userParts[j]), TRIM(userParts[j+1]))
						flag = true
						break
					}
				}
				if flag {
					break
				}

			}

		}
		if !flag {
			fmt.Println("El usuario no existe o se introdujeron datos incorrectos")
		}
	} else {
		fmt.Println("La particion solicitada no existe o no esta montada")
	}
}

func TRIM(entry string) string {
	entry = strings.Trim(entry, "\n")
	return strings.Trim(entry, " ")
}
