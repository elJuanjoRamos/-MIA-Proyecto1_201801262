package interpreter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	EXECUTE "../commands"
	FUNCTION "../functions"
)

// This func must be Exported, Capitalized, and comment added.

func GetCommand(commandEntry string) {

	var arCommand []string = strings.Split(FUNCTION.RemoveSpaces(commandEntry), " ")

	var command = strings.ToLower(arCommand[0])

	switch strings.TrimRight(command, "\n") {
	case "exec":
		ExecCommand(arCommand)
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
		RMDiskCommand(arCommand)
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
		UNMOUNTCommand(arCommand)
		break
	case "readdisk":
		fmt.Println("--" + commandEntry)
		ReadDiskCommand(arCommand[1])
		break
	case "mkfs":
		fmt.Println("--" + commandEntry)
		MKfsCommand(arCommand)
		break
	case "login":
		fmt.Println("--" + commandEntry)
		LoginCommand(arCommand)
		break
	case "logout":
		fmt.Println("--" + commandEntry)
		LogOutCommand()
		break
	case "mkgrp":
		fmt.Println("--" + commandEntry)
		MKGroupCommand(arCommand)
		break
	case "rmgrp":
		fmt.Println("--" + commandEntry)
		RMGroupCommand(arCommand)
		break
	case "mkusr":
		fmt.Println("--" + commandEntry)
		MKUsrCommand(arCommand)
		break
	case "mkdir":
		fmt.Println("--" + commandEntry)
		MKDirCommand(arCommand)
		break
	case "mkfile":
		fmt.Println("--" + commandEntry)
		MKFileCommand(arCommand)
		break
	case "chmod":
		fmt.Println("--" + commandEntry)
		ChmodCommand(arCommand)
		break
	case "rep":
		fmt.Println("--" + commandEntry)
		RepCommand(arCommand)
		break
	case "cat":
		fmt.Println("--" + commandEntry)
		CatCommand(arCommand)
		break

	default:
		if !strings.Contains(command, "#") {
			fmt.Println("╔══════════════════════════════════════════════════╗")
			fmt.Println("  COMANDO NO SOPORTADO")
			fmt.Println("╚══════════════════════════════════════════════════╝")
		}
		break
	}
}

//========================READ COMAND
func ReadDiskCommand(arCommand string) {

	var commandToExecute []string = strings.Split(arCommand, "->")
	var path string = FUNCTION.ReplaceAll(FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1])))
	EXECUTE.ReadFile(path)
	//EXECUTE.SendToFull(path)
}

//=========================EXEC COMMAND

func ExecCommand(arCommand []string) {
	for i := 1; i < len(arCommand); i++ {

		if arCommand[i] != "" {
			var commandToExecute []string
			commandToExecute = strings.Split(arCommand[i], "->")

			var fpath string = commandToExecute[1]

			readFile, err := os.Open(FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(fpath)))

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
			var comando string = ""
			for fileScanner.Scan() {

				commandTem := fileScanner.Text()
				if commandTem != "" {
					if strings.Contains(commandTem, "\\*") {
						remCaracter := strings.Replace(commandTem, "\\*", "", 1)

						comando += strings.TrimRight(remCaracter, "\n")
					} else {
						if comando != "" {
							comandoTrim := strings.TrimRight(commandTem, "\n")
							comando += comandoTrim
							GetCommand(comando)
							comando = ""
						} else {
							comando := strings.TrimRight(commandTem, "\n")
							GetCommand(comando)
						}
					}
				}
			}

			readFile.Close()
			break
		}
	}

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
			path = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
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
				name = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			} else {
				fmt.Println("No contiene la extension correcta")
				error = true
			}

		}
	}

	EXECUTE.AddAccion(1, "ARCHIVO", name, "creacion del archivo en la path: "+path)

	//Verifica que no haya error
	if !error {
		//si unit = 0, significa que no vino en el comando, por default es 1 Mb
		if unit == 0 {
			unit = 1024 * 1024
		}
		EXECUTE.CreateFile(name, path, unit*size)
		EXECUTE.WriteFile(path+name, name)

		//EXECUTE.CreateReport()
	} else {
		fmt.Println("Ha ocurrido un error al intentar ejecutar el comando MKDisk")
	}
}

