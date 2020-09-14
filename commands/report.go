package commands

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	CONTROLLER "../controllers"
	FUNCTION "../functions"
	STRUCTURES "../structures"
)

func MakeAReport(path, id, nombre, ruta string) {

	switch nombre {
	case "sb":
		MakeASBReport(path, id, ruta)
		break
	case "bm_arbdir":
		MakeABitMapReport(path, id, ruta, 1)
		break
	case "bm_detdir":
		MakeABitMapReport(path, id, ruta, 2)
		break
	case "bm_inode":
		MakeABitMapReport(path, id, ruta, 3)
		break
	case "bm_block":
		MakeABitMapReport(path, id, ruta, 4)
		break
	case "directorio":
		MakeAGenaralDirecoryReport(path, id, ruta)
		break
	case "tree_file":
		TreeFileReport(path, id, ruta)
		break
	case "bitacora":
		MakeALogReport(path, id, ruta)
		break
	}
}

//====================== REPORTE SUPER BOOT
func MakeASBReport(path, id, ruta string) { //REPORTE DEL SUPER BOOT
	if SearchPartitionById(id) { //VOY A BUSCAR LA PARTICION MONTADA, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
		var partition = GetPartitionById(id) //Obtengo la particion montada, ESTE METODO ESTA EN MOUNT_UMOUNT.GO

		//SE ABRE EL ARCHIVO
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error, no se pudo abrir el disco duro")
		}
		//OBTENGO EL SUPER BOOT
		var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO

		var s string
		for _, v := range sb.SB_hd_name {
			if v != 0 {
				s = s + string(v)
			}
		}
		var s1 string
		for _, v := range sb.SB_date {
			if v != 0 {
				s1 = s1 + string(v)
			}
		}
		var s2 string
		for _, v := range sb.SB_date_lstmount {
			if v != 0 {
				s2 = s2 + string(v)
			}
		}

		var body string = "digraph test { graph [ratio=fill];" +
			"node [label=\"Grafica\", fontsize=15, shape=plaintext];" +
			"graph [bb=\"0,0,352,154\"];" +
			"arset [label=<" +
			"<TABLE>" +
			"<TR>" + "<TD>Reporte: </TD>" + "<TD> Super Boot </TD>" + "</TR>" +
			"<TR>" + "<TD>SB_hd_name</TD>" + "<TD>" + s + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_date</TD>" + "<TD>" + s1 + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_date_lstmount</TD>" + "<TD>" + s2 + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_AVD_count</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_AVD_count)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_AVD_details_count</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_AVD_details_count)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_Inodes_count</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_Inodes_count)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_blocks_count</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_blocks_count)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_AVD_free</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_AVD_free)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_Inodes_free</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_Inodes_free)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_blocks_free</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_blocks_free)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_mount_count</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_mount_count)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_ap_bitmap_tree_dir</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_ap_bitmap_tree_dir)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_ap_tree_dir</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_ap_tree_dir)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_ap_bitmap_detail_dir</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_ap_bitmap_detail_dir)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_ap_detail_dir</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_ap_detail_dir)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_ap_bitmap_table_inode</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_ap_bitmap_table_inode)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_ap_table_inode</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_ap_table_inode)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_ap_bitmap_blocks</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_ap_bitmap_blocks)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_ap_blocks</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_ap_blocks)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_ap_log</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_ap_log)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_size_struct_tree_dir</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_size_struct_tree_dir)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_size_struct_detail_dir</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_size_struct_detail_dir)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_size_struct_inodo</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_size_struct_inodo)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_size_struct_block</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_size_struct_block)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_first_free_bit_tree_dir</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_first_free_bit_tree_dir)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_first_free_bit_detail_dir</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_first_free_bit_detail_dir)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_first_free_bit_table_dir</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_first_free_bit_table_dir)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_first_free_bit_block</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_first_free_bit_block)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_magic_num</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_magic_num)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>SB_free_space</TD>" + "<TD>" + strconv.Itoa(int(sb.SB_free_space)) + "</TD>" + "</TR>"
		body = body + "</TABLE>" + ">, ];}"

		dir, name := filepath.Split(path)

		GeneratePNG(name, body, dir)
	} else {
		fmt.Println("La particion con ID: " + id + " no esta montada")
	}
}

///===================== REPORTE BITMAP
func MakeABitMapReport(path, id, ruta string, types int) { ////REPORTES DE BITMAPS
	if SearchPartitionById(id) { //VOY A BUSCAR LA PARTICION MONTADA, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
		var partition = GetPartitionById(id) //Obtengo la particion montada, ESTE METODO ESTA EN MOUNT_UMOUNT.GO

		//SE ABRE EL ARCHIVO
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error, no se pudo abrir el disco duro")
		}
		//OBTENGO EL SUPER BOOT
		var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO

		var inicio int64 = 0
		var fin int64 = 0

		switch types {
		case 1:
			inicio = sb.SB_ap_bitmap_tree_dir
			fin = sb.SB_ap_tree_dir
			break
		case 2:
			inicio = sb.SB_ap_bitmap_detail_dir
			fin = sb.SB_ap_detail_dir
			break
		case 3:
			inicio = sb.SB_ap_bitmap_table_inode
			fin = sb.SB_ap_table_inode
			break
		case 4:
			inicio = sb.SB_ap_bitmap_blocks
			fin = sb.SB_ap_blocks
			break

		}

		file.Seek(inicio, 0)
		b1 := make([]byte, (fin - inicio))
		n1, err := file.Read(b1)
		if err != nil {

		}
		var contador int64 = 0
		var str = ""
		for i := 0; i < len(string(b1[:n1])); i++ {
			if contador == 19 {
				str = str + string(b1[:n1][i]) + "\n"
				contador = 0
			} else {
				str = str + string(b1[:n1][i]) + "|"
				contador = contador + 1
			}
		}

		dir, name := filepath.Split(path)

		FUNCTION.CreateADirectory(dir)
		FUNCTION.CreateAFile(dir+name, str)
	} else {
		fmt.Println("La particion con ID: " + id + " no esta montada")
	}
}

