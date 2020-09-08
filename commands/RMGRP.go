package commands

func RemoveGroup(name string, id string) {
	/*if SearchPartitionById(id) { //VOY A BUSCAR LA PARTICION MONTADA, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
		if CONTROLLER.IsRootLogged() { //VOY AL CONTROLADOR A VER SI HAY UN SUSIARIO LOGUEADO

			var partition = GetPartitionById(id) //Obtengo la particion montada, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
			var ifExist, idUser = VerifyGroupInFile(partition.Mount_usrtxt, name)

			//Significa que existe
			if ifExist && idUser != -1 && idUser != 0 {
				file, err := os.OpenFile(partition.Mount_usrtxt, os.O_RDWR|os.O_CREATE, os.ModeAppend)
				defer file.Close()
				if err != nil {
					log.Fatal(err)
				}

				scanner := bufio.NewScanner(file)
				var str = ""
				for scanner.Scan() {
					var userParts = strings.Split(strings.Trim(scanner.Text(), " "), ",")
					//Significa que el grupo ya existe
					if strings.Trim(userParts[2], " ") == name && strings.Trim(userParts[1], " ") == "G" {
						str = str + strconv.Itoa(0) + ",G," + name + "\n"
					} else {
						str = str + scanner.Text() + "\n"
					}
				}

				//SE BORRA EL ARCHIVO
				e := os.Remove(partition.Mount_usrtxt)
				if e != nil {
					log.Fatal(e)
				}

				//SE MANDA A CREAR OTRA VEZ
				FUNCTION.CreateAFile(partition.Mount_usrtxt, str)
				//file.WriteString(text)
			} else {
				fmt.Println("El grupo no existe")
			}

		} else {
			fmt.Println("Este comando solo puede ser ejecutado por un usiario root")
		}
	} else {
		fmt.Println("No hay particiones con el nombre " + id)
	}*/
}
