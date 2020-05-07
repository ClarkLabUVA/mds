package identifier

import (
	"net/http"
	"log"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"

	//"encoding/json"
	"github.com/rs/zerolog"
	"os"
)

var stardogLogger = zerolog.New(os.Stderr).With().Timestamp().Str("backend", "stardog").Logger()

var (
	errFailedPost = errors.New("Failed Stardog Post")
	errTXFailed = errors.New("Transaction Failed")
)

var jsonLD = "application/ld+json"

type StardogServer struct {
	URI      string
	Password string
	Username string
	Database string
	ValidationURI	string
}


// TODO: Fix Request Can't Find
func (s *StardogServer) CreateDatabase(databaseName string) (responseBody []byte, statusCode int, err error) {

	url := s.URI + "/admin/databases"

	data := []byte(`{"dbname": "`+ databaseName +`", "options": {}, "files": []}`)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// create the part root
	mime := make(textproto.MIMEHeader)

	mime.Add("content-type", "application/json")
	mime.Add("content-disposition", `form-data; name="root"`)


	root, _ := writer.CreatePart(mime)
	//root, _ := writer.CreateFormField("root")
	root.Write(data)
	writer.Close()


	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "createDatabase").
			Str("url", url).
			Msg("failed to acquire http request")

		return
	}

	req.SetBasicAuth(s.Username, s.Password)

	// Set Headers
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "multipart/form-data; boundary=" + writer.Boundary() )

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "createDatabase").
			Str("url", url).
			Msg("failed to preform request")

		return
	}

	responseBody, _ = ioutil.ReadAll(resp.Body)

	stardogLogger.Info().
		Str("operation", "createDatabase").
		Str("url", url).
		Int("statusCode", resp.StatusCode).
		Bytes("response", responseBody).
		Msg("preformed create database")

	return

}

func (s *StardogServer) DropDatabase(databaseName string) (response []byte, err error) {

	url := s.URI + "/admin/databases/" + databaseName


	req, err := http.NewRequest("POST", url)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "dropDatabase")
			Msg("failed to acquire http request")
		return
	}

	req.SetBasicAuth(s.Username, s.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "dropDatabase").
			Msg("failed to preform request")

		return
	}

	response, _ = ioutil.ReadAll(resp.Body)

	stardogLogger.Info().
		Str("operation", "dropDatabase").
		Str("url", url).
		Int("statusCode", resp.StatusCode).
		Bytes("response", responseBody).
		Msg("preformed create database")

	return

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

	txId, err := s.request(url, "POST", nil)
	t = string(txId)
	return

}

// POST /{db}/{txId}/remove -> void | text/plain
func (s *StardogServer) RemoveData(txId string, data []byte, namedGraphURI string) (err error) {

	url := s.URI + "/" + s.Database + "/" + txId + "/remove"

	if namedGraphURI != "" {
		url = url + "?graph-uri=" + namedGraphURI
	}

	_, err = s.request(url, "POST", data)

	return
}

// POST /{db}/{txId}/add → void | text/plain
func (s *StardogServer) AddData(txId string, data []byte, namedGraphURI string) (err error) {

	url := s.URI + "/" + s.Database + "/" + txId + "/add"

	if namedGraphURI != "" {
		url = url + "?graph-uri=" + namedGraphURI
	}

	_, err = s.request(url, "POST", data)

	return

}

// POST /{db}/transaction/commit/{txId}/ -> void | text/plain
func (s *StardogServer) Commit(txId string) (err error) {

	url := s.URI + "/" + s.Database + "/transaction/commit/" + txId
	_, err = s.request(url, "POST", nil)
	return

}

func (s *StardogServer) request(url string, method string, data []byte) (responseBody []byte, err error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
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
