//© 2020 By The Rector And Visitors Of The University Of Virginia

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"net/http"
	"os"

	"github.com/ClarkLabUVA/mds/pkg/identifier"
	"github.com/gorilla/mux"

	"log"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)


var server identifier.Backend

func main() {

	// set server to defaults for local testing
	server = identifier.Backend{
		Stardog: identifier.StardogServer{
			URI:      "http://localhost:5820",
			Password: "admin",
			Username: "admin",
			Database: "ors",
		},
		Mongo: identifier.MongoServer{
			URI:      "mongodb://mongoadmin:mongosecret@localhost:27017",
			Database: "ors",
			Collection: "ids",
		},

	}

	// if Environent Variables options are set, update backend server configuration
	if mongoURI, exists := os.LookupEnv("MONGO_URI"); exists {
		server.Mongo.URI = mongoURI
	}

	if mongoDB, exists := os.LookupEnv("MONGO_DB"); exists {
		server.Mongo.Database = mongoDB
	}

	if mongoCol, exists := os.LookupEnv("MONGO_COL"); exists {
		server.Mongo.Collection = mongoCol
	}

	if stardogURI, exists := os.LookupEnv("STARDOG_URI"); exists {
		server.Stardog.URI = stardogURI
	}

	if stardogDB, exists := os.LookupEnv("STARDOG_DATABASE"); exists {
		server.Stardog.Database = stardogDB
	}

	if stardogPassword, exists := os.LookupEnv("STARDOG_PASSWORD"); exists {
		server.Stardog.Password = stardogPassword
	}

	if stardogUsername, exists := os.LookupEnv("STARDOG_USERNAME"); exists {
		server.Stardog.Username = stardogUsername
	}


	// Log Initilization Variables
	zlog.Info().
		Dict("stardog", zerolog.Dict().
			Str("uri", server.Stardog.URI).
			Str("username", server.Stardog.Username).
			Str("password", server.Stardog.Password).
			Str("database", server.Stardog.Database),
		).
		Dict("mongo", zerolog.Dict().
			Str("uri", server.Mongo.URI).
			Str("database", server.Mongo.Database).
			Str("collection", server.Mongo.Collection),
		).
		Msg("initilization variables for server")

	server.Stardog.CreateDatabase(server.Stardog.Database)

	r := mux.NewRouter().StrictSlash(false)


	r.HandleFunc("/ark:{prefix}", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				server.CreateArkNamespaceHandler(w, r)
				return
			}

			if r.Method == "GET" {
				server.GetArkNamespaceHandler(w, r)
				return
			}

			if r.Method == "PUT" {
				server.UpdateArkNamespaceHandler(w, r)
				return
			} else {
				http.Error(w, "Method Not Allowed", 405)
				return
			}
		}))

	r.PathPrefix("/shoulder/ark:{prefix}").Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				server.ArkMintHandler(w, r)
				return
			} else {
				http.Error(w, "Method Not Allowed", 405)
				return
			}
		}))

	r.PathPrefix("/ark:{prefix}/{suffix}").Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				server.ArkCreateHandler(w, r)
				return
			}

			if r.Method == "GET" {
				server.ArkResolveHandler(w, r)
				return
			}

			if r.Method == "PUT" {
				server.ArkUpdateHandler(w, r)
				return
			}
			if r.Method == "DELETE" {
				server.ArkDeleteHandler(w, r)
				return
			} else {
				http.Error(w, "Method Not Allowed", 405)
				return
			}
		}))

	r.HandleFunc("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "ok"}`))
		}))

	log.Fatal(http.ListenAndServe(":8080", r))

}
