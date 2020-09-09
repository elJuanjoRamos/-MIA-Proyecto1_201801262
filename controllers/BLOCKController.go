package controllers

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"unsafe"

	STRUCTURES "../structures"
)

func BlockController_InsertText(sb STRUCTURES.SUPERBOOT, inicioInodo int64, texto string, file *os.File, idUsuario int64) {
	//OBTENEMOS EL INODO
	var m = GETInode(inicioInodo, file) // ESTA FUNCION ESTA EN INODOCONTROLLER.go
	//HACEMOS UN FOR EN EL ARREGLO DE BLOQUES PARA VERIFICAR QUE TENGA BLOQUES LIBRES
	var strtemp = ""
	for i := 0; i < len(m.I_array_bloques); i++ {
		var bloq, leido = GetBlockOcupado(m.I_array_bloques[i], file)
		if leido {
			strtemp = strtemp + GetStringInBlock(bloq.DB_data) //OBTENGO EL TEXTO DENTRO DEL BLOQUE
		}
	}

	if len(strtemp) < 100 { //SIGNIFICA QUE NO SE HAN OCUPADO TODOS LOS BLOQUES DEL INODO
		//HACEMOS UN FOR EN EL ARREGLO DE BLOQUES
		for i := 0; i < len(m.I_array_bloques); i++ {
			var bloq, leido = GetBlockOcupado(m.I_array_bloques[i], file)
			if leido { //SI FUE LEIDO, SIGNIFICA QUE EL BLOQUE YA ESTABA CREADO
				var str = GetStringInBlock(bloq.DB_data) //OBTENGO EL TEXTO DENTRO DEL BLOQUE
				if len(str) < 25 {                       //SI LA LONFITUD DEL BLOQUE ES MENOR A 25, SIGNIFICA QUE TODAVIA LE CABEN CARACTERES
					var cantidadCaracteres int = 25 - len(str)
					if len(texto) != 0 {

						if cantidadCaracteres <= len(texto) {
							for i := 0; i < cantidadCaracteres; i++ {
								str = str + string(texto[i])
							}
						} else {
							for i := 0; i < len(texto); i++ {
								str = str + string(texto[i])
							}
						}

						//MANDO A REESCRIBIR EL BLOQUE
						UpdateBlock(m.I_array_bloques[i], str, file)
						m.I_size_archivo = m.I_size_archivo + int64(len(str)) //UPDATE DEL SIZE DEL ARCHIVO

						var strtemp = ""
						for i := cantidadCaracteres; i < len(texto); i++ {
							strtemp = strtemp + string(texto[i])
						}
						texto = strtemp
					}
				}

			} else { //SI NO PUDO SER LEIDO, SIGNIFICA QUE EL BLOQUE NO ESTA CREADO
				var iterInString = 0
				var contador = 0
				var strTemp1 = ""
				var strTemp2 = ""

				if len(texto) != 0 {
					for j := 0; j < len(texto); j++ {
						strTemp1 = strTemp1 + string(texto[j])
						contador = contador + 1
						if contador == 25 || ((contador + iterInString) == int(len(texto))) {
							iterInString = iterInString + contador
							break
						}
					}
					m.I_array_bloques[i] = CreateABlock(sb, strTemp1, file)       //ASIGNO EL NUEVO BLOQUE EN EL ARRAY
					m.I_size_archivo = m.I_size_archivo + int64(len(strTemp1))    //UPDATE DEL SIZE DEL ARCHIVO
					m.I_count_bloques_asignados = m.I_count_bloques_asignados + 1 //AUMENTA LA CANTIDAD DE BLOQUES ASIGNADOS

					for i := len(strTemp1); i < len(texto); i++ {
						strTemp2 = strTemp2 + string(texto[i])
					}
					texto = strTemp2

				}
			}
		}

		if len(texto) > 0 { //SI ENTRA AQUI, SIGNIFICA QUE ENTRO AL PRIMER IF, ES DECIR, ENCONTRO ESPACIOS LIRBRES
			//DENTRO DE LOS BLOQUES, PERO TODAVIA SOBRO TEXTO QUE INSERTAR, ENTONCES ES NECESARIO CREAR BLOQUES INDIRECTOS
			if m.I_ap_indirecto == -1 { //SI NO TIENE UN APUNTADOR, LO CREA, SI YA TIENE UNO, LO MANDA
				m.I_ap_indirecto = INODOCONTROLLER_CreateINODO(sb, file, m.I_size_archivo, idUsuario)
				//METODO RECURSIVO PARA METER TEXTO DENTRO DEL INODO INDIRECTO
				BlockController_InsertText(sb, m.I_ap_indirecto, texto, file, idUsuario)
				//HAGO UN UPDATE DEL SIZE DEL ARCHIVO
				m.I_size_archivo = m.I_size_archivo + (GETInode(m.I_ap_indirecto, file)).I_size_archivo
			} else {
				BlockController_InsertText(sb, m.I_ap_indirecto, texto, file, idUsuario)
				m.I_size_archivo = m.I_size_archivo + (GETInode(m.I_ap_indirecto, file)).I_size_archivo
			}
		}

	} else { //SIGFINIFICA QUE YA NO YA ESPACIO LIBRE EN LOS BLOQUES ACTUALES Y HAY QUE CREAR UN INODO Y HACER REFERENCIA, (INODO INDIRECTO)
		if m.I_ap_indirecto == -1 { //SI NO TIENE UN APUNTADOR, LO CREA, SI YA TIENE UNO, LO MANDA
			m.I_ap_indirecto = INODOCONTROLLER_CreateINODO(sb, file, m.I_size_archivo, idUsuario)
			//METODO RECURSIVO PARA METER TEXTO DENTRO DEL INODO INDIRECTO
			BlockController_InsertText(sb, m.I_ap_indirecto, texto, file, idUsuario)
			//HAGO UN UPDATE DEL SIZE DEL ARCHIVO
			m.I_size_archivo = m.I_size_archivo + (GETInode(m.I_ap_indirecto, file)).I_size_archivo
		} else {
			BlockController_InsertText(sb, m.I_ap_indirecto, texto, file, idUsuario)
			m.I_size_archivo = m.I_size_archivo + (GETInode(m.I_ap_indirecto, file)).I_size_archivo
		}
	}

	//PREGAMOS LOS CAMBIOS HECHOS A LA TABLA INODO
	//NOS SITUAMOS AL INICIO DEL INODO
	file.Seek(inicioInodo, 0)
	//PEGAMOS OTRA VEZ EL INODO, ESTA VEZ, CON LA CANTIDAD DE BLOQUES NUEVOS
	s1 := &m
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

}

