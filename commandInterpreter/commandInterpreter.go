package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	EXECUTE "../commandExecute"
	FUNCTIONCONTROLLER "../controller"
)

// This func must be Exported, Capitalized, and comment added.

func GetCommand(commandEntry string) {
	var arCommand []string = strings.Split(FUNCTIONCONTROLLER.RemoveSpaces(commandEntry), " ")

	var command = strings.ToLower(arCommand[0])

	switch command {
	case "exec":
		ExecCommand(arCommand[1])
		break
	case "pause":
		Pause()
		break
	case "mkdisk":
		fmt.Println("--" + commandEntry)
		MKDiskCommand(arCommand)
		break
	case "rmdisk":
		fmt.Println("--" + commandEntry)
		RMDiskCommand(arCommand[1])
		break
	case "fdisk":
		fmt.Println("--" + commandEntry)
		FDiskCommand(arCommand)
		break
	case "mount":
		fmt.Println("--" + commandEntry)
		MOUNTCommand(arCommand)
		break
	case "unmount":
		fmt.Println("--" + commandEntry)
		fmt.Println("six")
		break
	case "readdisk":
		fmt.Println("--" + commandEntry)
		ReadDiskCommand(arCommand[1])
		break

	default:

	}
}

//========================READ COMAND
func ReadDiskCommand(arCommand string) {

	var commandToExecute []string = strings.Split(arCommand, "->")
	var path string = FUNCTIONCONTROLLER.ReplaceAll(FUNCTIONCONTROLLER.RemoveComilla(FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1])))
	EXECUTE.ReadFile(path)

}

//=========================EXEC COMMAND

func ExecCommand(arCommand string) {

	var commandToExecute []string
	commandToExecute = strings.Split(arCommand, "->")

	var fpath string = commandToExecute[1]

	readFile, err := os.Open(fpath)

	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	//var fileTextLines []string
	fmt.Println("==================================")
	fmt.Println("	COMMAND LIST		   ")
	fmt.Println("   list of running commands   ")
	fmt.Println("==================================")

	//var tieneSalto = false
	//var anteriorTenia = true

	for fileScanner.Scan() {
		/*var contiene bool = strings.ContainsAny(fileScanner.Text(), "\*");
		if  contiene {
			fileTextLines = append(fileTextLines, fileScanner.Text())

		} else {

		}*/

		GetCommand(fileScanner.Text())

	}

	readFile.Close()
	/*for _, eachline := range fileTextLines {
		fmt.Println(eachline)
	}*/

}

//============PAUSE COMMAND

func Pause() {
	EXECUTE.ExecutePause()
}

//-============MKDIR COMMAND
func MKDiskCommand(arCommand []string) {
	var size int64
	var path string
	var name string
	var unit int64 = 0

	//Manejar un error
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-size":
			//trata de covertir el size a numero
			temp, err := strconv.Atoi(commandToExecute[1])
			if err == nil {
				if temp <= 0 {
					fmt.Println("Error, El size no puede ser menor o igual a 0")
					error = true
				} else {
					size = int64(temp)
				}
			} else {
				fmt.Println("Error, el size establecido no se puede convertir a numero")
				error = true
			}
		case "-path":
			path = FUNCTIONCONTROLLER.RemoveComilla(FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1]))
		case "-unit":

			if strings.ToLower(commandToExecute[1]) == "k" {
				unit = 1024
			} else if strings.ToLower(commandToExecute[1]) == "m" {
				unit = 1024 * 1024
			} else {
				fmt.Println("Error, no se reconoce el tipo " + commandToExecute[1])
				error = true
			}
		case "-name":

			var nameArray []string = strings.Split(commandToExecute[1], ".")
			if nameArray[1] == "dsk" {
				name = FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1])
			} else {
				fmt.Println("No contiene la extension correcta")
				error = true
			}

		}
	}
	//Verifica que no haya error
	if !error {
		//si unit = 0, significa que no vino en el comando, por default es 1 Mb
		if unit == 0 {
			unit = 1024 * 1024
		}
		EXECUTE.WriteFile(name, path, unit*size)
	}

}

//=====RMDISK COMMAND

func RMDiskCommand(arCommand string) {

	var commandToExecute []string = strings.Split(arCommand, "->")
	var path string = FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1])
	EXECUTE.RemoveDisk(path)
}

//=======FDISK COMMAND
func FDiskCommand(arCommand []string) {
	//var comando string
	var size int64 = 0
	var unit int64 = 0
	var path string
	var name string
	var types byte = 0
	var fit byte = 0
	//var deletes string
	//var adds string
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-size":
			//trata de covertir el size a numero
			temp, err := strconv.Atoi(commandToExecute[1])
			if err == nil {
				if temp <= 0 {
					fmt.Println("Error, El size no puede ser menor o igual a 0")
					error = true
				} else {
					size = int64(temp)
				}
			} else {
				fmt.Println("Error, el size establecido no se puede convertir a numero")
				error = true
			}
		case "-path":
			path = FUNCTIONCONTROLLER.RemoveComilla(FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1]))
		case "-unit":
			if strings.ToLower(commandToExecute[1]) == "k" {
				unit = 1024
			} else if strings.ToLower(commandToExecute[1]) == "m" {
				unit = 1024 * 1024
			} else if strings.ToLower(commandToExecute[1]) == "b" {
				unit = 1
			} else {
				fmt.Println("Error, no se reconoce el tipo " + commandToExecute[1])
				error = true
			}
		case "-name":
			name = FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1])
		case "-type":
			types = commandToExecute[1][0]
		case "-fit":
			fit = commandToExecute[1][0]
		case "-delete":
			//deletes = commandToExecute[1]
		case "-add":
			//adds = commandToExecute[1]
		}

	}

	if !error {
		if unit == 0 {
			unit = 1024
		}
		if fit == 0 {
			fit = 'W'
		}
		if types == 0 {
			types = 'P'
		}
		EXECUTE.FormatDisk(path, unit*size, name, types, fit)

	}
	/*fmt.Println("------------------ ")
	fmt.Println("Comando: " + comando)
	fmt.Println("Size: ", size)
	fmt.Println("Path: " + path)
	fmt.Println("Name: " + name)
	fmt.Println("Unit: ", unit)
	fmt.Println("Delete: " + deletes)
	fmt.Println("Add: " + adds)*/

}

//=======MOUNT COMMAND
func MOUNTCommand(arCommand []string) {
	var comando string = arCommand[0]
	var path string
	var name string
	var id string

	if len(arCommand) > 1 {
		for i := 1; i < len(arCommand); i++ {

			var commandToExecute = strings.Split(arCommand[i], "->")
			var aux string = strings.ToLower(commandToExecute[0])
			switch aux {
			case "-path":
				path = FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1])
			case "-name":
				name = FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1])
			case "#id":
				id = commandToExecute[1]
			}

		}
	}

	fmt.Println("------------------ ")
	fmt.Println("Comando: " + comando)
	fmt.Println("Path: " + path)
	fmt.Println("Name: " + name)
	fmt.Println("ID: " + id)

}
