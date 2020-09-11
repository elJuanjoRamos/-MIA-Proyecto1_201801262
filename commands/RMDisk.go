package commands

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	FUNCTION "../functions"
)

func RemoveDisk(path string) {
	//Le quito las comillas al path
	var str = FUNCTION.RemoveComilla(FUNCTION.ReplaceAll(path))

	if _, err := os.Stat(str); os.IsNotExist(err) {

		fmt.Println("=====================================================")
		fmt.Println("	The file or directory does not exist")
		fmt.Println("=====================================================")
		fmt.Println("")

	} else {
		//Obtengo el nombre del disco
		var arCommand []string = strings.Split(str, "/")

		fmt.Println("=====================================================")
		fmt.Println("  Are you sure you want to delete " + arCommand[len(arCommand)-1] + "?")
		fmt.Println("=====================================================")
		fmt.Println("")
		fmt.Print("Press Y/N: ")

		reader := bufio.NewReader(os.Stdin)
		comando, _ := reader.ReadString('\n')
		input := ""
		if runtime.GOOS == "windows" {
			input = strings.TrimRight(comando, "\r\n")
		} else {
			input = strings.TrimRight(comando, "\n")
		}
		if strings.TrimRight(input, "\n") == "Y" || strings.TrimRight(input, "\n") == "y" {
			//Borra el disco
			err := os.Remove(str)

			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Disk '" + arCommand[len(arCommand)-1] + "' successfully deleted")

		} else {
			fmt.Println("The disk '" + arCommand[len(arCommand)-1] + "'  was not erased")

		}
	}
}