//======================= REPORTE LOG
func MakeALogReport(path, id, ruta string) {
	if SearchPartitionById(id) { //VOY A BUSCAR LA PARTICION MONTADA, ESTE METODO ESTA EN MOUNT_UMOUNT.GO

		var str = ""
		var bitacora = GetBitacora()

		for i := 0; i < len(bitacora); i++ {
			var temp = bitacora[i]

			if strings.Contains(temp.nombre, id) || strings.Contains(temp.contenido, id) {
				str = str +
					"Tipo Operacion:" + strconv.Itoa(temp.operacion) + "\nTipo:" + temp.tipo + "\nNombre:" + temp.nombre + "\nContenido:" + temp.contenido + "\nFecha:" + temp.fecha + "\n-----\n"
			}
		}

		dir, name := filepath.Split(path)

		FUNCTION.CreateADirectory(dir)
		FUNCTION.CreateAFile(dir+name, str)
	} else {
		fmt.Println("La particion con ID: " + id + " no esta montada")
	}
}

//===================REPORTE GENRAL
func MakeAGenaralDirecoryReport(path, id, ruta string) {
	if SearchPartitionById(id) { //VOY A BUSCAR LA PARTICION MONTADA, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
		var partition = GetPartitionById(id) //Obtengo la particion montada, ESTE METODO ESTA EN MOUNT_UMOUNT.GO

		//SE ABRE EL ARCHIVO
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error, no se pudo abrir el disco duro")
		}
		//OBTENGO EL SUPER BOOT
		var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO

		//OBTENGO EL ARBOL DE DIRECTORIOS
		var arbolRoot = CONTROLLER.GetArbolVirual(sb.SB_ap_tree_dir, file)
		var body = "digraph H { rankdir=\"LR" + "\"  "

		//ROOT
		var root = "parent [ shape=plaintext label=<\n<table border='1' cellborder='1'> \n<tr>\n\t<td bgcolor=\"chartreuse\" colspan=\"2\">" + GetAllName(arbolRoot.Avd_nombre_directorio) + "</td>\n</tr>\n" +
			"<tr>\n\t<td colspan=\"2\"> Proper: " + GetAllName(arbolRoot.Avd_proper) + "</td>\n</tr>\n" + "<tr>\n\t<td colspan=\"2\"> Proper: " + GetFecha(arbolRoot.Avd_fecha_creacion) + "</td>\n</tr>\n"

		var insideRoot = ""
		var enlacesRoot = ""
		var childRoot = ""

		for i := 0; i < len(arbolRoot.Avd_ap_array_subdirectorios); i++ {
			var data = arbolRoot.Avd_ap_array_subdirectorios[i]
			if data == -1 {
				insideRoot = insideRoot + "\n<tr>\n\t<td>AVD apun</td>\n\t<td port='port'>-1</td>\n</tr>\n"
			} else {
				insideRoot = insideRoot + "\n<tr>\n\t<td>AVD apun</td>\n\t<td port='port" + strconv.Itoa(int(data)) + "'>" + strconv.Itoa(int(data)) + "</td>\n</tr>"
				enlacesRoot = enlacesRoot + "\nparent:port" + strconv.Itoa(int(data)) + "   -> child" + strconv.Itoa(int(data)) + ";"
				childRoot = childRoot + GetAlDirectory(data, "child"+strconv.Itoa(int(data)), file)
			}
		}

		/// A DIRECOTRY
		var detalle = ""
		fmt.Println("---------")
		var detallebit = arbolRoot.Avd_ap_detalle_directorio
		if detallebit != -1 {
			insideRoot = insideRoot + "\n<tr><td>DD apun</td><td port='port" + strconv.Itoa(int(detallebit)) + "'>" + strconv.Itoa(int(detallebit)) + "</td></tr>"
			detalle = detalle + GetAlDetails(detallebit, "child"+strconv.Itoa(int(detallebit)), file)
			enlacesRoot = enlacesRoot + "\nparent:port" + strconv.Itoa(int(detallebit)) + "   -> child" + strconv.Itoa(int(detallebit)) + ";"
		}

		// APUNTADOR INDIRECTO
		if arbolRoot.Avd_ap_arbol_virtual_directorio != -1 {
			var dt = arbolRoot.Avd_ap_arbol_virtual_directorio
			insideRoot = insideRoot + "\n<tr><td>AVD INDIR</td><td port='port" + strconv.Itoa(int(dt)) + "'>" + strconv.Itoa(int(dt)) + "</td></tr>"
			enlacesRoot = enlacesRoot + "\nparent:port" + strconv.Itoa(int(dt)) + "   -> child" + strconv.Itoa(int(dt)) + ";"
			childRoot = childRoot + GetAlDirectory(dt, "child"+strconv.Itoa(int(dt)), file)
		}

		root = root + insideRoot
		root = root + "\n</table>>];"
		body = body + root + childRoot + enlacesRoot + detalle

		body = body + "\n}"

		dir, name := filepath.Split(path)

		GeneratePNG(name, body, dir)
	} else {
		fmt.Println("La particion con ID: " + id + " no esta montada")
	}
}

