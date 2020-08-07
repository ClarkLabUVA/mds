package auth

import (
	"net/http"
	"strings"
)

var (
	AuthURI = "http://auth/"
	AuthInspect = AuthURI + "inspect"
)

// AuthMiddleware is a handler for the Fairscape auth service
// it checks that token is present and valid for a user
// implemented as negroni middleware
func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// read bearer token from request
	authHeader := strings.TrimPrefix(
		r.Header.Get("Authorization"),
		"Bearer",
		)


	// if bearer token doesn't exist
	if authHeader == "" {
		w.Write([]byte(`{"error": "missing Authorization bearer token"`))
		w.WriteHeader(400)
		return
	}

	// call authorization service

	client := &http.Client{}

	req, err := http.NewRequest("POST", AuthInspect, nil)


	if err != nil {
		w.Write([]byte(`{"error": "error creating http request"`))
		w.WriteHeader(500)
		return
	}

	res, err := client.Do(req)

	// if there is an error in preforming the service call
	if err != nil {
		w.Write([]byte(`{"error": "error creating http request"`))
		w.WriteHeader(500)
		return
	}

	if res.StatusCode == 204 {
		// Call the next handler 
		next(w, r)
	} else {
		w.Write([]byte(`{"error": "user not authorized"}`))
		w.WriteHeader(401)
		return
	}


}
