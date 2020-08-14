//© 2020 By The Rector And Visitors Of The University Of Virginia

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package identifier

import (
	"github.com/buger/jsonparser"
	"testing"
)

func TestMetadata(t *testing.T) {

	data := []byte(`{
	"@id": "10.70200/1234",
	"@context": "https://schema.org/",
	"@type": "Dataset",
	"name": "Test ID",
	"url": "https://example.org",
	"author": "Max Levinson",
	"creativeWorkStatus": "Draft",
	"expires": "2011-07-14T19:43:37+0100",
	"keywords": ["testing", "tdd"]
    }`)

	// get the @id value
	val, dataType, offset, err := jsonparser.Get(data, "@id")

	if err != nil {
		t.Fatalf("Failed to Retrieve @id\n\tError: %s", err.Error())
	}
	t.Logf("Value: %s\tdataType: %d\toffset: %d", val, dataType, offset)

	// set the @id value to a full url
	data, err = jsonparser.Set(data, []byte("https://mds.clark-lab.org/10.70200/1234"), "@id")

	if err != nil {
		t.Fatalf("Failed to Set @id\n\tError: %s", err.Error())
	}
	t.Logf("Metadata: %s", string(data))

}

type Identifier struct {
	GUID         string
	Namespace    string
	URL          string
	Name         string
	version      int
	author       interface{}
	dateCreated  interface{}
	dateModified interface{}
	metadata     interface{}
}

// func NewIdentifier() (Identifier, error) {}
