package lib

import (
	"fmt"
	"strconv"
	"strings"

	EXECUTE "../commandExecute"
	FUNCTIONCONTROLLER "../controller"
)

// This func must be Exported, Capitalized, and comment added.
func Demo() {
	fmt.Println("HI")
}

func GetCommand(commandEntry string) {
	var arCommand []string = strings.Split(FUNCTIONCONTROLLER.RemoveSpaces(commandEntry), " ")

	var command = strings.ToLower(arCommand[0])

	switch command {
	case "exec":
		ExecCommand(arCommand[1])
	case "mkdisk":
		MKDiskCommand(arCommand)
	case "rmdisk":
		RMDiskCommand(arCommand[1])
	case "fdisk":
		FDiskCommand(arCommand)
	case "mount":
		MOUNTCommand(arCommand)
	case "unmount":
		fmt.Println("six")
	}
}

//=========================EXEC COMMAND

func ExecCommand(arCommand string) {

	var commandToExecute []string
	commandToExecute = strings.Split(arCommand, "->")

	var path string
	var ubication string
	for i := 0; i < len(commandToExecute); i++ {

		if i == 0 {
			path = commandToExecute[0]
		}
		if i == 1 {
			ubication = commandToExecute[1]
		}
	}

	fmt.Println("aux " + path)
	fmt.Println("la path es  " + ubication)
}

//-============MKDIR COMMAND
func MKDiskCommand(arCommand []string) {
	var comando string
	var size int
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
			size, err := strconv.Atoi(commandToExecute[1])
			if err == nil {
				if size <= 0 {
					fmt.Println("Error, El size no puede ser menor o igual a 0")
					error = true
				} else {
					fmt.Println(size)
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
				unit = 0
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
		EXECUTE.WriteFile(path+name, unit)
	}

	fmt.Println("Comando: " + comando)
	fmt.Println("Size: ", size)
	fmt.Println("Path: " + path)
	fmt.Println("Name: " + name)
	fmt.Println("Unit: ", unit)

}

//=====RMDISK COMMAND

func RMDiskCommand(arCommand string) {

	var commandToExecute []string = strings.Split(arCommand, "->")
	var path string = FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1])
	fmt.Println(path)
}

//=======FDISK COMMAND
func FDiskCommand(arCommand []string) {
	var comando string
	var size string
	var path string
	var name string
	var unit string
	var types string
	var fit string
	var deletes string
	var adds string

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-size":
			size = commandToExecute[1]
		case "-path":
			path = FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1])
		case "-unit":
			unit = commandToExecute[1]
		case "-name":
			name = FUNCTIONCONTROLLER.ReplaceAll(commandToExecute[1])
		case "-type":
			types = commandToExecute[1]
		case "-fit":
			fit = commandToExecute[1]
		case "-delete":
			deletes = commandToExecute[1]
		case "-add":
			adds = commandToExecute[1]
		}

	}
	fmt.Println("------------------ ")
	fmt.Println("Comando: " + comando)
	fmt.Println("Size: " + size)
	fmt.Println("Path: " + path)
	fmt.Println("Name: " + name)
	fmt.Println("Unit: " + unit)
	fmt.Println("Delete: " + deletes)
	fmt.Println("Fit: " + fit)
	fmt.Println("Add: " + adds)
	fmt.Println("Type: " + types)

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
