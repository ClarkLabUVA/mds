package main

import (
    "errors"
    "fmt"
)

var (
	errRequestAquire    = errors.New("net/http failed to aquire request")
	errRequestExecute = errors.New("net/http failed in executing the request")
	errRequestReadBody  = errors.New("Failed to read body")
)

type APIError struct {
    Message string
    TargetURL string
    Method string
    ResponseStatusCode int
    ResponseBody []byte
}

func (err APIError) Error() string {

    return fmt.Sprintf(`{"message": "%s", "requestMethod": %s, "requestTarget": "%s", "responseStatusCode": %d,  "responseBody": "%s"}`, err.Message, err.Method, err.TargetURL, err.ResponseStatusCode, err.ResponseBody)

}