//=====RMDISK COMMAND

func RMDiskCommand(arCommand []string) {

	for i := 1; i < len(arCommand); i++ {
		if arCommand[i] != "" {
			var commandToExecute []string = strings.Split(arCommand[i], "->")
			var path string = FUNCTION.ReplaceAll(FUNCTION.ReplaceAll(FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))))
			file := filepath.Base(path)
			EXECUTE.AddAccion(2, "ARCHIVO", file, "eliminacion del archivo en la path: "+path)
			EXECUTE.RemoveDisk(path)
			break
		}
	}
}

//=======FDISK COMMAND
func FDiskCommand(arCommand []string) {
	//var comando string
	var size int64 = 0
	var unit int64 = 0
	var path string
	var name string
	var types string = "P"
	var fit string = "WF"
	var deletes string = ""
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
			path = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
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
			name = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
		case "-type":
			if strings.ToLower(commandToExecute[1]) == "p" || strings.ToLower(commandToExecute[1]) == "l" || strings.ToLower(commandToExecute[1]) == "e" {
				types = commandToExecute[1]
			} else {
				error = true
				fmt.Println("No se reconoce el tipo de particion a crear")
			}
		case "-fit":
			if strings.ToLower(commandToExecute[1]) == "bf" || strings.ToLower(commandToExecute[1]) == "ff" || strings.ToLower(commandToExecute[1]) == "wf" {
				fit = commandToExecute[1]
			} else {
				error = true
				fmt.Println("No se reconoce el ajuste de particion a crear")
			}
		case "-delete":
			deletes = commandToExecute[1]
		case "-add":
			//adds = commandToExecute[1]
		}

	}

	if !error {
		if unit == 0 {
			unit = 1024
		}

		if deletes != "" {
			if strings.ToLower(deletes) == "fast" || strings.ToLower(deletes) == "full" {
				EXECUTE.DeletePartition(path, name, deletes)
			} else {
				fmt.Println("El parametro establecido para delete no es valido")
				fmt.Println("Error al ejecutar el comando FDISK")
			}
		} else {
			EXECUTE.AddAccion(3, "ARCHIVO", name, "formateo de del archivo en la path: "+path)
			EXECUTE.FormatDisk(path, unit*size, name, types, fit)
		}

	} else {
		fmt.Println("Error al ejecutar el comando FDISK")
	}
}

//=======MOUNT COMMAND
func MOUNTCommand(arCommand []string) {

	var path string = "null"
	var name string = "null"

	if len(arCommand) > 1 {
		for i := 1; i < len(arCommand); i++ {

			var commandToExecute = strings.Split(arCommand[i], "->")
			var aux string = strings.ToLower(commandToExecute[0])
			switch aux {
			case "-path":
				path = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			case "-name":
				name = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			}

		}
	}

	if path != "null" && name != "null" {
		if path != "" && name != "" {
			EXECUTE.AddAccion(4, "PARTICION", name, "Montar la particion con id :'"+name+"'  dada la path: "+path)
			EXECUTE.Mount(path, name)
		} else {
			fmt.Println("Ha ocurrido un error al intentar ejecutar el comando Mount")
		}
	} else {
		EXECUTE.MountPrint()
	}
}

//======== UMOUNT COMMAND
func UNMOUNTCommand(arCommand []string) {
	var ids []string

	if len(arCommand) > 1 {
		for i := 1; i < len(arCommand); i++ {

			var commandToExecute = strings.Split(arCommand[i], "->")
			var aux string = strings.ToLower(commandToExecute[0])
			if strings.Contains(aux, "-id") {
				ids = append(ids, FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1])))
				EXECUTE.AddAccion(4, "PARTICION", commandToExecute[1], "Desmontar Particion "+commandToExecute[1])
			}
		}
	}

	if len(ids) > 0 {
		EXECUTE.UMount(ids)
	} else {
		fmt.Println("Error al intentar ejecutar el comando unmount")
	}
}

/*					PARTE 2 		*/

//==== MKFS COMMAND