func GetAlDirectory(bitInicio int64, name string, file *os.File) string {

	var arbolRoot = CONTROLLER.GetArbolVirual(bitInicio, file)

	var root = "\n" + name + " [ shape=plaintext label=<\n<table border='1' cellborder='1'> \n<tr>\n\t<td bgcolor=\"chartreuse\"  colspan=\"2\">" + GetAllName(arbolRoot.Avd_nombre_directorio) + "</td>\n</tr>\n" +
		"<tr>\n\t<td colspan=\"2\"> Proper: " + GetAllName(arbolRoot.Avd_proper) + "</td>\n</tr>" + "\n<tr>\n\t<td colspan=\"2\"> Proper: " + GetFecha(arbolRoot.Avd_fecha_creacion) + "</td>\n</tr>\n"

	var insideRoot = ""
	var enlacesRoot = ""
	var childRoot = ""

	//SUBDIRECTORIOS DEL ROOT
	for i := 0; i < len(arbolRoot.Avd_ap_array_subdirectorios); i++ {
		var data = arbolRoot.Avd_ap_array_subdirectorios[i]
		if data == -1 {
			insideRoot = insideRoot + "\n<tr><td>AVD apun</td><td port='port'>-1</td></tr>"
		} else {
			insideRoot = insideRoot + "\n<tr><td>AVD apun</td><td port='port" + strconv.Itoa(int(data)) + "'>" + strconv.Itoa(int(data)) + "</td></tr>"
			enlacesRoot = enlacesRoot + "\n" + name + ":port" + strconv.Itoa(int(data)) + "   -> child" + strconv.Itoa(int(data)) + ";"
			childRoot = childRoot + GetAlDirectory(data, "child"+strconv.Itoa(int(data)), file)
		}
	}

	/// A DIRECOTRY
	var detalle = ""
	var detallebit = arbolRoot.Avd_ap_detalle_directorio
	if detallebit != -1 {
		//insideRoot = insideRoot + "\n<td port='port" + strconv.Itoa(int(detallebit)) + "'>" + strconv.Itoa(int(detallebit)) + "</td>"
		insideRoot = insideRoot + "\n<tr><td>DD apun</td><td port='port" + strconv.Itoa(int(detallebit)) + "'>" + strconv.Itoa(int(detallebit)) + "</td></tr>"

		detalle = detalle + GetAlDetails(detallebit, "child"+strconv.Itoa(int(detallebit)), file)
		enlacesRoot = enlacesRoot + "\n" + name + ":port" + strconv.Itoa(int(detallebit)) + "   -> child" + strconv.Itoa(int(detallebit)) + ";"
	}

	// APUNTADOR INDIRECTO
	if arbolRoot.Avd_ap_arbol_virtual_directorio != -1 {
		var dt = arbolRoot.Avd_ap_arbol_virtual_directorio
		//insideRoot = insideRoot + "\n<td port='port" + strconv.Itoa(int(dt)) + "'>I: " + strconv.Itoa(int(dt)) + "</td>"
		insideRoot = insideRoot + "\n<tr><td>AVD INDIR</td><td port='port" + strconv.Itoa(int(dt)) + "'>I:" + strconv.Itoa(int(dt)) + "</td></tr>"

		enlacesRoot = enlacesRoot + "\n" + name + ":port" + strconv.Itoa(int(dt)) + "   -> child" + strconv.Itoa(int(dt)) + ";"
		childRoot = childRoot + GetAlDirectory(dt, "child"+strconv.Itoa(int(dt)), file)
	}

	root = root + insideRoot
	root = root + "\n</table>>];\n" + childRoot + enlacesRoot + detalle

	return root
}

