package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var respuesta *http.Request
var writer http.ResponseWriter

func main() {

	fmt.Printf("Starting server at port 8080\n")

	http.HandleFunc("/leerExe", readExe)
	http.HandleFunc("/leerComando", readComando)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func convertResponseJSON(respuesta string) string {
	responseJSON := &RespuestaJSON{Respuesta: respuesta}
	b, err := json.Marshal(responseJSON)
	if err != nil {
		return string(err.Error())
	}
	return string(b)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8100")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func readExe(rw http.ResponseWriter, req *http.Request) {
	enableCors(&rw)
	writer = rw
	respuesta = req

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	data := CmdStruct{}
	json.Unmarshal([]byte(string(body)), &data)

	leerExe(data)
}

func leerExe(data CmdStruct) {
	b, err := ioutil.ReadFile(data.Command)

	if err != nil {
		fmt.Print(err)
	}

	str := string(b)

	fmt.Fprintf(writer, convertResponseJSON(str))
}

func readComando(rw http.ResponseWriter, req *http.Request) {
	enableCors(&rw)
	writer = rw
	respuesta = req

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	data := CmdStruct{}
	json.Unmarshal([]byte(string(body)), &data)

	elegirComando(data.Command)
}

func elegirComando(comando string) {
	comandoStr := strings.Split(comando, " -")
	comandoStrLower := strings.TrimSpace(strings.ToLower(comandoStr[0]))

	if comandoStrLower == "mkdisk" {
		parametrosMKDISK(comando)
	} else if comandoStrLower == "rmdisk" {
		parametrosRMDISK(comando)
	} else if comandoStrLower == "fdisk" {
		parametrosFDISK(comando)
	}

}

func parametrosMKDISK(comandoMKDISK string) {
	comandoStr := strings.Split(comandoMKDISK, " -")
	size := ""
	fit := "ff"
	unit := "m"
	path := ""

	for i := 0; i < len(comandoStr); i++ {
		valorEncontrado := strings.Split(comandoStr[i], "->")

		if valorEncontrado[0] == "path" {
			path = strings.Replace(valorEncontrado[1], "\"", "", 2)
		} else if valorEncontrado[0] == "u" {
			unit = valorEncontrado[1]
		} else if valorEncontrado[0] == "s" {
			size = valorEncontrado[1]
		} else if valorEncontrado[0] == "f" {
			fit = valorEncontrado[1]
		}
	}

	if len(size) > 0 {
		if len(path) > 0 {
			crear := createDisk(path, unit, size, fit)
			fmt.Fprintf(writer, convertResponseJSON(crear))
		} else {
			fmt.Fprintf(writer, convertResponseJSON("No se encontro la variable path"))
		}
	} else {
		fmt.Fprintf(writer, convertResponseJSON("No se encontro la variable size"))
	}
}

func parametrosRMDISK(comandoRMDISK string) {
	comandoStr := strings.Split(comandoRMDISK, " -")
	path := ""

	for i := 0; i < len(comandoStr); i++ {
		valorEncontrado := strings.Split(comandoStr[i], "->")

		if valorEncontrado[0] == "path" {
			path = strings.Replace(valorEncontrado[1], "\"", "", 2)
		}
	}

	if len(path) > 0 {
		fmt.Fprintf(writer, convertResponseJSON(deleteDisk(path)))
	} else {
		fmt.Fprintf(writer, convertResponseJSON("No se encontro la variable path"))
	}
}

func parametrosFDISK(comandoFDISK string) {
	comandoStr := strings.Split(comandoFDISK, " -")
	size := ""
	unit := "k"
	path := ""
	tipe := "p"
	fit := "wf"
	name := ""

	fmt.Println(comandoFDISK)
	fmt.Println(comandoStr)

	for i := 0; i < len(comandoStr); i++ {
		valorEncontrado := strings.Split(comandoStr[i], "->")

		if valorEncontrado[0] == "size" {
			size = valorEncontrado[1]
		} else if valorEncontrado[0] == "unit" {
			unit = valorEncontrado[1]
		} else if valorEncontrado[0] == "path" {
			path = strings.Replace(valorEncontrado[1], "\"", "", 2)
		} else if valorEncontrado[0] == "type" {
			tipe = valorEncontrado[1]
		} else if valorEncontrado[0] == "fit" {
			fit = valorEncontrado[1]
		} else if valorEncontrado[0] == "name" {
			name = valorEncontrado[1]
		}
	}

	fmt.Println("termino")

	if len(path) > 0 {
		if len(name) > 0 {
			fmt.Fprintf(writer,
				convertResponseJSON(crearParticion(size, unit, path, tipe, fit, name)))
		} else {
			fmt.Fprintf(writer, convertResponseJSON("No se encontro la variable name"))
		}
	} else {
		fmt.Fprintf(writer, convertResponseJSON("No se encontro la variable path"))
	}
}
