package main

import (
	"strconv"
	"strings"
)

func crearParticion(size string, unit string, path string, tipe string, fit string, name string) string {
	mbru := returnMBR(path)

	tamanoReal, err := strconv.Atoi(size)
	if err == nil {
		if strings.ToLower(unit) == "k" {
			tamanoReal = tamanoReal * 1024
		} else if strings.ToLower(unit) == "m" {
			tamanoReal = tamanoReal * 1024 * 1024
		}
	}

	if tipe == "p" || tipe == "e" {
		if mbru.Mbr_partition_1.Part_status == 0 || mbru.Mbr_partition_2.Part_status == 0 ||
			mbru.Mbr_partition_3.Part_status == 0 || mbru.Mbr_partition_4.Part_status == 0 {
			parR := crearParticionPE(mbru, tamanoReal)
			parR.Part_size = tamanoReal
			parR.Part_name = name
			parR.Part_fit = "ff" //fit
			parR.Part_status = 1
			parR.Part_type = tipe

			if parR.Part_start > 0 {

				if mbru.Mbr_partition_1.Part_status == 0 {
					mbru.Mbr_partition_1 = parR
				} else if mbru.Mbr_partition_2.Part_status == 0 {
					mbru.Mbr_partition_2 = parR
				} else if mbru.Mbr_partition_3.Part_status == 0 {
					mbru.Mbr_partition_3 = parR
				} else {
					mbru.Mbr_partition_4 = parR
				}

				dataOut := EncodeToBytes(mbru)
				WriteToFile(dataOut, path, 0)

				if parR.Part_type == "e" {
					var ebr EBR
					ebr.Part_fit = "ff"
					ebr.Part_name = ""
					ebr.Part_next = 0
					ebr.Part_size = 0
					ebr.Part_start = parR.Part_start
					ebr.Part_status = 0

					guardarEBR(ebr, path, int64(parR.Part_start))
				}

				return "Se creo la particion correctamente"

			} else {
				return "No se creo la particion"
			}

		} else {
			return "No se pueden crear mas particiones"
		}
	} else if tipe == "l" {
		creadoL := false

		if mbru.Mbr_partition_1.Part_type == "e" && mbru.Mbr_partition_1.Part_status == 1 {
			creadoL = crearParticionL(mbru, tamanoReal,
				path, fit, name, mbru.Mbr_partition_1)
		} else if mbru.Mbr_partition_2.Part_type == "e" && mbru.Mbr_partition_2.Part_status == 1 {
			creadoL = crearParticionL(mbru, tamanoReal,
				path, fit, name, mbru.Mbr_partition_2)
		} else if mbru.Mbr_partition_3.Part_type == "e" && mbru.Mbr_partition_3.Part_status == 1 {
			creadoL = crearParticionL(mbru, tamanoReal,
				path, fit, name, mbru.Mbr_partition_3)
		} else if mbru.Mbr_partition_4.Part_type == "e" && mbru.Mbr_partition_4.Part_status == 1 {
			creadoL = crearParticionL(mbru, tamanoReal,
				path, fit, name, mbru.Mbr_partition_4)
		} else {
			return "No existe particion extendida"
		}

		if creadoL {
			return "Se creo la particion logica correctamente"
		} else {
			return "No se creo la particion logica correctamente"
		}
	}

	return "Ocurrio un error al crear la particion"
}

