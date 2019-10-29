package main

import (
	"net/http"
	//"encoding/json"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
)

var errFailedPost = errors.New("Failed Stardog Post")
var errTXFailed = errors.New("Transaction Failed")
var jsonLD = "application/ld+json"

type StardogServer struct {
	URI      string
	Password string
	Username string
	Database string
}

func (s *StardogServer) AddIdentifier(payload []byte) (err error) {

	txId, err := s.NewTransaction()
	if err != nil {
		return
	}

	err = s.AddData(txId, payload, "")
	if err != nil {
		return
	}

	err = s.Commit(txId)
	return
}

func (s *StardogServer) RemoveIdentifier(payload []byte) (err error) {

	txId, err := s.NewTransaction()
	if err != nil {
		return
	}

	err = s.RemoveData(txId, payload, "")
	if err != nil {
		return
	}

	err = s.Commit(txId)
	return
}

// POST /{db}/transaction/begin -> text/plain
func (s *StardogServer) NewTransaction() (t string, err error) {

	url := s.URI + "/" + s.Database + "/transaction/begin"

	txId, err := s.postStardog(url, nil)

	t = string(txId)
	return

}

// POST /{db}/{txId}/remove -> void | text/plain
func (s *StardogServer) RemoveData(txId string, data []byte, namedGraphURI string) (err error) {

	url := s.URI + "/" + s.Database + "/" + txId + "/remove"

	if namedGraphURI != "" {
		url = url + "?graph-uri=" + namedGraphURI
	}

	_, err = s.postStardog(url, data)

	return
}

// POST /{db}/{txId}/add → void | text/plain
func (s *StardogServer) AddData(txId string, data []byte, namedGraphURI string) (err error) {

	url := s.URI + "/" + s.Database + "/" + txId + "/add"

	if namedGraphURI != "" {
		url = url + "?graph-uri=" + namedGraphURI
	}

	_, err = s.postStardog(url, data)

	return

}

// POST /{db}/transaction/commit/{txId}/ -> void | text/plain
func (s *StardogServer) Commit(txId string) (err error) {

	url := s.URI + "/" + s.Database + "/transaction/commit/" + txId
	_, err = s.postStardog(url, nil)
	return

}

func (s *StardogServer) postStardog(url string, data []byte) (responseBody []byte, err error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	req.SetBasicAuth(s.Username, s.Password)
	req.Header.Add("Content-Type", jsonLD)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode == 200 {
		responseBody, _ = ioutil.ReadAll(resp.Body)
		return
	}

	responseBody, _ = ioutil.ReadAll(resp.Body)
	err = fmt.Errorf("%w: %s\tStatusCode: %d\tResponse:%s", errFailedPost, url, resp.StatusCode, responseBody)
	return
}
