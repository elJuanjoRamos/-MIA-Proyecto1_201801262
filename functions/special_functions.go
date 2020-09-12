package functions

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

//Reemplaza los * de las path por espacios
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

func Contains(s string, substr string) bool {
	return strings.Contains(s, substr)

}

//================FUNCIONES ESPECIALES

//Verifica si existe una ruta

func IfExistDirectoryOrPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("La path no existe")
		return false
	} else {
		return true
	}
}

//Verifica si existe o no un archivo
func IfExistFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//Crea un directorio si no existe
func CreateADirectory(path string) {
	//Revisa si existe o no el directorio
	if !IfExistDirectoryOrPath(path) {
		//Create a folder/directory at a full qualified path
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("El directorio fue creado correctamente")
		}
	}
}

func CreateAFile(path string, text string) {

	if !IfExistFile(path) {
		//Se crea el archivo user txt que contiene
		file, err := os.Create(path)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		file.WriteString(text)
	} else {
		fmt.Println("El archivo ya existe")
	}

}

//Retorna la ruta del proyecto
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func FileSize(path string) int64 {
	if IfExistFile(path) {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
			return 0
		} else {
			fi, err := os.Stat(path)
			if err != nil {
				fmt.Println(err)
			}
			return fi.Size()

		}
	}
	return 0
}

//REVISA EL CONTENIDO DEL ARCHIVO BINARIO
func ContentPrint(path string) {
	file2, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file2.Close()
	if err != nil {
		fmt.Println(err)
	}
	stats, statsErr := file2.Stat()
	if statsErr != nil {
		fmt.Println("erro")
	}

	var sizes int64 = stats.Size()
	bytess := make([]byte, sizes)

	bufr := bufio.NewReader(file2)
	_, err = bufr.Read(bytess)

	fmt.Println(bytess)
}
