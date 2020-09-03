package main

import (
	"fmt"
	"strings"

	//EXECUTE "./commandExecute"
	//EXECCOMAND "./commandExecute"
	INTERPRETE "./interpreter"
)

func main() {
	interpret()
	//EXECUTE.WriteFile()
	//fmt.Println("Reading File: ")
	//EXECUTE.ReadFile()
}

//FUNCION INTERPRETE
func interpret() {
	//finalizar := 0
	//inicio := "IMPLEMENTACION DE ARCHIVOS.\n ('x' FINALIZAR)"
	//comando := "exec –path->/home/Desktop/calificacion.mia"

	comando1 := "exec -path->\"/home/eljuanjoramos/Documentos/MANEJO E IMPLEMENTACION DE ARCHIVOS/-MIA-Proyecto1_201801262/archivo\""
	//comando1 := "Mkdisk -size->32 -path->\"/home/eljuanjoramos/Documentos/MANEJO E IMPLEMENTACION DE ARCHIVOS/-MIA-Proyecto1_201801262/main/\" -name->Disco1.dsk -uniT->k"
	//var arCommand []string = strings.Split(comando1, " ")

	//var command = strings.ToLower(arCommand[0])

	//comando := "rmDisk -path->\"/home/eljuanjoramos/Documentos/MANEJO E IMPLEMENTACION DE ARCHIVOS/-MIA-Proyecto1_201801262/main/Disco16.dsk\""
	//comando := "mount -path->/home/Disco1.dsk -name->Part1 #id->vd0a1"
	//comando2 := "mount -path->/home/Disco2.dsk -name->Part1 #id->vdb1"
	//comando3 := "mount -path->/home/Disco3.dsk -name->Part2 #id->vdc1"
	//comando4 := "mount"
	//comando5 := "mount -path->/home/Disco1.dsk -name->Part2 #id->vda2"

	/*if command == "exec" {
		EXECCOMAND.ReadEntryFile(arCommand[1])
	} else {
		INTERPRETE.GetCommand(comando1)
	}*/

	INTERPRETE.GetCommand(comando1)

	//INTERPRETE.GetCommand(comando)
	//fmt.Println(inicio)

	/*for finalizar != 1 {
		fmt.Println("Ingresar Comandos: ")
		reader := bufio.NewReader(os.Stdin)

		//PARSEO DEPENDIENDO EL SISTEMA OPERATIVO PARA QUITAR SALTO DE LINEA Y RETORNO
		comando, _ := reader.ReadString('\n')
		input := ""
		if runtime.GOOS == "windows" {
			input = strings.TrimRight(comando, "\r\n")
		} else {
			input = strings.TrimRight(comando, "\n")
		}

		if strings.TrimRight(input, "\n") == "x" {
			finalizar = 1
		} else {
			if input != "" {
				commandLine(input)
			}
		}
	}*/
}

//RECONOCE LOS COMANDOS Y LOS MANDA A LA FUNCION CORRESPONDIENTE
func commandLine(input string) {
	var commandArray []string
	commandArray = strings.Split(input, " ")

	var command = strings.ToLower(commandArray[0])

	switch command {
	case "exec":
		runExec(commandArray[1])
	case "mkdisk":
		runMk(commandArray)
	case "rmdisk":
		runRm(commandArray[1])
	case "fdisk":
		runFd(commandArray)
	case "mount":
		fmt.Println("five")
	case "unmount":
		fmt.Println("six")
	}
}

func runExec(commandArray string) {

	var commandRead []string
	commandRead = strings.Split(commandArray, "->")
	fmt.Println(commandRead[1])
}

func runMk(commandArray []string) {
	var comando string
	var size string
	var path string
	var name string
	var unit string

	for i := 1; i < len(commandArray); i++ {
		var commandRead = strings.Split(commandArray[i], "->")
		var aux string = strings.ToLower(commandRead[0])
		switch aux {
		case "-size":
			size = commandRead[1]
		case "-path":
			path = commandRead[1]
		case "-unit":
			unit = commandRead[1]
		case "-name":
			name = commandRead[1]
		}
	}
	fmt.Println("Comando: " + comando)
	fmt.Println("Size: " + size)
	fmt.Println("Path: " + path)
	fmt.Println("Name: " + name)
	fmt.Println("Unit: " + unit)
}

func runRm(commandArray string) {

	var commandRead []string
	commandRead = strings.Split(commandArray, "->")

	fmt.Println(commandRead[1])
}

func runFd(commandArray []string) {
	var comando string
	var size string
	var path string
	var name string
	var unit string
	var types string
	var fit string
	var deletes string
	var adds string

	for i := 1; i < len(commandArray); i++ {
		var commandRead = strings.Split(commandArray[i], "->")
		var aux string = strings.ToLower(commandRead[0])
		switch aux {
		case "-size":
			size = commandRead[1]
		case "-path":
			path = commandRead[1]
		case "-unit":
			unit = commandRead[1]
		case "-name":
			name = commandRead[1]
		case "-type":
			types = commandRead[1]
		case "-fit":
			fit = commandRead[1]
		case "-delete":
			deletes = commandRead[1]
		case "-add":
			adds = commandRead[1]
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

func contieneComilla(comando string) bool {
	if strings.ContainsAny(comando, "\"") {
		return true
	}
	return false
}

func corregirEspacio(comando string) string {
	var nuevaCadena string = ""

	for i := 0; i < len(comando); i++ {
		nuevaCadena = nuevaCadena + string(comando[i])

		if string(comando[i]) == "\"" {

			for j := i + 1; j < len(comando); j++ {

				if string(comando[j]) == " " {
					nuevaCadena = nuevaCadena + "-"
				} else {
					nuevaCadena = nuevaCadena + string(comando[j])
				}

				if string(comando[j]) == "\"" {
					i = j
					break
				}
			}
		}
	}
	return nuevaCadena
}
