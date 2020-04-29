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
