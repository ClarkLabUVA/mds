package bolognese

import (
	"testing"
)

func TestDataciteConversion(t *testing.T) {
	var err error

	// create a temporary file from ioutil
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

	dataciteXML, err := bologneseConvertXML(content)

	if err != nil {
		t.Fatalf("bologneseConverXML Error: %s", err.Error())
	}

	if len(dataciteXML) == 0 {
		t.Fatalf("bologneseConvertXML Error: conversion output is null")
	}

	t.Logf("bologneseConvertXML: Success\n\tOutput: %s", string(dataciteXML))
}
