package commands

import (
	"fmt"

	CONTROLLER "../controllers"
)

func LogOut() {
	var isLogged, user = CONTROLLER.AddLogOut()

	if isLogged {
		fmt.Println("La sesion de " + user + " ha terminado")
	} else {
		fmt.Println("No hay ningun usuario logueado")
	}
}
