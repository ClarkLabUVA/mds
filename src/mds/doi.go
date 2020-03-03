package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"
)

var DataciteUser = "VIVA.UVA-TEST"
var DatacitePassword = "Lib#14Books"
var DatacitePrefix = "10.70020"
var DataciteBasicAuth = base64.StdEncoding.EncodeToString([]byte("Basic " + base64.StdEncoding.EncodeToString([]byte(DataciteUser+":"+DatacitePassword))))

var (
	errRequestAquire    = errors.New("net/http failed to aquire request")
	errRequestExecution = errors.New("net/http failed in executing the request")
)

type Datacite struct {
	Username string
	Password string
	Prefix   string
	Auth     string
}

type DOI struct {
	Identifier  string
	URL         string
	Content     []byte
	DataciteXML []byte
}

func NewDOI(identifier string, content []byte, url string) (doi DOI, err error) {
	doi.Identifier = identifier
	doi.Content = content
	doi.URL = url

	// convert metadata
	doi.DataciteXML, err = bologneseConvertXML(content)

	return
}

func (doi *DOI) dataciteCreate() (err error) {
	// create metadata
	err = doi.datacitePutMetadata()

	if err != nil {
		return
	}

	// create datacite resolver link
	err = doi.datacitePutResolver()

	if err != nil {
		return
	}

	return
}

func (doi *DOI) dataciteDeleteDOI() (err error) {

	// mark metadata
	doi.dataciteDeleteMetadata()

	return
}

func (doi *DOI) dataciteUpdate() (err error) {
	// PUT https://mds.test.datacite.org/metadata/:doi

	// update metadata
	doi.datacitePutMetadata()

	return
}

func (doi *DOI) datacitePutMetadata() (err error) {
	// PUT https://mds.test.datacite.org/metadata/10.5072/0000-03VC

	url := "https://mds.test.datacite.org/metadata/" + doi.Identifier

	client := &http.Client{}

	bodyBuffer := bytes.NewBuffer(doi.DataciteXML)
	req, err := http.NewRequest("PUT", url, bodyBuffer)

	if err != nil {
		return
	}

	req.Header.Add("Authorization", DataciteBasicAuth)

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	// determine success of request
	if resp.StatusCode == 400 {
		return
	}

	return
}

func (doi *DOI) datacitePutResolver() (err error) {
	// PUT https://mds.test.datacite.org/
	url := "https://mds.test.datacite.org/doi/" + doi.Identifier

	payload := []byte("doi=" + doi.Identifier + "\nurl=" + doi.URL)

	client := http.Client{}

	bodyBuffer := bytes.NewBuffer(payload)
	req, err := http.NewRequest("PUT", url, bodyBuffer)

	if err != nil {
		return
	}

	req.Header.Add("Authorization", DataciteBasicAuth)

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	// log state of request

	// determine success of request
	if resp.StatusCode == 400 {
		return
	}

	return
}

func (doi *DOI) dataciteDeleteMetadata() (err error) {
	// DELETE https://mds.test.datacite.org/metadata/:doi

	return
}
