package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"unsafe"

	STRUCTURES "../structures"
)

//	FUNCION PARA LEER ARCHIVOS
func ReadFile(path string) {

	/*file2, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file2.Close()
	if err != nil {
		fmt.Println(err)
	}

	stats, statsErr := file2.Stat()
	if statsErr != nil {
		fmt.Println("erro")
	}

	var sizes int64 = stats.Size()
	bytess := make([]byte, sizes)

	bufr := bufio.NewReader(file2)
	_, err = bufr.Read(bytess)

	fmt.Println(bytess)*/

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
	fmt.Printf("SIZE: %s\nFECHA: %s\nSIGNATURE: %s\nPARTITIONS: %s\n", strconv.Itoa(int(m.Mbr_size)), m.Mbr_creation_date, strconv.Itoa(int(m.Mbr_disk_signature)),
		strconv.Itoa(int(m.Mbr_count)))
	fmt.Println("Partition1:")

	s := reflect.ValueOf(&m.Mbr_partition_1).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s = %v\n",
			typeOfT.Field(i).Name, f.Interface())
	}

	fmt.Println("Partition2:")
	s = reflect.ValueOf(&m.Mbr_partition_2).Elem()
	typeOfT = s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s = %v\n",
			typeOfT.Field(i).Name, f.Interface())
	}

	fmt.Println("Partition3:")
	s = reflect.ValueOf(&m.Mbr_partition_3).Elem()
	typeOfT = s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s = %v\n",
			typeOfT.Field(i).Name, f.Interface())
	}

	fmt.Println("Partition4:")
	s = reflect.ValueOf(&m.Mbr_partition_4).Elem()
	typeOfT = s.Type()

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