func GetAlDetails(bitInicio int64, name string, file *os.File) string {
	var detalle = CONTROLLER.GETDetails(bitInicio, file) // ESTA FUNCION ESTA EN DIRECTORY DETAIL CONTROLLER
	var body = "\n" + name + " [ shape=plaintext label=<  <table border='1' cellborder='1'> " +
		" <tr><td bgcolor=\"cadetblue1\" colspan=\"4\">Detail</td><td bgcolor=\"cadetblue1\" colspan=\"4\">" + strconv.Itoa(int(detalle.DD_num)) + "</td></tr>"

	var interior = ""
	var inodos = ""
	var apuntadores = ""
	var childs = ""
	if detalle.DD_file_lleno == true {
		interior = interior + "\n<tr>\n<td colspan=\"4\">" + GetAllName(detalle.DD_file_nombre) + "</td>\n<td colspan=\"4\" port='port" + strconv.Itoa(int(detalle.DD_file_ap_inodo)) + "'>" +
			strconv.Itoa(int(detalle.DD_file_ap_inodo)) + "</td>\n</tr>\n"

		if detalle.DD_file_ap_inodo != -1 {
			inodos = inodos + GetAlInodes(detalle.DD_file_ap_inodo, "child"+strconv.Itoa(int(detalle.DD_file_ap_inodo)), file)
			apuntadores = apuntadores + "\n" + name + ":port" + strconv.Itoa(int(detalle.DD_file_ap_inodo)) + "   -> child" + strconv.Itoa(int(detalle.DD_file_ap_inodo)) + ";"
		}

	}
	if detalle.DD_file_lleno2 == true {
		interior = interior + "\n<tr>\n<td colspan=\"4\">" + GetAllName(detalle.DD_file_nombre2) + "</td>\n<td colspan=\"4\" port='port" + strconv.Itoa(int(detalle.DD_file_ap_inodo2)) + "'>" +
			strconv.Itoa(int(detalle.DD_file_ap_inodo2)) + "</td>\n</tr>\n"

		if detalle.DD_file_ap_inodo2 != -1 {
			inodos = inodos + GetAlInodes(detalle.DD_file_ap_inodo2, "child"+strconv.Itoa(int(detalle.DD_file_ap_inodo2)), file)
			apuntadores = apuntadores + "\n" + name + ":port" + strconv.Itoa(int(detalle.DD_file_ap_inodo2)) + "   -> child" + strconv.Itoa(int(detalle.DD_file_ap_inodo2)) + ";"
		}
	}
	if detalle.DD_file_lleno3 == true {
		interior = interior + "\n<tr>\n<td colspan=\"4\">" + GetAllName(detalle.DD_file_nombre3) + "</td>\n<td colspan=\"4\" port='port" + strconv.Itoa(int(detalle.DD_file_ap_inodo3)) + "'>" +
			strconv.Itoa(int(detalle.DD_file_ap_inodo3)) + "</td>\n</tr>\n"

		if detalle.DD_file_ap_inodo3 != -1 {
			inodos = inodos + GetAlInodes(detalle.DD_file_ap_inodo3, "child"+strconv.Itoa(int(detalle.DD_file_ap_inodo3)), file)
			apuntadores = apuntadores + "\n" + name + ":port" + strconv.Itoa(int(detalle.DD_file_ap_inodo3)) + "   -> child" + strconv.Itoa(int(detalle.DD_file_ap_inodo3)) + ";"
		}

	}
	if detalle.DD_file_lleno4 == true {
		interior = interior + "\n<tr>\n<td colspan=\"4\">" + GetAllName(detalle.DD_file_nombre4) + "</td>\n<td colspan=\"4\" port='port" + strconv.Itoa(int(detalle.DD_file_ap_inodo4)) + "'>" +
			strconv.Itoa(int(detalle.DD_file_ap_inodo4)) + "</td>\n</tr>\n"

		if detalle.DD_file_ap_inodo4 != -1 {
			inodos = inodos + GetAlInodes(detalle.DD_file_ap_inodo4, "child"+strconv.Itoa(int(detalle.DD_file_ap_inodo4)), file)
			apuntadores = apuntadores + "\n" + name + ":port" + strconv.Itoa(int(detalle.DD_file_ap_inodo4)) + "   -> child" + strconv.Itoa(int(detalle.DD_file_ap_inodo4)) + ";"
		}
	}
	if detalle.DD_file_lleno5 == true {
		interior = interior + "\n<tr>\n<td colspan=\"4\">" + GetAllName(detalle.DD_file_nombre5) + "</td>\n<td colspan=\"4\" port='port" + strconv.Itoa(int(detalle.DD_file_ap_inodo5)) + "'>" +
			strconv.Itoa(int(detalle.DD_file_ap_inodo5)) + "</td>\n</tr>\n"

		if detalle.DD_file_ap_inodo5 != -1 {
			inodos = inodos + GetAlInodes(detalle.DD_file_ap_inodo5, "child"+strconv.Itoa(int(detalle.DD_file_ap_inodo5)), file)
			apuntadores = apuntadores + "\n" + name + ":port" + strconv.Itoa(int(detalle.DD_file_ap_inodo5)) + "   -> child" + strconv.Itoa(int(detalle.DD_file_ap_inodo5)) + ";"
		}
	}

	if detalle.DD_ap_detalle_directorio != -1 {
		interior = interior + "\n<tr>\n<td colspan=\"4\">" + "Indirecto" + "</td>\n<td colspan=\"4\" port='port" + strconv.Itoa(int(detalle.DD_ap_detalle_directorio)) + "'>" +
			strconv.Itoa(int(detalle.DD_ap_detalle_directorio)) + "</td>\n</tr>\n"

		apuntadores = apuntadores + "\n" + name + ":port" + strconv.Itoa(int(detalle.DD_ap_detalle_directorio)) + "   -> child" + strconv.Itoa(int(detalle.DD_ap_detalle_directorio)) + ";"
		childs = childs + GetAlDetails(detalle.DD_ap_detalle_directorio, "child"+strconv.Itoa(int(detalle.DD_ap_detalle_directorio)), file)
	}

	body = body + interior

	body = body + "</table> >]" + inodos + apuntadores + childs
	return body
}

