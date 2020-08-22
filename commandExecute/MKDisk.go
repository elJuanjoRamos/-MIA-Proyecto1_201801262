package commandExecute

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	FUNCTIONCONTROLLER "../controller"
	STRUCTURES "../structures"
)

//	FUNCION PARA LEER ARCHIVOS
func ReadFile(path string) {

	//Abrimos/creamos un archivo.
	file, err := os.Open(path)
	defer file.Close()
	if err != nil { //validar que no sea nulo.
		log.Fatal(err)
	}

	//Declaramos variable de tipo mbr
	m := STRUCTURES.MBR{}
	//Obtenemos el tamanio del mbr
	var size int = int(unsafe.Sizeof(m))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	//Se imprimen los valores guardados en el struct
	fmt.Println(m)
	fmt.Printf("SIZE: %s\nFECHA: %s\nSIGNATURE: %s\n", strconv.Itoa(int(m.Mbr_size)), m.Mbr_creation_date, strconv.Itoa(int(m.Mbr_disk_signature)))
	fmt.Println("Partition:")

	s := reflect.ValueOf(&m.Mbr_partition_1).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s = %v\n",
			typeOfT.Field(i).Name, f.Interface())
	}
}

//Función que lee del archivo, se especifica cuantos bytes se quieren leer.
func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes

	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

//Método para escribir en un archivo
func WriteFile(name string, path string, size int64) {

	//Mando a crear el directorio
	FUNCTIONCONTROLLER.CreateADirectory(path)

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

	/*disco.Mbr_partition_1 = STRUCTURES.PARTITION{
	Part_status: 'F',
	Part_type:   'F',
	Part_fit:    'F',
	Part_start:  30,
	Part_size:   1024}*/

	//Escribimos struct.
	s1 := &disco
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

}

//Método para escribir en un archivo.
func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}
