package main

import (
	"testing"
)


func TestDOIConstructor(t *testing.T) {

    // test DOI constructor
    t.Run("Success", func(tester *testing.T) {

	content := []byte(
`{  
    "@context":"http://schema.org",
    "@type":"SoftwareSourceCode",
    "@id": "https://doi.org/10.5438/qeg0-3gm3",
    "url":"https://github.com/datacite/maremma",
    "name":"Maremma: a Ruby library for simplified network calls",
    "author":{  
    "@type":"person",
    "@id":"http://orcid.org/0000-0003-0077-4738",
    "name":"Martin Fenner"
    },
    "description":"Simplifies network calls, including json/xml parsing and error handling. Based on Faraday.",
    "keywords":"faraday, excon, net/http",
    "dateCreated":"2015-11-28",
    "datePublished":"2017-02-24",
    "dateModified":"2017-02-24",
    "publisher":{  
      "@type":"Organization",
      "name":"DataCite"
    }
}`)

	url = "https://example.org"
	identifier = "10.5072/1234"
	doi, err := NewDOI(identifier, content, url)

	if err != nil {
	    tester.Fatalf("Constructor Failed to Produce DOI\n\tError: %s", err.Error())
	}
    })

    t.Run("NoMetadata", func(testing *testing.T) {

    })

}

func TestDOIDatacitePutMetadata(t *testing.T) {

}

func TestDOIDatacitePutResolver(t *testing.T) {

}

func TestDOIUpdate(t *testing.T) {

}

func TestDOICreate(t *testing.T) {

}

func TestDOIDelete(t *testing.T) {

}


