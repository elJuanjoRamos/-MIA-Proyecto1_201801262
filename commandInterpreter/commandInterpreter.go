package lib

import (
	"fmt"
	"strings"

	FUNCTIONCONTROLLER "../controller"
)

// This func must be Exported, Capitalized, and comment added.
func Demo() {
	fmt.Println("HI")
}

func GetCommand(commandEntry string) {
	var arCommand []string

	arCommand = strings.Split(FUNCTIONCONTROLLER.RemoveSpaces(commandEntry), " ")

	var command = strings.ToLower(arCommand[0])

	switch command {
	case "exec":
		ExecCommand(arCommand[1])
	case "mkdisk":
		MKDirCommand(arCommand)
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
func MKDirCommand(arCommand []string) {
	var comando string
	var size string
	var path string
	var name string
	var unit string

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
		}
	}
	fmt.Println("Comando: " + comando)
	fmt.Println("Size: " + size)
	fmt.Println("Path: " + path)
	fmt.Println("Name: " + name)
	fmt.Println("Unit: " + unit)

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
