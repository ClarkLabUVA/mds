//© 2020 By The Rector And Visitors Of The University Of Virginia

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package identifier

import (
	"bytes"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"

	//"encoding/json"
	"github.com/rs/zerolog"
	"os"
)

var stardogLogger = zerolog.New(os.Stderr).With().Timestamp().Str("backend", "stardog").Logger()

var (
	errFailedPost = errors.New("Failed Stardog Post")
	errTXFailed   = errors.New("Transaction Failed")
    errStardogPingFail = errors.New("Stardog Ping failed to return status 200")
)

var jsonLD = "application/ld+json"

type StardogServer struct {
	URI           string
	Password      string
	Username      string
	Database      string
	ValidationURI string
}

// Ping simply checks the health of the stardog server
func (s *StardogServer) Ping() (err error) {
    url := s.URI + "/admin/healthcheck"
    
    req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "Ping").
			Str("url", url).
			Msg("failed to preform request")

		return
	}

    if resp.StatusCode != 200 {
        err = errStardogPingFail
    }

    return


}

// TODO: Fix Request Can't Find
func (s *StardogServer) CreateDatabase(databaseName string) (responseBody []byte, statusCode int, err error) {

	url := s.URI + "/admin/databases"

	data := []byte(`{"dbname": "` + databaseName + `", "options": {}, "files": []}`)

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
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())

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

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "dropDatabase").
			Str("url", url).
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
			Str("url", url).
			Msg("failed to preform request")

		return
	}

	response, _ = ioutil.ReadAll(resp.Body)

	stardogLogger.Info().
		Str("operation", "dropDatabase").
		Str("url", url).
		Int("statusCode", resp.StatusCode).
		Bytes("response", response).
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

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "newTransaction").
			Str("url", url).
			Msg("failed to acquire http request")
		return
	}

	req.SetBasicAuth(s.Username, s.Password)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "newTransaction").
			Str("url", url).
			Msg("failed to preform request")

		return
	}

	responseBody, _ := ioutil.ReadAll(response.Body)
	t = string(responseBody)

	stardogLogger.Info().
		Str("operation", "newTransaction").
		Str("url", url).
		Int("statusCode", response.StatusCode).
		Str("transaction", t).
		Msg("created transaction")

	return

}

// POST /{db}/{txId}/remove -> void | text/plain
func (s *StardogServer) RemoveData(txId string, data []byte, namedGraphURI string) (err error) {

	url := s.URI + "/" + s.Database + "/" + txId + "/remove"

	if namedGraphURI != "" {
		url = url + "?graph-uri=" + namedGraphURI
	}

	body := &bytes.Buffer{}
	body.Write(data)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "removeData").
			Str("transaction", txId).
			Str("url", url).
			Str("data", string(data)).
			Msg("failed to acquire http request")
		return
	}

	req.SetBasicAuth(s.Username, s.Password)
	req.Header.Add("Content-Type", "application/ld+json")

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "removeData").
			Str("transaction", txId).
			Str("url", url).
			Str("data", string(data)).
			Msg("failed to preform request")

		return
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	stardogLogger.Info().
		Str("operation", "removeData").
		Str("transaction", txId).
		Str("url", url).
		Str("data", string(data)).
		Int("statusCode", response.StatusCode).
		Str("responseBody", string(responseBody)).
		Msg("created transaction")

	return
}

// POST /{db}/{txId}/add → void | text/plain
func (s *StardogServer) AddData(txId string, data []byte, namedGraphURI string) (err error) {

	url := s.URI + "/" + s.Database + "/" + txId + "/add"

	if namedGraphURI != "" {
		url = url + "?graph-uri=" + namedGraphURI
	}

	body := &bytes.Buffer{}
	body.Write(data)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "addData").
			Str("transaction", txId).
			Str("url", url).
			Str("data", string(data)).
			Msg("failed to acquire http request")
		return
	}

	req.SetBasicAuth(s.Username, s.Password)
	req.Header.Add("Content-Type", "application/ld+json")

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "addData").
			Str("transaction", txId).
			Str("url", url).
			Str("data", string(data)).
			Msg("failed to preform request")

		return
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	stardogLogger.Info().
		Str("operation", "addData").
		Str("transaction", txId).
		Str("url", url).
		Str("data", string(data)).
		Int("statusCode", response.StatusCode).
		Str("responseBody", string(responseBody)).
		Msg("created transaction")

	return

	return

}

// POST /{db}/transaction/commit/{txId}/ -> void | text/plain
func (s *StardogServer) Commit(txId string) (err error) {

	url := s.URI + "/" + s.Database + "/transaction/commit/" + txId

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "commitTransaction").
			Str("transaction", txId).
			Str("url", url).
			Msg("failed to acquire http request")
		return
	}

	req.SetBasicAuth(s.Username, s.Password)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		stardogLogger.Error().
			Err(err).
			Str("operation", "commitTransaction").
			Str("transaction", txId).
			Str("url", url).
			Msg("failed to preform request")

		return
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	stardogLogger.Info().
		Str("operation", "commitTransaction").
		Str("transaction", txId).
		Str("url", url).
		Int("statusCode", response.StatusCode).
		Str("responseBody", string(responseBody)).
		Msg("created transaction")

	return

}
