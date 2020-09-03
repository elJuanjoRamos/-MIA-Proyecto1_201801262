package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	CONTROLLER "../controllers"
)

func Login(usr string, password string, idPartition string) {

	//Busco particiones montadas
	if SearchPartitionById(idPartition) { //Esta funcion esta en mount_umount go
		var partition = GetPartitionById(idPartition) // Obtego la particion montada
		file, err := os.OpenFile(partition.Mount_usrtxt, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error")
		}
		//bandera que servira para buscar si existe el usiario
		var flag = false
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var userParts = strings.Split(strings.Trim(scanner.Text(), " "), ",")

			if strings.Trim(userParts[1], " ") == "U" { //Se hace un for solo en los usuarios

				for i := 3; i < len(userParts)-1; i = i + 2 {
					if strings.Trim(userParts[i], " ") == usr && strings.Trim(userParts[i+1], " ") == password {
						fmt.Println("Usuario y password insertados correctamente")
						CONTROLLER.AddLogedUser(TRIM(userParts[0]), TRIM(userParts[1]), TRIM(userParts[2]), TRIM(userParts[i]), TRIM(userParts[i+1]))
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
	return strings.Trim(entry, " ")
}
