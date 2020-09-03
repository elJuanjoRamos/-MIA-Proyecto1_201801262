package commands

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	CONTROLLER "../controllers"
	FUNCTION "../functions"
)

func CreateGroup(name string, id string) {
	if SearchPartitionById(id) { //VOY A BUSCAR LA PARTICION MONTADA, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
		if CONTROLLER.IsRootLogged() { //VOY AL CONTROLADOR A VER SI HAY UN SUSIARIO LOGUEADO

			var partition = GetPartitionById(id) //Obtengo la particion montada, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
			var ifExist, idUser = VerifyGroupInFile(partition.Mount_usrtxt, name)
			//Significa que no existe
			if !ifExist {
				var str = "\n" + strconv.Itoa(idUser+1) + ",G," + name + "\n"
				WriteNewGroup(str, partition.Mount_usrtxt)
			} else {
				fmt.Println("El grupo ya existe")
			}

		} else {
			fmt.Println("Este comando solo puede ser ejecutado por un usuario root")
		}
	} else {
		fmt.Println("No hay particiones con el nombre " + id)
	}
}

func GetStringByBytes(name []byte) string {
	var s string = ""
	for _, v := range name {
		if v != 0 {
			s = s + string(v)
		}
	}
	return s
}

//Verifica si existe el grupo en el archivo, si no existe, retorna el ultimo id para agregar uno nuevo
func VerifyGroupInFile(path string, name string) (bool, int) {
	file, err := os.Open(path) //Abro el archivo de usuarios asociado a esa particion
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var scanner = bufio.NewScanner(file)
	var flag bool = false
	for scanner.Scan() {

		if FUNCTION.Contains(scanner.Text(), ",G,"+name) { //Verifico que exista el grupo en las lineas
			var userParts = strings.Split(strings.Trim(scanner.Text(), " "), ",")
			var id, err = strconv.Atoi(userParts[0])
			if err == nil {
				return true, id
			} else {
				return true, -1
			}
		}
	}
	//Si no existe, me posiciono en la ultima linea y mando el id de la linea para crear un nuevo grupo con el id+1
	if !flag {
		var lastLine = strings.Split(GetLastLine(file), ",")
		var d, erro = strconv.Atoi(string(lastLine[0]))
		if erro != nil {
			d = -1
			return true, d
		}
		return false, 0
	}
	return true, -1
}

func WriteNewGroup(text string, path string) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanner.Text()
	}
	file.WriteString(text)

}

func GetLastLine(file *os.File) string {
	line := ""
	var cursor int64 = 0
	stat, _ := file.Stat()
	filesize := stat.Size()
	for {
		cursor -= 1
		file.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		file.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			break
		}

		line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

		if cursor == -filesize { // stop if we are at the begining
			break
		}

	}
	return line
}
