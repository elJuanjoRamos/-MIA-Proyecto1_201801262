package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	CONTROLLER "../controllers"
	STRUCTURES "../structures"
)

func MostarFiles(files []string, id string) {

	if SearchPartitionById(id) {
		//se obtiene la particion montada
		var partition = GetPartitionById(id) //esta funcion se encuenetra en commands/Moun_Umount.go
		//Se abre el archivo para modificarlo

		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		} else {

			//OBTENGO EL SUPER BOOT
			var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO
			//NOS SITUAMOS AL INICIO DEL BITMAP DE DETALLE DIRECTORIO
			file.Seek(sb.SB_ap_bitmap_detail_dir, 0)

			//ITERO EN EL BITMAP DE BLOQUES PARA SABER CUANTOS HAY
			b1 := make([]byte, (sb.SB_ap_detail_dir - sb.SB_ap_bitmap_detail_dir))
			n1, err := file.Read(b1)
			if err != nil {

			}
			var contador = 0
			for i := 0; i < len(string(b1[:n1])); i++ {
				if string(b1[:n1][i]) == "1" {
					contador = contador + 1
				}
			}

			if len(files) == 1 {

				var archivo = files[0]
				dir, name := filepath.Split(archivo)

				var texto, bandera = GetFileContent(sb, contador, name, dir, file)

				if bandera {
					fmt.Println("=================")
					fmt.Println("> El contendido de File:", name) //File: file.name
					fmt.Println("=================")
					fmt.Println(texto)
				} else {
					fmt.Println("=================")
					fmt.Println("> No existe el archivo:", name) //File: file.name
					fmt.Println("=================")

				}

			} else {
				var arregloTemp []string
				var strToMostrar = ""

				for i := 0; i < len(files); i++ {
					var archivo = files[i]
					dir, name := filepath.Split(archivo)
					strToMostrar = dir
					arregloTemp = append(arregloTemp, name)
				}
				strToMostrar = ""
				for i := 0; i < len(arregloTemp); i++ {

					var texto, bandera = GetFileContent(sb, contador, arregloTemp[i], "/", file)
					if bandera {
						strToMostrar = strToMostrar + "\n#" + texto + " from file: " + arregloTemp[i]
					} else {
						fmt.Println("=================")
						fmt.Println("> No existe el archivo:", arregloTemp[i]) //File: file.name
						fmt.Println("=================")
					}
				}

				fmt.Println("===================================")
				fmt.Println("> El contendido de los archivos es:") //File: file.name
				fmt.Println("===================================")
				fmt.Println(strToMostrar)

			}
			//OBTENGO EL SUPER BOOT
			//var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO

			//CONTROLLER.MakeAnDirectoryInDisk(sb, path, p, file, (CONTROLLER.GetLogedUser()).User_username)

		}

	} else {
		fmt.Println("No existe una particion con el id " + id)
	}
}

func GetFileContent(sb STRUCTURES.SUPERBOOT, contador int, filename string, dir string, file *os.File) (string, bool) {
	var texto = ""
	var bandera = false
	for i := 0; i < contador; i++ {

		var inicio = sb.SB_ap_detail_dir + int64(i)*int64(unsafe.Sizeof(STRUCTURES.DIRECTORYDETAIL{}))
		var directory = CONTROLLER.GETDetails(inicio, file)

		if directory.DD_file_lleno && !bandera {
			texto, bandera = ValidateBloque(directory.DD_file_nombre, filename, directory.DD_file_ap_inodo, file, texto, bandera)
		}
		if directory.DD_file_lleno2 && !bandera {
			texto, bandera = ValidateBloque(directory.DD_file_nombre2, filename, directory.DD_file_ap_inodo2, file, texto, bandera)
		}
		if directory.DD_file_lleno3 && !bandera {
			texto, bandera = ValidateBloque(directory.DD_file_nombre3, filename, directory.DD_file_ap_inodo3, file, texto, bandera)

		}
		if directory.DD_file_lleno4 && !bandera {
			texto, bandera = ValidateBloque(directory.DD_file_nombre4, filename, directory.DD_file_ap_inodo4, file, texto, bandera)
		}
		if directory.DD_file_lleno5 && !bandera {
			texto, bandera = ValidateBloque(directory.DD_file_nombre5, filename, directory.DD_file_ap_inodo5, file, texto, bandera)
		}

	}
	return texto, bandera
}

func ValidateBloque(DD_file_nombre [20]byte, filename string, DD_file_ap_inodo int64, file *os.File, texto string, bandera bool) (string, bool) {
	if GetAllName(DD_file_nombre) == filename {

		if DD_file_ap_inodo != -1 {
			var inodo = CONTROLLER.GETInode(DD_file_ap_inodo, file)

			for i := 0; i < len(inodo.I_array_bloques); i++ {

				if inodo.I_array_bloques[i] != 0 {
					var bloque, encontrado = CONTROLLER.GetBlockOcupado(inodo.I_array_bloques[i], file)
					if encontrado {
						texto = texto + CONTROLLER.GetStringInBlock(bloque.DB_data)
						bandera = true
					}
				}
			}
		}
	}
	return texto, bandera
}
