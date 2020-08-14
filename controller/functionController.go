package controller

import "strings"

//Re
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

func RemoveComilla(command string) string {
	if strings.ContainsAny(command, "\"") {
		var str string = strings.ReplaceAll(command, "\"", " ")
		return strings.Trim(str, " ")
	}
	return command
}

func IfContains(str string, strContains string) bool {
	return strings.ContainsAny(str, strContains)
}
