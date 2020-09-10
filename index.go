package main

import (
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

	comando1 := "exec -path->\"/home/eljuanjoramos/Documentos/MANEJO E IMPLEMENTACION DE ARCHIVOS/-MIA-Proyecto1_201801262/archivo\""
	INTERPRETE.GetCommand(comando1)
	/*finalizar := 0
	inicio := "╔══════════════════════════════════════════════════╗\n"
	inicio += "  UNIVERSIDAD DE SAN CARLOS DE GUATEMALA\n"
	inicio += "  MANEJO E IMPLEMENTACIÓN DE ARCHIVOS A-\n"
	inicio += "  JUAN JOSE RAMOS CAMPOS\n"
	inicio += "  201801262\n"
	inicio += "╠══════════════════════════════════════════════════╣\n"
	inicio += "  Escriba X para finalizar.\n"
	inicio += "╚══════════════════════════════════════════════════╝"
	fmt.Println(inicio)
	var comando string = ""

	for finalizar != 1 {
		lecturaBuffer := bufio.NewReader(os.Stdin)
		fmt.Print("->")
		strComando, _ := lecturaBuffer.ReadString('\n')

		if strings.TrimRight(strComando, "\n") == "X" {
			finalizar = 1
		} else {
			if strComando != "" {
				if strings.Contains(strComando, "\\*") {
					remover := strings.Replace(strComando, "\\*", "", 1)
					comando += strings.TrimRight(remover, "\n")
				} else {
					if comando != "" {
						comandoTrim := strings.TrimRight(strComando, "\n")
						comando += comandoTrim
						//fmt.Print("COMANDO: " + comando)
						INTERPRETE.GetCommand(comando)
						comando = ""
					} else {
						//fmt.Print("COMANDO NUEVO: " + strComando)
						comando := strings.TrimRight(strComando, "\n")
						INTERPRETE.GetCommand(comando)
					}
				}

			}
		}
	}*/

}
