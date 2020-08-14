package commandExecute

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
)

type MBR struct {
	Number    uint8
	Character byte
	Sttring   [20]byte
}

//	FUNCION PARA LEER ARCHIVOS
func ReadFile() {
	//Abrimos/creamos un archivo.
	file, err := os.Open("test.bin")
	defer file.Close()
	if err != nil { //validar que no sea nulo.
		log.Fatal(err)
	}

	//Declaramos variable de tipo MBR
	m := MBR{}
	//Obtenemos el tamanio del MBR
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
	fmt.Printf("Character: %c\nSttring: %s\n", m.Character, m.Sttring)
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

//Método para escribir en un archivo, en este caso se escribe un archivo binario de 1kb, 1024 bytes.
func WriteFile(name string, size int64) {
	file, err := os.Create(name)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var otro int8 = 0

	s := &otro

	fmt.Println(unsafe.Sizeof(otro))
	//Escribimos un 0 en el inicio del archivo.
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	escribirBytes(file, binario.Bytes())
	//Nos posicionamos en el byte 1023 (primera posicion es 0)
	file.Seek(size, 0) // segundo parametro: 0, 1, 2.     0 -> Inicio, 1-> desde donde esta el puntero, 2 -> Del fin para atras

	//Escribimos un 0 al final del archivo.
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	escribirBytes(file, binario2.Bytes())

	//----------------------------------------------------------------------- //
	//Escribimos nuestro struct en el inicio del archivo

	file.Seek(0, 0) // nos posicionamos en el inicio del archivo.

	//Asignamos valores a los atributos del struct.
	disco := MBR{}
	//disco.Character = 'a'

	// Igualar Sttrings a array de bytes (array de chars)
	//cadenita := "Hola Amigos"
	//copy(disco.Sttring[:], cadenita)

	s1 := &disco

	//Escribimos struct.
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
