package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"log"
	"os"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"strings"
	"github.com/ClarkLabUVA/mds/pkg/identifier"
)


var server identifier.Backend

func init() {

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

	if stardogDB, exists := os.LookupEnv("STARDOG_URI"); exists {
		server.Stardog.Database = stardogDB
	}

	if stardogPassword, exists := os.LookupEnv("STARDOG_PASSWORD"); exists {
		server.Stardog.Password = stardogPassword
	}

	if stardogUsername, exists := os.LookupEnv("STARDOG_USERNAME"); exists {
		server.Stardog.Username = stardogUsername
	}

	server.Stardog.CreateDatabase(server.Stardog.Database)

	// Log Initilization Variables
	log.Printf("StardogURI: %s\tStardogUsername: %s\tStardogPassword: %s\tStardogDatabase: %s",
		server.Stardog.URI, server.Stardog.Username, server.Stardog.Password, server.Stardog.Database)

	log.Printf("MongoURI: %s\tMongoDatabase: %s\tMongoCollection: %s",
		server.Mongo.URI, server.Mongo.Database, server.Mongo.Collection)

}

func main() {

	r := mux.NewRouter().StrictSlash(false)

	n := negroni.New()

	r.HandleFunc("/ark:{prefix}", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				CreateArkNamespaceHandler(w, r)
				return
			}

			if r.Method == "GET" {
				GetArkNamespaceHandler(w, r)
				return
			}

			if r.Method == "PUT" {
				UpdateArkNamespaceHandler(w, r)
				return
			} else {
				http.Error(w, "Method Not Allowed", 405)
				return
			}
		}))

	r.PathPrefix("/shoulder/ark:{prefix}").Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				ArkMintHandler(w, r)
				return
			} else {
				http.Error(w, "Method Not Allowed", 405)
				return
			}
		}))

	r.PathPrefix("/ark:{prefix}/{suffix}").Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				ArkCreateHandler(w, r)
				return
			}

			if r.Method == "GET" {
				ArkResolveHandler(w, r)
				return
			}

			if r.Method == "PUT" {
				ArkUpdateHandler(w, r)
				return
			}
			if r.Method == "DELETE" {
				ArkDeleteHandler(w, r)
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

	n.UseHandler(r)

	log.Fatal(http.ListenAndServe(":80", n))

}

func CreateArkNamespaceHandler(w http.ResponseWriter, r *http.Request) {

	// read in response from request
	payload, err := ioutil.ReadAll(r.Body)

	if err != nil {
		serveJSON(w, 400, map[string]interface{}{"error": err.Error(), "message": "Error reading in payload"})
		return
	}

	// get vars from path
	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"]

	err = server.CreateNamespace(guid, payload)
	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write([]byte(`{"created": "` + guid + `"}`))

	case identifier.ErrAlreadyExists:
		serveJSON(w, 400, map[string]interface{}{"error": err.Error()})

	default:
		serveJSON(w, 500, map[string]interface{}{"error": err.Error(), "message": "Error Creating Namespace"})

	}

	return

}

func GetArkNamespaceHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"]

	ns, err := server.GetNamespace(guid)

	switch err {

	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(ns)

	case identifier.ErrNilDocument:
		serveJSON(w, 404, map[string]interface{}{"error": "Namespace Not Found"})

	default:
		serveJSON(w, 500, map[string]interface{}{"error": err.Error(), "message": "Error Resolving Namespace"})

	}

	return

}

func UpdateArkNamespaceHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"]

	update, err := ioutil.ReadAll(r.Body)

	response, err := server.UpdateNamespace(guid, update)

	switch err {
	case nil:
		w.WriteHeader(200)
		w.Write(response)

	default:
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		w.WriteHeader(500)

	}

	return

}

func ArkResolveHandler(w http.ResponseWriter, r *http.Request) {

	guid := strings.TrimPrefix(r.RequestURI, "/")

	identifier, err := server.GetIdentifier(guid)

	if err != nil {
		serveJSON(w, 500, map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(identifier)
	return

}

func ArkCreateHandler(w http.ResponseWriter, r *http.Request) {

	guid := strings.TrimPrefix(r.RequestURI, "/")

	splitPath := strings.Split(guid, "/")
	namespace := splitPath[0]

	// read in response from request
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		serveJSON(w, 400, map[string]interface{}{"error": err.Error(), "message": "Error reading in payload"})
		return
	}

	var u identifier.User
	err = server.CreateIdentifier(guid, bodyBytes, u)

	switch err {
	case nil:
		serveJSON(w, 201, map[string]interface{}{"created": guid})

	case identifier.ErrNoNamespace:
		serveJSON(w, 404, map[string]interface{}{"error": "Namespace ark:" + namespace + " does not exist"})

	case identifier.ErrAlreadyExists:
		serveJSON(w, 400, map[string]interface{}{"error": "Identifier ark:" + guid + " already exists"})

	case identifier.ErrInvalidMetadata:
		serveJSON(w, 400, map[string]interface{}{"error": err.Error(), "message": "Invalid Metadata"})

	default:
		serveJSON(w, 500, map[string]interface{}{"error": err.Error(), "message": "Error Creating Identifier"})
	}

	return
}

func ArkMintHandler(w http.ResponseWriter, r *http.Request) {

	// read in response from request
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		serveJSON(w, 400, map[string]interface{}{"error": err.Error(), "message": "Error reading in payload"})
		return
	}

	// get vars from path
	vars := mux.Vars(r)

	// create a uuid
	identifierUUID := uuid.New()

	// append to identifier
	guid := "ark:" + vars["prefix"] + "/" + identifierUUID.String()

	// store identifier record
	var u identifier.User
	err = server.CreateIdentifier(guid, bodyBytes, u)

	switch err {

	case nil:
		serveJSON(w, 201, map[string]interface{}{"created": guid})

	case identifier.ErrNoNamespace:
		serveJSON(w, 404, map[string]interface{}{"error": "Namespace ark:" + vars["prefix"] + " does not exist"})

	case identifier.ErrInvalidMetadata:
		serveJSON(w, 400, map[string]interface{}{"error": err.Error(), "message": "Invalid Metadata"})

	default:
		serveJSON(w, 500, map[string]interface{}{"error": err.Error(), "message": "Error Creating Identifier"})
	}

	return

}

func ArkUpdateHandler(w http.ResponseWriter, r *http.Request) {

	// read in response from request
	update, err := ioutil.ReadAll(r.Body)

	if err != nil {
		serveJSON(w, 400, map[string]interface{}{"error": err.Error(), "message": "Error reading in payload"})
		return
	}

	// get vars from path
	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"] + "/" + vars["suffix"]

	identifier, err := server.UpdateIdentifier(guid, update)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(`{"error": ` + err.Error() + `}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"updated": ` + string(identifier) + `}`))

	return

}

func ArkDeleteHandler(w http.ResponseWriter, r *http.Request) {

	// get vars from path
	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"] + "/" + vars["suffix"]

	identifier, err := server.DeleteIdentifier(guid)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(`{"error": ` + err.Error() + `}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"deleted": ` + string(identifier) + `}`))

	return

}

func serveJSON(w http.ResponseWriter, statusCode int, payload interface{}) {

	b, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)

}
