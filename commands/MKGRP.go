package commands

import (
	"fmt"

	CONTROLLER "../controllers"
)

func CreateGroup(name string, id string) {
	if SearchPartitionById(id) {
		if CONTROLLER.IsRootLogged() {

		} else {
			fmt.Println("Este comando solo puede ser ejecutado por un usiario root")
		}
	} else {
		fmt.Println("No hay particiones con el nombre " + id)
	}
}
