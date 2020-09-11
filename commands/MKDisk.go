package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	CONTROLLER "../controllers"
	FUNCTION "../functions"
	REPORTS "../reports"
	STRUCTURES "../structures"
)

//Método para escribir en un archivo
func CreateFile(name string, path string, size int64) {
	//Mando a crear el directorio
	FUNCTION.CreateADirectory(path)
	//Se crea el archivo
	file, err := os.Create(path + name)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//Se manda a guardar el nombre del disco por si se crean particiones extendidas
	CONTROLLER.AddNewExtendedDisk(name)
	var init int8 = 0

	s := &init

	//Escribimos un 0 en el inicio del archivo.
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	escribirBytes(file, binario.Bytes())

	//Nos posicionamos en el byte final(primera posicion es 0)
	file.Seek(size-1, 0) // segundo parametro: 0, 1, 2.     0 -> Inicio, 1-> desde donde esta el puntero, 2 -> Del fin para atras

	//Escribimos un 0 al final del archivo.
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	escribirBytes(file, binario2.Bytes())

}

func WriteFile(path string) {

	//Abrimos el archivo
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	//----------------------------------------------------------------------- //
	//Escribimos nuestro struct en el inicio del archivo

	file.Seek(0, 0) // nos posicionamos en el inicio del archivo.

	var random int64 = rand.Int63()
	var time = time.Now()

	//Asignamos valores a los atributos del struct.
	disco := STRUCTURES.MBR{Mbr_disk_signature: random, Mbr_count: 4, Mbr_size: FUNCTION.FileSize(path)}
	copy(disco.Mbr_creation_date[:], time.Format("2006-01-02 15:04:05"))

	//Escribimos struct.
	s1 := &disco
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

	REPORTS.CreateMBRReport(disco)
}

//Método para escribir en un archivo.
func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}