func crearParticionPE(mbru MBRU, tamanoReal int) PartitionU {
	arrayInicioOriginal := []int{mbru.Mbr_partition_1.Part_start, mbru.Mbr_partition_2.Part_start,
		mbru.Mbr_partition_3.Part_start, mbru.Mbr_partition_4.Part_start}
	arrayParticiones := []int{1, 2, 3, 4}

	arrayOrdenado := insertion_sort(arrayInicioOriginal, arrayParticiones)

	par1 := mbru.Mbr_partition_1
	par2 := mbru.Mbr_partition_2
	par3 := mbru.Mbr_partition_3
	par4 := mbru.Mbr_partition_4

	if arrayOrdenado[0] == 2 {
		par1 = mbru.Mbr_partition_2
	} else if arrayOrdenado[0] == 3 {
		par1 = mbru.Mbr_partition_3
	} else {
		par1 = mbru.Mbr_partition_4
	}

	if arrayOrdenado[1] == 1 {
		par2 = mbru.Mbr_partition_1
	} else if arrayOrdenado[1] == 3 {
		par2 = mbru.Mbr_partition_3
	} else {
		par2 = mbru.Mbr_partition_4
	}

	if arrayOrdenado[2] == 1 {
		par3 = mbru.Mbr_partition_1
	} else if arrayOrdenado[2] == 2 {
		par3 = mbru.Mbr_partition_2
	} else {
		par3 = mbru.Mbr_partition_4
	}

	if arrayOrdenado[3] == 1 {
		par4 = mbru.Mbr_partition_1
	} else if arrayOrdenado[3] == 2 {
		par4 = mbru.Mbr_partition_2
	} else {
		par4 = mbru.Mbr_partition_3
	}

	i1 := mbru.Mbr_tamano
	i2 := mbru.Mbr_tamano
	i3 := mbru.Mbr_tamano
	i4 := mbru.Mbr_tamano

	if par1.Part_start > 0 {
		i1 = par1.Part_start
	}

	if par2.Part_start > 0 {
		i2 = par2.Part_start
	}

	if par3.Part_start > 0 {
		i3 = par3.Part_start
	}

	if par4.Part_start > 0 {
		i4 = par4.Part_start
	}

	e1 := i1 - 720
	e2 := i2 - (i1 + par1.Part_size)
	e3 := i3 - (i2 + par2.Part_size)
	e4 := i4 - (i3 + par3.Part_size)

	var parR PartitionU

	if strings.ToLower(mbru.Dsk_fit) == "ff" {
		if e1 >= tamanoReal {
			parR.Part_start = 721
		} else if e2 >= tamanoReal {
			parR.Part_start = (i1 + par1.Part_size)
		} else if e3 >= tamanoReal {
			parR.Part_start = (i2 + par2.Part_size)
		} else if e4 >= tamanoReal {
			parR.Part_start = (i3 + par3.Part_size)
		}
	} else if strings.ToLower(mbru.Dsk_fit) == "bf" {
		if e1 >= tamanoReal && e2 >= e1 && e3 >= e1 && e4 >= e1 {
			parR.Part_start = 721
		} else if e2 >= tamanoReal && e1 >= e2 && e3 >= e2 && e4 >= e2 {
			parR.Part_start = (i1 + par1.Part_size)
		} else if e3 >= tamanoReal && e2 >= e3 && e1 >= e3 && e4 >= e3 {
			parR.Part_start = (i2 + par2.Part_size)
		} else if e4 >= tamanoReal && e2 >= e4 && e3 >= e4 && e1 >= e4 {
			parR.Part_start = (i3 + par3.Part_size)
		}
	} else {
		if e1 >= tamanoReal && e1 >= e2 && e1 >= e3 && e1 >= e4 {
			parR.Part_start = 721
		} else if e2 >= tamanoReal && e2 >= e1 && e2 >= e3 && e2 >= e4 {
			parR.Part_start = (i1 + par1.Part_size)
		} else if e3 >= tamanoReal && e3 >= e1 && e3 >= e2 && e3 >= e4 {
			parR.Part_start = (i2 + par2.Part_size)
		} else if e4 >= tamanoReal && e4 >= e3 && e4 >= e2 && e4 >= e1 {
			parR.Part_start = (i3 + par3.Part_size)
		}
	}

	return parR
}

func guardarEBR(ebr EBR, path string, posicion int64) {
	dataOut := EncodeToBytes(ebr)
	WriteToFile(dataOut, path, posicion)
}

func crearParticionL(mbru MBRU, tamanoReal int, path string,
	fit string, name string, extendida PartitionU) bool {

	inicioExtendida := extendida.Part_start
	tamanoExtendida := extendida.Part_size

	dataIn := ReadFromFileEBR(path, int64(inicioExtendida))
	ebr := DecodeToEBR(dataIn)

	var ebrGuardar EBR
	ebrGuardar.Part_fit = fit
	ebrGuardar.Part_name = name
	ebrGuardar.Part_status = 1
	ebrGuardar.Part_start = inicioExtendida
	ebrGuardar.Part_next = 0
	ebrGuardar.Part_size = tamanoReal

	//inicioBusqueda := inicioExtendida
	//var ebrActual EBR
	//var ebrAnterior EBR

	encontrado := false
	//ebrActual = ebr
	//ebrAnterior = ebr

	if ebr.Part_next == 0 {
		if tamanoExtendida >= tamanoReal {
			guardarEBR(ebrGuardar, path, int64(inicioExtendida))
			return true
		} else {
			return false
		}
	} else {

		for {
			if extendida.Part_fit == "f" {

			} else if extendida.Part_fit == "w" {

			} else if extendida.Part_fit == "b" {

			}
			//var tamano1 int = 0
			//var tamano2 int = 0

			if encontrado {
				break
			}
		}

	}

	return true
}

func insertion_sort(arreglo []int, particiones []int) []int {
	for i := 1; i < len(arreglo); i++ {
		j := i
		for j > 0 && arreglo[j-1] > arreglo[j] {
			swap(j-1, j, &arreglo, &particiones)
			j--
		}
	}

	return particiones
}

func swap(previo, actual int, puntero_arreglo *[]int, puntero_particiones *[]int) {
	arreglo := *puntero_arreglo
	copia := arreglo[actual]
	arreglo[actual] = arreglo[previo]
	arreglo[previo] = copia

	particiones := *puntero_particiones
	copiaP := particiones[actual]
	particiones[actual] = particiones[previo]
	particiones[previo] = copiaP
}
