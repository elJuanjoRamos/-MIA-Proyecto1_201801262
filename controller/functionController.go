package controller

import (
	"fmt"
	"os"
	"strings"
)

//Re
func ReplaceAll(str string) string {

	return strings.ReplaceAll(str, "*", " ")
}

//Remueve los espacios
func RemoveSpaces(command string) string {
	if strings.ContainsAny(command, "\"") {
		var nuevaCadena string = ""

		for i := 0; i < len(command); i++ {
			nuevaCadena = nuevaCadena + string(command[i])

			if string(command[i]) == "\"" {

				for j := i + 1; j < len(command); j++ {

					if string(command[j]) == " " {
						nuevaCadena = nuevaCadena + "*"
					} else {
						nuevaCadena = nuevaCadena + string(command[j])
					}

					if string(command[j]) == "\"" {
						i = j
						break
					}
				}
			}
		}

		return nuevaCadena
	}
	return command
}

//Remueve las comillas
func RemoveComilla(command string) string {
	if strings.ContainsAny(command, "\"") {
		var str string = strings.ReplaceAll(command, "\"", " ")
		return strings.Trim(str, " ")
	}
	return command
}

//Funcion si conteiene, revisa si una cadena contiene algo
func IfContains(str string, strContains string) bool {
	return strings.ContainsAny(str, strContains)
}

//================FUNCIONES ESPECIALES

//Crea un directorio si no existe
func CreateADirectory(path string) {

	//Revisa si existe o no el directorio
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("La path no existe")
		//Create a folder/directory at a full qualified path
		err := os.Mkdir(path, 0755)
		if err != nil {
			fmt.Println("No se puede crear el directorio ", err)
		} else {
			fmt.Println("El directorio fue creado correctamente")
		}
	}
}
