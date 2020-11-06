//© 2020 By The Rector And Visitors Of The University Of Virginia

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package bolognese

import (
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
)

func init() {
	_, err := exec.LookPath("bolognese")
	if err != nil {
		log.Fatal("Bologonese not installed")
	}
	//log.Printf("INIT BolognesePath: %s\n", path)
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
	//log.Printf("bologneseConvertXML Output: %q", out.String())

	// obtain the xml from the string
	convertedMetadata = out.Bytes()
	return
}
