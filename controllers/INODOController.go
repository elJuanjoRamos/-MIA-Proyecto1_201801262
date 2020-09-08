package controllers

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"unsafe"

	STRUCTURES "../structures"
)

func INODOCONTROLLER_CreateINODO(superboot STRUCTURES.SUPERBOOT, file *os.File, lenArchivo int64, idProper int64) int64 {
	//NOS SITUAMOS AL INICIO DEL BITMAP DE INODOS
	file.Seek(superboot.SB_ap_bitmap_table_inode, 0)

	//ITERO EN EL BITMAP DE BLOQUES PARA SABER CUANTOS HAY Y EN QUE POSICION DEBO CREAR EL NUEVO BLOQUE
	b1 := make([]byte, (superboot.SB_ap_table_inode - superboot.SB_ap_bitmap_table_inode))
	n1, err := file.Read(b1)
	if err != nil {

	}
	var contador int64 = 0
	for i := 0; i < len(string(b1[:n1])); i++ {
		if string(b1[:n1][i]) == "1" {
			contador = contador + 1
		}
	}
	//NOS SITUAMOS EN EL BITMAP DE INODOS correspondiente Y ESCRIBIMOS UNO
	file.Seek(superboot.SB_ap_bitmap_table_inode+contador, 0)
	var unit int8 = '1'
	s1 := &unit
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

	//NOS SITUAMOS EN EL BIT CORRESPONDIENTE PARA CREAR EL NUEVO BLOQUE
	file.Seek(superboot.SB_ap_table_inode+(contador*int64(unsafe.Sizeof(STRUCTURES.TABLEINODE{}))), 0)

	var inodo = STRUCTURES.TABLEINODE{
		I_count_inodo:             contador + 1,
		I_size_archivo:            lenArchivo,
		I_count_bloques_asignados: 0,
		I_ap_indirecto:            -1,
		I_id_proper:               idProper,
	}

	block11 := &inodo
	var binario7 bytes.Buffer
	binary.Write(&binario7, binary.BigEndian, block11)
	escribirBytes(file, binario7.Bytes())

	return superboot.SB_ap_table_inode + (contador * int64(unsafe.Sizeof(STRUCTURES.TABLEINODE{})))
}

func GETInode(inodoInicio int64, file *os.File) STRUCTURES.TABLEINODE {

	//NOS SITUAMOS AL INICIO DEL INODO
	file.Seek(inodoInicio, 0)

	//LEEMOS LA ESTRUCTURA TABLAINODO
	m := STRUCTURES.TABLEINODE{}
	//Obtenemos el tamanio del TABLE INODO
	var size int = int(unsafe.Sizeof(m))
	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	return m
}
