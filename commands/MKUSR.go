package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	CONTROLLER "../controllers"
)

func MakeAUser(usr string, pwd string, id string, grp string) {
	if SearchPartitionById(id) { //VOY A BUSCAR LA PARTICION MONTADA, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
		if CONTROLLER.IsRootLogged() { //VOY AL CONTROLADOR A VER SI HAY UN SUSIARIO LOGUEADO

			///VERIFICO SI LA PASS, EL USERNAME Y EL GRUPO TENGAN LA LONG TENGAN LA LONGITUD CORRECTA
			if len(usr) <= 10 && len(pwd) <= 10 && len(grp) <= 10 {

				var partition = GetPartitionById(id) //Obtengo la particion montada, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
				//SE ABRE EL ARCHIVO
				file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
				defer file.Close()
				if err != nil {
					fmt.Println("Hay un error, no se pudo abrir el disco duro")
				}
				//OBTENGO EL SUPER BOOT
				var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO
				var str = GetContentInINodes(file, sb.SB_ap_table_inode)

				var users []string = strings.Split(str, "\n")
				var flagGrupo bool = false
				var lstId int = 0
				for i := 0; i < len(users)-1; i++ {
					var userParts = strings.Split(strings.Trim(string(users[i]), " "), ",")
					if strings.Trim(userParts[1], " ") == "G" && strings.Trim(userParts[2], " ") == grp { //Se hace un for solo en los usuarios
						ar, er := strconv.Atoi(userParts[0]) //OBTENGO EL ID DEL GRUPO
						lstId = ar
						if er != nil {
							fmt.Println(er)
						}
						flagGrupo = true
						break
					}
				}

				if flagGrupo { //SIGNIFICA QUE EL GRUPO YA ESTA CREADO
					var contadorUsuarios = 0
					for i := 0; i < len(users)-1; i++ {
						var userParts = strings.Split(strings.Trim(string(users[i]), " "), ",")
						if strings.Trim(userParts[1], " ") == "U" && strings.Trim(userParts[2], " ") == grp { //Verifica que el los usuarios ya esten creados
							contadorUsuarios = contadorUsuarios + 1
						}
					}

					var str = strconv.Itoa(lstId+contadorUsuarios) + ",U," + grp + "," + usr + "," + pwd + "\n"

					var id, err = strconv.Atoi(CONTROLLER.GetLogedUser().User_id)
					if err != nil {

					}
					CONTROLLER.BlockController_InsertText(sb, sb.SB_ap_table_inode, str, file, int64(id))

				} else {
					fmt.Println("El grupo " + grp + " no ha sido creado todavia.")
				}

			} else {
				fmt.Println("La longitud del username, password y el grupo no debe exceder los 10 caracteress")
			}

		} else {
			fmt.Println("Este comando solo puede ser ejecutado por un usiario root")
		}
	} else {
		fmt.Println("No hay particiones con el nombre " + id)
	}
}
