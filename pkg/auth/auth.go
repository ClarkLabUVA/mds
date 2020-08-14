//© 2020 By The Rector And Visitors Of The University Of Virginia

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package auth

import (
	"net/http"
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
    var authHeader string

	authHeader = r.Header.Get("Authorization")

	// if bearer token doesn't exist
	if authHeader == "" {

        // check cookies of reqest
        authCookie, err := r.Cookie("fairscapeAuth") 

        if err != nil {
            w.Write([]byte(`{"error": "request missing authorization token"}`))
            w.WriteHeader(400)
            return
        }

        authHeader = authCookie.Value
        
	}

	// call authorization service

	client := &http.Client{}

	req, err := http.NewRequest("POST", AuthInspect, nil)

	req.Header.Set("Authorization",  authHeader)

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
