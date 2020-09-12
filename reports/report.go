package reports

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	CONTROLLER "../controllers"
	FUNCTION "../functions"
	STRUCTURES "../structures"
)

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

	GeneratePNG("mbr", body)

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

	var body string = "digraph {" +
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
					body = body + "<td WIDTH=\"" + strconv.Itoa(int(part1len)) + "\" >" + GetString(part.Part_name) + " " + strconv.Itoa(int(part.Part_size)) + "</td>"
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
									body = body + "<td WIDTH=\"" + strconv.Itoa(int(ebrAuxiliar.Part_size)) + "\" > " + GetString(ebrAuxiliar.Part_name) + ":" + strconv.Itoa(int(ebrAuxiliar.Part_size)) + "</td >"
									var logicaAuxiliar = extended.Part_partition[i]
									var lenlogic = logicaAuxiliar.Part_size / 100
									body = body + "<td WIDTH=\"" + strconv.Itoa(int(lenlogic)) + "\"> " + GetString(logicaAuxiliar.Part_name) + " :" + strconv.Itoa(int(logicaAuxiliar.Part_size)) + "</td >"

									espacioExtendidoLIbre = espacioExtendidoLIbre - ebrAuxiliar.Part_size - logicaAuxiliar.Part_size

								}
							}
						}
						body = body + "<td WIDTH=\"" + strconv.Itoa(int(espacioExtendidoLIbre/100)) + "\" > Ext Libre" + strconv.Itoa(int(espacioExtendidoLIbre)) + "</td >"

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

	GeneratePNG("disk", body)
}

func GeneratePNG(nombre string, body string) {

	if FUNCTION.IfExistFile(rootDir + "/reports/dots/" + nombre + ".dot") {
		err := os.Remove(rootDir + "/reports/dots/" + nombre + ".dot")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	f, err := os.Create(rootDir + "/reports/dots/" + nombre + ".dot")
	if err != nil {
		fmt.Println(err)
		return
	}

	f.WriteString(body)

	app := "dot -Tpng "
	arg1 := "\"" + rootDir + "/reports/dots/" + nombre + ".dot" + "\""
	arg2 := " -o "
	arg3 := " \"" + rootDir + "/reports/pngs/" + nombre + ".png" + "\""

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

func GetString(str [16]byte) string {
	var cadena = ""
	for i := 0; i < len(str); i++ {
		if string(str[i]) == "+" {
			break
		}
		cadena = cadena + string(str[i])
	}
	return cadena
}