func GetAlInodes(bitInicio int64, name string, file *os.File) string {

	var inodo = CONTROLLER.GETInode(bitInicio, file)

	var body = "\n" + name + " [shape=plaintext label=<<table border='1' cellborder='1'>\n <tr><td bgcolor=\"chocolate1\" colspan=\"2\">Inodo</td></tr>"

	var enlaces = ""
	var bloques = ""
	var childs = ""
	body = body + "\n<tr><td colspan=\"2\"> Numero:" + strconv.Itoa(int(inodo.I_count_inodo)) + "</td></tr>"
	body = body + "\n<tr><td colspan=\"2\"> Size File:" + strconv.Itoa(int(inodo.I_size_archivo)) + "</td></tr>"
	body = body + "\n<tr><td colspan=\"2\"> B. Asig:" + strconv.Itoa(int(inodo.I_count_bloques_asignados)) + "</td></tr>"

	for i := 0; i < len(inodo.I_array_bloques); i++ {
		if inodo.I_array_bloques[i] != 0 {
			body = body + " <tr><td>Block</td> <td port='port" + strconv.Itoa(int(inodo.I_array_bloques[i])) + "'> " + strconv.Itoa(int(inodo.I_array_bloques[i])) + " </td> " + "</tr>"
			enlaces = enlaces + "\n" + name + ":port" + strconv.Itoa(int(inodo.I_array_bloques[i])) + "   -> child" + strconv.Itoa(int(inodo.I_array_bloques[i])) + ";"
			bloques = bloques + GetAlBloques(inodo.I_array_bloques[i], "child"+strconv.Itoa(int(inodo.I_array_bloques[i])), file)

		}
	}
	if inodo.I_ap_indirecto != -1 {
		body = body + " <tr><td>Indirecto </td> <td port='port" + strconv.Itoa(int(inodo.I_ap_indirecto)) + "'> " + strconv.Itoa(int(inodo.I_ap_indirecto)) + " </td> " + "</tr>"
		enlaces = enlaces + "\n" + name + ":port" + strconv.Itoa(int(inodo.I_ap_indirecto)) + "   -> child" + strconv.Itoa(int(inodo.I_ap_indirecto)) + ";"
		childs = childs + GetAlInodes(inodo.I_ap_indirecto, "child"+strconv.Itoa(int(inodo.I_ap_indirecto)), file)
	}

	body = body + "</table>>];" + enlaces + bloques + childs
	return body
}

func GetAlBloques(bitInicio int64, name string, file *os.File) string {
	var bloque, encontrado = CONTROLLER.GetBlockOcupado(bitInicio, file)

	var body = "\n" + name + " [ shape=plaintext label=<<table border='1' cellborder='1'>"

	if encontrado {
		body = body + "<tr><td bgcolor=\"deepskyblue1\" colspan=\"2\">Block</td></tr>"

		body = body + "<tr><td colspan=\"2\">" + GetContentBlock(bloque.DB_data) + "</td></tr>"

	}
	body = body + "</table> >];"
	return body
}

///======================TREE FILE

func TreeFileReport(path, id, ruta string) {

	if SearchPartitionById(id) { //VOY A BUSCAR LA PARTICION MONTADA, ESTE METODO ESTA EN MOUNT_UMOUNT.GO
		var partition = GetPartitionById(id) //Obtengo la particion montada, ESTE METODO ESTA EN MOUNT_UMOUNT.GO

		//SE ABRE EL ARCHIVO
		file, err := os.OpenFile(partition.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error, no se pudo abrir el disco duro")
		}
		//OBTENGO EL SUPER BOOT
		var sb = GetSuperBoot(partition.Mount_part.Part_start, file) //GET SUPER BOOT SE ENCUENTRA FORMAT_FIRSTTIME.GO

		//ITERO EN EL BITMAP DE DETALLES PARA SABER CUANTOS HAY Y EN QUE POSICION DEBO CREAR EL NUEVO BLOQUE
		file.Seek(sb.SB_ap_bitmap_detail_dir, 0)
		b1 := make([]byte, (sb.SB_ap_detail_dir - sb.SB_ap_bitmap_detail_dir))
		n1, err := file.Read(b1)
		if err != nil {

		}
		var contador int = 0
		for i := 0; i < len(string(b1[:n1])); i++ {
			if string(b1[:n1][i]) == "1" {
				contador = contador + 1
			}
		}

		//var fileNames = ""

		GetFileNamesInDetails(sb.SB_ap_detail_dir, contador, file)
		/*//OBTENGO EL ARBOL DE DIRECTORIOS
		var arbolRoot = CONTROLLER.GetArbolVirual(sb.SB_ap_tree_dir, file)
		var body = "digraph H { rankdir=\"LR\" "

		//ROOT
		var root = "parent [ shape=plaintext \n label=<\n<table border='1' cellborder='1'> \n<tr><td colspan=\"8\">" + GetAllName(arbolRoot.Avd_nombre_directorio) + "</td></tr> " +
			"\n<tr><td colspan=\"8\">Proper: " + GetAllName(arbolRoot.Avd_proper) + "</td></tr>\n" + "\n<tr><td colspan=\"8\">" + GetFecha(arbolRoot.Avd_fecha_creacion) + "</td></tr>\n<tr>"

		var insideRoot = ""
		var enlacesRoot = ""
		var childRoot = ""

		for i := 0; i < len(arbolRoot.Avd_ap_array_subdirectorios); i++ {
			var data = arbolRoot.Avd_ap_array_subdirectorios[i]
			if data == -1 {
				insideRoot = insideRoot + "\n<td port='port'>-1</td>"
			} else {
				insideRoot = insideRoot + "\n<td port='port" + strconv.Itoa(int(data)) + "'>" + strconv.Itoa(int(data)) + "</td>"
				enlacesRoot = enlacesRoot + "\nparent:port" + strconv.Itoa(int(data)) + "   -> child" + strconv.Itoa(int(data)) + ";"
				childRoot = childRoot + GetAlDirectory(data, "child"+strconv.Itoa(int(data)), file)
			}
		}

		/// A DIRECOTRY
		var detalle = ""
		fmt.Println("---------")
		var detallebit = arbolRoot.Avd_ap_detalle_directorio
		if detallebit != -1 {
			insideRoot = insideRoot + "\n<td port='port" + strconv.Itoa(int(detallebit)) + "'>" + strconv.Itoa(int(detallebit)) + "</td>"
			detalle = detalle + GetAlDetails(detallebit, "child"+strconv.Itoa(int(detallebit)), file)
			enlacesRoot = enlacesRoot + "\nparent:port" + strconv.Itoa(int(detallebit)) + "   -> child" + strconv.Itoa(int(detallebit)) + ";"
		}

		// APUNTADOR INDIRECTO
		if arbolRoot.Avd_ap_arbol_virtual_directorio != -1 {
			var dt = arbolRoot.Avd_ap_arbol_virtual_directorio
			insideRoot = insideRoot + "\n<td port='port" + strconv.Itoa(int(dt)) + "'>" + strconv.Itoa(int(dt)) + "</td>"
			enlacesRoot = enlacesRoot + "\nparent:port" + strconv.Itoa(int(dt)) + "   -> child" + strconv.Itoa(int(dt)) + ";"
			childRoot = childRoot + GetAlDirectory(dt, "child"+strconv.Itoa(int(dt)), file)
		}

		root = root + insideRoot
		root = root + "\n</tr></table>>];"
		body = body + root + childRoot + enlacesRoot + detalle

		body = body + "\n}"

		dir, name := filepath.Split(path)

		GeneratePNG(name, body, dir)*/
	} else {
		fmt.Println("La particion con ID: " + id + " no esta montada")
	}
}