func MKfsCommand(arCommand []string) {

	var id string
	var types string = "full"
	var add int64 = 0
	var unit int64 = 0

	//Manejar un error
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-id":
			id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
		case "-type":
			types = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
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
		case "-add":

			temp, err := strconv.Atoi(commandToExecute[1])
			if err == nil {
				add = int64(temp)
			} else {
				fmt.Println("Error, el size establecido no se puede convertir a numero")
				error = true
			}

		}
	}
	//Verifica que no haya error
	if !error {
		//Cuando las unidades son diferentes a cero, significa que se quiere agregar o quitar espacio de la particion
		//especificcada con el id
		if unit != 0 && add != 0 {

			//si las unidades son cero, significa que solo se quiere formatear la unidad
		} else {
			EXECUTE.AddAccion(5, "PARTICION", id, "Formateo de particion del tipo "+types)

			EXECUTE.MKFSFormatPartition(id, types)
			EXECUTE.MakeADirFirsTime("/", id)

		}

	}
}

//===== LOGIN
func LoginCommand(arCommand []string) {
	var usr string = ""
	var pwd string = ""
	var id string = ""

	//Manejar un error
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-id":
			id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "&id":
			id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-pwd":
			pwd = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "&pwd":
			pwd = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-usr":
			usr = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "&usr":
			usr = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		}
	}
	if usr == "" || pwd == "" || id == "" {
		error = true
	}
	//Verifica que no haya error
	if !error {
		EXECUTE.AddAccion(6, "LOGIN", id+"/"+pwd, "Login de usuario en la paticion "+id)
		EXECUTE.Login(usr, pwd, id)
	} else {
		fmt.Println("Error al ejecutar el comando login")
	}
}

//===== LOGOUT

func LogOutCommand() {
	EXECUTE.LogOut()
}

//====== MKGRP

func MKGroupCommand(arCommand []string) {
	var name string = ""
	var id string = ""

	//Manejar un error
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-id":
			id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-name":
			name = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		}
	}
	if name == "" || id == "" {
		error = true
	}
	//Verifica que no haya error
	if !error {
		EXECUTE.AddAccion(7, "GRPUP", name, "Crear grupo en la paticion "+id)
		EXECUTE.CreateGroup(name, id)
	} else {
		fmt.Println("Error al ejecutar el comando MKGRP")
	}
}

//===== RMGRP

func RMGroupCommand(arCommand []string) {
	var name string = ""
	var id string = ""

	//Manejar un error
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-id":
			id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-name":
			name = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		}
	}
	if name == "" || id == "" {
		error = true
	}
	//Verifica que no haya error
	if !error {
		EXECUTE.RemoveGroup(name, id)
	} else {
		fmt.Println("Error al ejecutar el comando MKGRP")
	}
}

///==== MKUSR
func MKUsrCommand(arCommand []string) {
	var usr string = ""
	var pwd string = ""
	var id string = ""
	var grp string = ""
	//Manejar un error
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-id":
			id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-pwd":
			pwd = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-usr":
			usr = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-grp":
			grp = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		}
	}
	if grp == "" || usr == "" || pwd == "" || id == "" {
		error = true
	}
	//Verifica que no haya error
	if !error {
		EXECUTE.AddAccion(8, "USER", usr+"/"+pwd, "Crear usuario en la paticion "+id)
		EXECUTE.MakeAUser(usr, pwd, id, grp)
	} else {
		fmt.Println("Error al ejecutar el comando MKUSR")
	}
}

//===== CHMOD
func ChmodCommand(arCommand []string) {
	/*var ugo string = ""
	var path string = ""
	var id string = ""
	var r string = ""
	//Manejar un error
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-id":
			id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-path":
			path = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-ugo":
			ugo = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-R":
			r = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		}
	}
	if ugo == "" || id == "" || path == "" {
		error = true
	}
	//Verifica que no haya error
	if !error {

		//EXECUTE.MakeAUser(ugo, path, id)
	} else {
		fmt.Println("Error al ejecutar el comando MKUSR")
	}*/
}

//=== MKDIR

