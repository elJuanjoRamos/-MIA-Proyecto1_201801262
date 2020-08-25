package reports

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	CONTROLLER "../controller"
	STRUCTURES "../structures"
)

//se guarda la direccion del proyecto
var rootDir = CONTROLLER.RootDir()

func CreateMBRReport(mbr STRUCTURES.MBR) {
	//Se crea los directorios que almacenara los dots
	CONTROLLER.CreateADirectory(rootDir + "/reports/dots")
	//Se crea los directorios que almacenara las imagenes
	CONTROLLER.CreateADirectory(rootDir + "/reports/pngs")

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

	f, err := os.Create(rootDir + "/reports/dots/mbr.dot")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(body)

	app := "dot"
	arg0 := "-Tpng"
	arg1 := "\"" + rootDir + "/reports/dots/mbr.dot" + "\""
	arg2 := "-o"
	arg3 := "\"" + rootDir + "/reports/pngs/mbr.png" + "\""

	fmt.Println(app + " " + arg0 + " " + arg1 + " " + arg2 + " " + arg3)
	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))

}

func verifyPartition(partition STRUCTURES.PARTITION, body string, size int) string {
	if partition.Part_active {

		var s string
		for _, v := range partition.Part_name {
			if v != 0 {
				s = s + string(v)
			}
		}
		body = body + "<TR>" + "<TD>" + s + "</TD>" + "</TR>" +
			"<TR>" + "<TD>Partition Number</TD>" + "<TD>" + strconv.Itoa(size) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>Part_status</TD>" + "<TD>" + string(partition.Part_status) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>Part_type</TD>" + "<TD>" + string(partition.Part_type) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>Part_fit</TD>" + "<TD>" + string(partition.Part_fit) + "</TD>" + "</TR>" +

			"<TR>" + "<TD>Part_start</TD>" + "<TD>" + strconv.Itoa(int(partition.Part_start)) + "</TD>" + "</TR>" +
			"<TR>" + "<TD>Part_size</TD>" + "<TD>" + strconv.Itoa(int(partition.Part_size)) + "</TD>" + "</TR>"

	}
	return body
}