func GetFileNamesInDetails(bitInicio int64, contador int, file *os.File) {
	fmt.Println("========================")
	fmt.Println("> Listado de Archivos")
	fmt.Println("========================")

	for i := 0; i < contador; i++ {

		var detalle = CONTROLLER.GETDetails((bitInicio + int64(i)*int64(unsafe.Sizeof(STRUCTURES.DIRECTORYDETAIL{}))), file)
		fmt.Println("." + GetAllName(detalle.DD_file_nombre))
	}
}

//se guarda la direccion del proyecto
var rootDir = FUNCTION.RootDir()

func CreateMBRReport(mbr STRUCTURES.MBR) {

	//Se crea los directorios que almacenara los dots
	FUNCTION.CreateADirectory(rootDir + "/reports/dots")
	//Se crea los directorios que almacenara las imagenes
	FUNCTION.CreateADirectory(rootDir + "/reports/pngs")

	var body string = "digraph test { graph [ratio=fill];" +
		"node [label=\"Grafica\", fontsize=15, shape=plaintext];" +
		"graph [bb=\"0,0,352,154\"];" +
		"arset [label=<" +
		"<TABLE>" +
		"<TR>" + "<TD>Mbr_size</TD>" + "<TD>" + strconv.Itoa(int(mbr.Mbr_size)) + "</TD>" + "</TR>" +
		"<TR>" + "<TD>Mbr_creation_date</TD>" + "<TD>" + string(mbr.Mbr_creation_date[:]) + "</TD>" + "</TR>" +
		"<TR>" + "<TD>Mbr_disk_signature</TD>" + "<TD>" + strconv.Itoa(int(mbr.Mbr_disk_signature)) + "</TD>" + "</TR>"

	body = verifyPartition(mbr.Mbr_partition_1, body, 1)
	body = verifyPartition(mbr.Mbr_partition_2, body, 2)
	body = verifyPartition(mbr.Mbr_partition_3, body, 3)
	body = verifyPartition(mbr.Mbr_partition_4, body, 4)

	body = body + "</TABLE>" +
		">, ];}"

	GeneratePNG("mbr.jpg", body, rootDir+"/reports/dots/")

}

