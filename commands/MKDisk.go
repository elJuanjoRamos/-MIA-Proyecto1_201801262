package commands

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"os"
	"time"
	"unsafe"

	FUNCTION "../functions"
	STRUCTURES "../structures"
)

//Método para escribir en un archivo
func WriteFile(name string, path string, size int64) {

	//Mando a crear el directorio
	FUNCTION.CreateADirectory(path)
	//Se crea el archivo
	file, err := os.Create(path + name)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var otro int8 = 0

	s := &otro

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

	//----------------------------------------------------------------------- //
	//Escribimos nuestro struct en el inicio del archivo

	file.Seek(0, 0) // nos posicionamos en el inicio del archivo.

	var random int64 = rand.Int63()
	var time = time.Now()

	//Asignamos valores a los atributos del struct.
	disco := STRUCTURES.MBR{Mbr_disk_signature: random, Mbr_count: 4}
	copy(disco.Mbr_creation_date[:], time.Format("2006-01-02 15:04:05"))
	var sizeDisk int64 = int64(unsafe.Sizeof(disco))
	disco.Mbr_size = sizeDisk

	//Escribimos struct.
	s1 := &disco
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())
	//REPORTS.CreateMBRReport(disco)
}

//Método para escribir en un archivo.
func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}
