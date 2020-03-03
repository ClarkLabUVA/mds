package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
)

func init() {
	path, err := exec.LookPath("bolognese")
	if err != nil {
		log.Fatal("Bologonese not installed")
	}
	log.Printf("INIT BolognesePath: %s\n", path)
}

func bologneseConvertXML(inputMetadata []byte) (convertedMetadata []byte, err error) {
	// randomly generate a filename
	filename := "identifier" + string(rand.Intn(1000)) + ".json"

	err = ioutil.WriteFile(filename, inputMetadata, 0644)
	if err != nil {
		log.Printf("ERROR bologneseConvertXML: Failed to Write Input to Temp File\n\tError: %s", err.Error())
		return
	}
	defer os.Remove(filename) // clean up

	cmd := exec.Command("bolognese", filename, "-t", "datacite")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Printf("bologneseConvertXML: Failed to Run the Conversion Command\n\tError: %s", err.Error())
		return
	}
	log.Printf("bologneseConvertXML Output: %q", out.String())

	// obtain the xml from the string
	convertedMetadata = out.Bytes()
	return
}