func verifyPartition(partition STRUCTURES.PARTITION, body string, size int) string {

	if partition.Part_isEmpty != 0 {
		var s string
		for _, v := range partition.Part_name {
			if v != 0 {
				s = s + string(v)
			}
		}
		body = body + "<TR>" + "<TD>" + s + "</TD>" + "</TR>" +
			"<TR>" + "<TD>Partition Number</TD>" + "<TD>" + strconv.Itoa(size) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>Part_status</TD>" + "<TD>" + strconv.Itoa(int(partition.Part_status)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>Part_type</TD>" + "<TD>" + string(partition.Part_type) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>Part_fit</TD>" + "<TD>" + string(partition.Part_fit) + "</TD>" + "</TR>" +

			"<TR>" + "<TD>Part_start</TD>" + "<TD>" + strconv.Itoa(int(partition.Part_start)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>Part_size</TD>" + "<TD>" + strconv.Itoa(int(partition.Part_size)) + "</TD>" + "</TR>"
	}
	return body
}

func CreateDiskReport(mbr STRUCTURES.MBR, diskName string) {
	//Se crea los directorios que almacenara los dots
	FUNCTION.CreateADirectory(rootDir + "/reports/dots")
	//Se crea los directorios que almacenara las imagenes
	FUNCTION.CreateADirectory(rootDir + "/reports/pngs")

	var tableSize = mbr.Mbr_disk / 100

	var body string = "digraph {  " +
		"tbl [ " +
		"shape=plaintext " +
		"label=<" +
		"<table WIDTH=\"" + strconv.Itoa(int(tableSize)) + "\" >"

	var heads = "<tr><td></td>"

	var PartitionArray = [4]STRUCTURES.PARTITION{}
	PartitionArray[0] = mbr.Mbr_partition_1
	PartitionArray[1] = mbr.Mbr_partition_2
	PartitionArray[2] = mbr.Mbr_partition_3
	PartitionArray[3] = mbr.Mbr_partition_4

	var totalDiscoOcupado = mbr.Mbr_size

	///SIRVE PARA PONER LOS ENCABEZADOS
	var temporal = PartitionArray[0].Part_start
	for i := 0; i < len(PartitionArray); i++ {
		var part = PartitionArray[i]

		var libreTemporal = part.Part_start - temporal
		if libreTemporal > 0 {

			heads = heads + "<td> </td>"
		}

		temporal = part.Part_end

		if part.Part_isEmpty == 1 && string(part.Part_type) == "P" {
			heads = heads + "<td>Primaria</td>"
		} else if string(part.Part_type) == "E" {
			heads = heads + "<td>Extendida</td>"
		}
	}
	heads = heads + "</tr>"
	///////////////

	body = body + heads

	/// COLOCA EL MBR
	body = body + "<tr><td WIDTH=\"" + strconv.Itoa(int(mbr.Mbr_size)) + "\" >MBR: " + strconv.Itoa(int(mbr.Mbr_size)) + "</td>"

	if mbr.Mbr_partition_1.Part_isEmpty == 0 &&
		mbr.Mbr_partition_2.Part_isEmpty == 0 &&
		mbr.Mbr_partition_3.Part_isEmpty == 0 &&
		mbr.Mbr_partition_4.Part_isEmpty == 0 {

		var libre1 int64 = (mbr.Mbr_disk - mbr.Mbr_size) / 100
		body = body + "<td WIDTH=\"" + strconv.Itoa(int(libre1)) + "\" >Libre: " + strconv.Itoa(int(libre1*100)) + " </td>"
	} else {
		var libreEntreParticiones = mbr.Mbr_size

		for i := 0; i < 4; i++ {
			var part = PartitionArray[i]
			if part.Part_isEmpty == 1 {

				/////////// ME SIRVE PARA VER SI HAY FRAGMENTACION, ESPACIO ENTRE PARTICIONES
				var libreTemporal = part.Part_start - libreEntreParticiones
				if libreTemporal > 0 {

					body = body + "<td WIDTH=\"" + strconv.Itoa(int(libreTemporal/100)) + "\" >Libre: " + strconv.Itoa(int(libreTemporal)) + " </td>"
				}
				///////////////////////////////////////////

				var part1len = part.Part_size / 100

				//VERIFICO LAS PARTICIONES PRIMARIAS
				if string(part.Part_type) == "P" {
					body = body + "<td WIDTH=\"" + strconv.Itoa(int(part1len)) + "\" >" + GetString1(part.Part_name) + " " + strconv.Itoa(int(part.Part_size)) + "</td>"
					totalDiscoOcupado = totalDiscoOcupado + part.Part_size

					//VERIFICO LAS PARTICIONES EXTENDIDAS
				} else { // PARTICION EXTENDIDA

					body = body + "<td WIDTH=\"" + strconv.Itoa(int(part1len)) + "\" ><table><tr>" // CREO LA TABLA

					var extended, encontrado = CONTROLLER.GetExtendedPartition(diskName) //ESTE METODO ESTA EN STRUCT CONTROLLER,
					//ME RETORNA LAS LOGICAS CREADASE EN ESA PARTICION

					if encontrado { // SI ENCUENTRA QUE TIENE EXTENDIDA

						var espacioExtendidoLIbre int64 = part.Part_size

						if len(extended.Part_ebr) == 1 && len(extended.Part_partition) == 0 {

							body = body + "<td WIDTH=\"" + strconv.Itoa(int(extended.Part_ebr[0].Part_size)) + "\" > EBR:" + strconv.Itoa(int(extended.Part_ebr[0].Part_size)) + "</td >"
							espacioExtendidoLIbre = espacioExtendidoLIbre - extended.Part_ebr[0].Part_size
						} else {

							for i := 0; i < len(extended.Part_ebr); i++ {

								if extended.Part_partition[i].Part_isEmpty == 1 {
									var ebrAuxiliar = extended.Part_ebr[i]
									body = body + "<td WIDTH=\"" + strconv.Itoa(int(ebrAuxiliar.Part_size)) + "\" > " + GetString1(ebrAuxiliar.Part_name) + ":" + strconv.Itoa(int(ebrAuxiliar.Part_size)) + "</td >"
									var logicaAuxiliar = extended.Part_partition[i]
									var lenlogic = logicaAuxiliar.Part_size / 100
									body = body + "<td WIDTH=\"" + strconv.Itoa(int(lenlogic)) + "\"> " + GetString1(logicaAuxiliar.Part_name) + " :" + strconv.Itoa(int(logicaAuxiliar.Part_size)) + "</td >"

									espacioExtendidoLIbre = espacioExtendidoLIbre - ebrAuxiliar.Part_size - logicaAuxiliar.Part_size

								}
							}

						}
						body = body + "<td WIDTH=\"" + strconv.Itoa(int(espacioExtendidoLIbre/100)) + "\" > Ext Libre" + strconv.Itoa(int(espacioExtendidoLIbre)) + "</td >"
						CreateEBRReport(extended)
					}

					body = body + "</tr></table></td>"

					totalDiscoOcupado = totalDiscoOcupado + part.Part_size
				}
				libreEntreParticiones = part.Part_end

			}
		}

		if totalDiscoOcupado < mbr.Mbr_disk {
			var libre1 int64 = (mbr.Mbr_disk - totalDiscoOcupado) / 100
			body = body + "<td WIDTH=\"" + strconv.Itoa(int(libre1)) + "\" >Libre: " + strconv.Itoa(int(libre1*100)) + " </td>"
		}
	}

	body = body + "</tr></table>>];}"

	GeneratePNG("disk.jpg", body, rootDir+"/reports/dots/")
}

func CreateEBRReport(extended STRUCTURES.EXTENDED) {

	var espacioExtendidoLIbre int64 = extended.Part_size
	var body = "digraph H {" +
		"rankdir=\"LR\" " +
		"parent [ shape=plaintext " +
		"label=<<table border='1' cellborder='1'> "
	if len(extended.Part_ebr) == 1 && len(extended.Part_partition) == 0 {

		body = body + "<tr><td colspan=\"2\">EBR PARTITIONS</td></tr><tr>"
		body = body + "<td WIDTH=\"" + strconv.Itoa(int(extended.Part_ebr[0].Part_size)) + "\" > EBR:" + strconv.Itoa(int(extended.Part_ebr[0].Part_size)) + "</td >"
		espacioExtendidoLIbre = espacioExtendidoLIbre - extended.Part_ebr[0].Part_size

	} else {

		var leng = strconv.Itoa(len(extended.Part_ebr)*2 + 1)
		body = body + "<tr><td colspan=\"" + leng + "\">EBR PARTITIONS</td></tr><tr>"

		for i := 0; i < len(extended.Part_ebr); i++ {

			if extended.Part_partition[i].Part_isEmpty == 1 {
				var ebrAuxiliar = extended.Part_ebr[i]
				body = body + "<td WIDTH=\"" + strconv.Itoa(int(ebrAuxiliar.Part_size)) + "\" > " + GetString1(ebrAuxiliar.Part_name) + ":" + strconv.Itoa(int(ebrAuxiliar.Part_size)) + "</td >"
				var logicaAuxiliar = extended.Part_partition[i]

				var lenlogic = logicaAuxiliar.Part_size / 100
				body = body + "<td WIDTH=\"" + strconv.Itoa(int(lenlogic)) + "\"> " + GetString1(logicaAuxiliar.Part_name) + " :" + strconv.Itoa(int(logicaAuxiliar.Part_size)) + "</td >"

				espacioExtendidoLIbre = espacioExtendidoLIbre - ebrAuxiliar.Part_size - logicaAuxiliar.Part_size

			}
		}

	}
	body = body + "<td WIDTH=\"" + strconv.Itoa(int(espacioExtendidoLIbre/100)) + "\" > Ext Libre: " + strconv.Itoa(int(espacioExtendidoLIbre)) + "</td >"

	body = body + "</tr></table>>];}"

	GeneratePNG("ebr.jpg", body, rootDir+"/reports/dots/")
}

func GeneratePNG(nombre string, body string, path string) {

	FUNCTION.CreateADirectory(path)

	if FUNCTION.IfExistFile(path + nombre + ".dot") {
		err := os.Remove(path + nombre + ".dot")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	f, err := os.Create(path + nombre + ".dot")
	if err != nil {
		fmt.Println(err)
		return
	}

	f.WriteString(body)

	app := "dot -Tjpg "
	arg1 := "\"" + path + nombre + ".dot" + "\""
	arg2 := " -o "
	arg3 := " \"" + path + nombre + "\""

	err, out, errout := Shellout(app + arg1 + arg2 + arg3)
	if err != nil {
		log.Printf("error: %v\n", err)
	} else {
		fmt.Println(out)
		fmt.Println(errout)

	}
}

const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func GetString1(str [16]byte) string {
	var cadena = ""
	for i := 0; i < len(str); i++ {
		if string(str[i]) == "+" {
			break
		}
		cadena = cadena + string(str[i])
	}
	return cadena
}

func GetAllName(str [20]byte) string {
	var cadena = ""
	for i := 0; i < len(str); i++ {
		if string(str[i]) == "+" {
			break
		}
		cadena = cadena + string(str[i])
	}
	return cadena
}
func GetContentBlock(str [25]byte) string {
	var s string
	for _, v := range str {
		if v != 0 {
			s = s + string(v)
		}
	}
	return s
}

func GetFecha(str [19]byte) string {
	var cadena = ""
	for i := 0; i < len(str); i++ {
		if string(str[i]) == "+" {
			break
		}
		cadena = cadena + string(str[i])
	}
	return cadena
}