func GetBlockOcupado(bitInicio int64, file *os.File) (STRUCTURES.DATABLOCK, bool) {
	if bitInicio != 0 {
		file.Seek(bitInicio, 0)
		//LEEMOS LA ESTRUCTURA BLOQUE
		m1 := STRUCTURES.DATABLOCK{}
		//Obtenemos el tamanio del BLOQUE
		var size1 int = int(unsafe.Sizeof(m1))
		//Lee la cantidad de <size> bytes del archivo
		data1 := leerBytes(file, size1)
		//Convierte la data en un buffer,necesario para
		//decodificar binario
		buffer1 := bytes.NewBuffer(data1)

		//Decodificamos y guardamos en la variable m
		err1 := binary.Read(buffer1, binary.BigEndian, &m1)
		if err1 != nil {

			log.Fatal("binary.Read failed", err1)
		}
		return m1, true
	}

	return STRUCTURES.DATABLOCK{}, false
}

func GetStringInBlock(name [25]byte) string {
	var s string = ""
	for _, v := range name {
		if v != 0 {
			s = s + string(v)
		}
	}
	return s
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes

	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

//////////CRUD

func CreateABlock(sb STRUCTURES.SUPERBOOT, texto string, file *os.File) int64 {

	//ITERO EN EL BITMAP DE BLOQUES PARA SABER CUANTOS HAY Y EN QUE POSICION DEBO CREAR EL NUEVO BLOQUE
	file.Seek(sb.SB_ap_bitmap_blocks, 0)
	b1 := make([]byte, (sb.SB_ap_blocks - sb.SB_ap_bitmap_blocks))
	n1, err := file.Read(b1)
	if err != nil {

	}
	var contador int64 = 0
	for i := 0; i < len(string(b1[:n1])); i++ {
		if string(b1[:n1][i]) == "1" {
			contador = contador + 1
		}
	}
	//NOS SITUAMOS EN EL BITMAP DE BLOQUE correspondiente Y ESCRIBIMOS UNO
	file.Seek(sb.SB_ap_bitmap_blocks+contador, 0)
	var unit int8 = '1'
	s1 := &unit
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

	//NOS SITUAMOS EN EL BIT CORRESPONDIENTE PARA CREAR EL NUEVO BLOQUE
	file.Seek(sb.SB_ap_blocks+(contador*int64(unsafe.Sizeof(STRUCTURES.DATABLOCK{}))), 0)

	var block = STRUCTURES.DATABLOCK{}

	copy(block.DB_data[:], texto)
	//MANDO A ESCRIBIR EL BLOQUE
	block11 := &block
	var binario7 bytes.Buffer
	binary.Write(&binario7, binary.BigEndian, block11)
	escribirBytes(file, binario7.Bytes())

	return sb.SB_ap_blocks + (contador * int64(unsafe.Sizeof(STRUCTURES.DATABLOCK{}))) //RETORNAMOS EL BIT CORRESPONDIENTE AL NUEVO BLOQUE
}

func UpdateBlock(inicio int64, texto string, file *os.File) {
	file.Seek(inicio, 0)
	var block = STRUCTURES.DATABLOCK{}

	copy(block.DB_data[:], texto)
	//MANDO A ESCRIBIR EL BLOQUE
	block11 := &block
	var binario7 bytes.Buffer
	binary.Write(&binario7, binary.BigEndian, block11)
	escribirBytes(file, binario7.Bytes())
}
