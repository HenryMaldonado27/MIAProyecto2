package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func returnMBR(path string) MBRU {
	dataIn := ReadFromFile(path)
	mbr := DecodeToPerson(dataIn)

	return mbr
}

func WriteToFile(s []byte, file string, pos int64) {

	f, err := os.OpenFile(file, os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Seek(pos, io.SeekStart)
	if err != nil {
		panic(err)
	}

	f.Write(s)
}

func ReadFromFile(path string) []byte {

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 720)
	n, err := f.Read(buf[:cap(buf)])
	buf = buf[:n]
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	return buf
}

func ReadFromFileEBR(path string, ofset int64) []byte {

	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}

	return b[721:900]
}

func EncodeToBytes(p interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func DecodeToPerson(s []byte) MBRU {

	mbr := MBRU{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&mbr)
	if err != nil {
		log.Fatal(err)
	}

	return mbr
}

func DecodeToEBR(s []byte) EBR {

	ebr := EBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&ebr)
	if err != nil {
		log.Fatal(err)
	}

	return ebr
}