func MKDirCommand(arCommand []string) {
	var p string = ""
	var path string = ""
	var id string = ""
	//Manejar un error
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-id":
			id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-path":
			path = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-p":
			p = "P"
			break
		}
	}
	if id == "" || path == "" {
		error = true
	}
	//Verifica que no haya error
	if !error {

		EXECUTE.AddAccion(9, "DIRECTORIO", id, "Crear directorio "+path)
		EXECUTE.MakeADir(path, id, p)
	} else {
		fmt.Println("Error al ejecutar el comando MKUSR")
	}
}

// ===== MKFILE
func MKFileCommand(arCommand []string) {
	var p string = ""
	var path string = ""
	var id string = ""
	var size int64 = 0
	var cont string = ""
	//Manejar un error
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-id":
			id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-path":
			path = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-cont":
			cont = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-p":
			p = "P"
			break
		case "-size":
			//trata de covertir el size a numero
			temp, err := strconv.Atoi(commandToExecute[1])
			if err == nil {
				if temp < 0 {
					fmt.Println("Error, El size no puede ser menor a 0")
					error = true
				} else {
					size = int64(temp)
				}
			} else {
				fmt.Println("Error, el size establecido no se puede convertir a numero")
				error = true
			}
			break
		}
	}
	if id == "" || path == "" {
		error = true
	}
	//Verifica que no haya error
	if !error {

		var abc = " abcdefghijklmnopqrstuvwzyz"
		var strtemp = ""

		if cont == "" {
			cont = abc
		}
		if size == 0 && cont != "" {
			size = int64(len(cont))
		}
		if size > int64(len(cont)) {
			var iTmp = 0
			for i := int64(len(cont)); i < size; i++ {

				if iTmp == 26 {
					iTmp = 0
				}
				cont = cont + string(abc[iTmp])
				iTmp = iTmp + 1
			}

		} else if size < int64(len(cont)) {
			for i := 0; i < int(size); i++ {
				strtemp = strtemp + string(cont[i])
			}
			cont = strtemp
		}
		EXECUTE.AddAccion(10, "ARCHIVO", id, "Crear archivo "+path+" con size "+strconv.Itoa(int(size))+" y contenido:"+cont)
		EXECUTE.MakeAFileInLogicalDisk(path, id, p, size, cont)
	} else {
		fmt.Println("Error al ejecutar el comando MKUSR")
	}
}

//==== REPORTES
func RepCommand(arCommand []string) {
	var nombre string = ""
	var path string = ""
	var id string = ""
	var ruta string = ""
	//Manejar un error
	var error bool = false

	for i := 1; i < len(arCommand); i++ {
		var commandToExecute = strings.Split(arCommand[i], "->")
		var aux string = strings.ToLower(commandToExecute[0])
		switch aux {
		case "-id":
			id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-path":
			path = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-ruta":
			ruta = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		case "-nombre":
			nombre = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			break
		}
	}

	if id == "" || path == "" || nombre == "" {
		error = true
	}
	//Verifica que no haya error
	if !error {
		EXECUTE.AddAccion(11, "ARCHIVO", id, "Generar reporte: "+nombre)

		EXECUTE.MakeAReport(path, id, strings.ToLower(nombre), ruta)
	} else {
		fmt.Println("Error al ejecutar el comando MKUSR")
	}
}

//=== CAT

func CatCommand(arCommand []string) {
	var files []string
	var id = ""
	var temporal = ""
	if len(arCommand) > 1 {
		for i := 1; i < len(arCommand); i++ {

			var commandToExecute = strings.Split(arCommand[i], "->")
			var aux string = strings.ToLower(commandToExecute[0])
			if strings.Contains(aux, "-file") {
				files = append(files, FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1])))
				temporal = temporal + "-" + FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			} else if strings.Contains(aux, "-id") {
				id = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(commandToExecute[1]))
			}

		}
	}

	if len(files) > 0 && id != "" {
		EXECUTE.AddAccion(12, "ARCHIVO", id, "Mostrar contenido en el o los archivos: "+temporal)
		EXECUTE.MostarFiles(files, id)
	} else {
		fmt.Println("Error al intentar ejecutar el comando Cat")
	}
}
