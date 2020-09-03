package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	CONTROLLER "../controllers"
	FUNCTION "../functions"
)

func MakeAUser(usr string, pwd string, id string, grp string) {
	if SearchPartitionById(id) { //VOY A BUSCAR LA PARTICION MONTADA, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
		if CONTROLLER.IsRootLogged() { //VOY AL CONTROLADOR A VER SI HAY UN SUSIARIO LOGUEADO

			///VERIFICO SI LA PASS, EL USERNAME Y EL GRUPO TENGAN LA LONG TENGAN LA LONGITUD CORRECTA
			if len(usr) <= 10 && len(pwd) <= 10 && len(grp) <= 10 {

				var partition = GetPartitionById(id)                                 //Obtengo la particion montada, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
				var ifExist, idUser = VerifyGroupInFile(partition.Mount_usrtxt, grp) //Verifico si existe el grupo, ESTE METODO ESTA EN MKGRP.go
				//Significa que no existe
				if ifExist {
					var str = strconv.Itoa(idUser) + ",U," + grp + "," + usr + "," + pwd + "\n"
					var str2 = "," + usr + "," + pwd + "\n"

					WriteNewUser(str, str2, partition.Mount_usrtxt, grp) //MANDO A ESCRIBIR AL USUARIO NUEVO,
				} else {
					fmt.Println("El grupo no existe")
				}

			} else {
				fmt.Println("La longitud del username, password y el grupo no debe exceder los 10 caracteress")
			}

		} else {
			fmt.Println("Este comando solo puede ser ejecutado por un usiario root")
		}
	} else {
		fmt.Println("No hay particiones con el nombre " + id)
	}
}

//Verifica si existe el los usuarios en el archivo,

func WriteNewUser(textIfNotExistGroup string, textIfExistGroup string, path string, grp string) {

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	var text = ""
	var bandera = false
	for scanner.Scan() {
		if FUNCTION.Contains(scanner.Text(), ",U,"+grp) { //Verifico que exista el grupo en las lineas
			fmt.Println("entro")
			bandera = true
			text = text + scanner.Text() + textIfExistGroup + "\n"
			break
		} else {
			text = text + scanner.Text() + "\n"
		}
	}

	if !bandera {
		file.WriteString(textIfNotExistGroup)
	} else {
		//SE BORRA EL ARCHIVO
		e := os.Remove(path)
		if e != nil {
			log.Fatal(e)
		}

		FUNCTION.CreateAFile(path, text)
	}
}
