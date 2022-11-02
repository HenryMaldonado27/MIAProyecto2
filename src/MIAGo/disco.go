package main

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func createDisk(path string, unit string, size string, fit string) string {

	pathDir := strings.Split(path, "/")
	directorio := ""
	for i := 0; i < len(pathDir)-1; i++ {
		directorio = directorio + "/" + pathDir[i]
	}
	err := os.Mkdir(directorio, os.ModePerm)

	tamanoReal, err := strconv.Atoi(size)
	if strings.ToLower(unit) == "k" {
		tamanoReal = tamanoReal * 1024
	} else {
		tamanoReal = tamanoReal * 1024 * 1024
	}

	t := time.Now()

	file, err := os.Create(path)
	if err != nil {
	}

	var part PartitionU
	part.Part_fit = "ff"
	part.Part_name = "a"
	part.Part_size = 0
	part.Part_start = 0
	part.Part_status = 0
	part.Part_type = "p"

	var mbru MBRU
	mbru.Dsk_fit = fit
	mbru.Mbr_dsk_signature = rand.Intn(100)
	mbru.Mbr_fecha_creacion = t.String()[0:16]
	mbru.Mbr_tamano = tamanoReal
	mbru.Mbr_partition_1 = part
	mbru.Mbr_partition_2 = part
	mbru.Mbr_partition_3 = part
	mbru.Mbr_partition_4 = part

	file.Close()

	dataOut := EncodeToBytes(mbru)
	WriteToFile(dataOut, path, 0)

	for i := 0; i < tamanoReal-720; i++ {
		pos := int64(720 + i)
		WriteToFile([]byte{'0'}, path, pos)
	}

	if err != nil {
		return "Fallo la creacion del disco"
	} else {
		return "Se creo el disco correctamente"
	}

}

func deleteDisk(path string) string {
	e := os.Remove(path)
	if e != nil {
		return e.Error()
	}

	return "Se elimino el disco correctamente"
}
