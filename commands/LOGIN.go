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
			fmt.Println("no existe")
		}
		//bandera que servira para buscar si existe el usiario
		var flag = false
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var userParts = strings.Split(strings.Trim(scanner.Text(), " "), ",")
			if len(userParts) > 3 {
				if strings.Trim(userParts[3], " ") == usr && strings.Trim(userParts[4], " ") == password {
					fmt.Println("Usuario y password insertados correctamente")
					CONTROLLER.AddLogedUser(TRIM(userParts[0]), TRIM(userParts[1]), TRIM(userParts[2]), TRIM(userParts[3]), TRIM(userParts[4]))
					flag = true
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
