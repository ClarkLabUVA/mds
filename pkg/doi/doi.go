package doi

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var DataciteUser = "VIVA.UVA-TEST"
var DatacitePassword = "Lib#14Books"
var DatacitePrefix = "10.70020"
var DataciteBasicAuth = "Basic " + base64.StdEncoding.EncodeToString([]byte(DataciteUser+":"+DatacitePassword))

type DataciteCredentials struct {
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
		return fmt.Errorf("Error: %w  Value: %s", errRequestAquire, err.Error())
	}

	req.Header.Add("Authorization", DataciteBasicAuth)

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Error: %w Value: %s", errRequestExecute, err.Error())
	}

	statusCode := resp.StatusCode
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error: %w Value: %s", errRequestReadBody, err.Error())
	}

	// log status code
	// log.Printf("Status Code: %d", statusCode)
	// log.Printf("Response Body: %s", responseBody)

	// determine success of request
	if statusCode == 201 {
		return nil
	}

	apiErr := APIError{
		TargetURL:          url,
		Method:             "PUT",
		ResponseStatusCode: statusCode,
		ResponseBody:       responseBody,
	}

	if statusCode == 422 {
		apiErr.Message = "DOI Missing Required Metadata"
		return apiErr
	}

	return apiErr
}

func (doi *DOI) datacitePutResolver() (err error) {
	// PUT https://mds.test.datacite.org/
	url := "https://mds.test.datacite.org/doi/" + doi.Identifier

	payload := []byte("doi=" + doi.Identifier + "\nurl=" + doi.URL)

	client := http.Client{}

	bodyBuffer := bytes.NewBuffer(payload)

	req, err := http.NewRequest("PUT", url, bodyBuffer)
	if err != nil {
		return fmt.Errorf("Error: %w  Value: %s", errRequestAquire, err.Error())
	}

	req.Header.Add("Authorization", DataciteBasicAuth)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error: %w Value: %s", errRequestExecute, err.Error())
	}

	statusCode := resp.StatusCode
	responseBody, _ := ioutil.ReadAll(resp.Body)

	// Log Response
	// log.Printf("PutResolver StatusCode: %d ResponseBody: %s", statusCode, responseBody)

	// determine success of request
	if resp.StatusCode == 201 {
		return nil
	}

	apiErr := APIError{
		TargetURL:          url,
		Method:             "PUT",
		ResponseStatusCode: statusCode,
		ResponseBody:       responseBody,
	}

	return apiErr
}

func (doi *DOI) dataciteDeleteMetadata() (err error) {
	// DELETE https://mds.test.datacite.org/metadata/:doi
	url := "https://mds.test.datacite.org/metadata/" + doi.Identifier
	client := http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Error: %w  Value: %s", errRequestAquire, err.Error())
	}

	req.Header.Add("Authorization", DataciteBasicAuth)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error: %w Value: %s", errRequestExecute, err.Error())
	}

	statusCode := resp.StatusCode
	responseBody, _ := ioutil.ReadAll(resp.Body)

	// Log Response
	log.Printf("DeleteMetadata StatusCode: %d ResponseBody: %s", statusCode, responseBody)

	// determine success of request
	if statusCode == 200 {
		return nil
	}

	apiErr := APIError{
		TargetURL:          url,
		Method:             "DELETE",
		ResponseStatusCode: statusCode,
		ResponseBody:       responseBody,
	}

	return apiErr
}
