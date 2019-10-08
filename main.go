package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"

	"log"
	"os"

	"encoding/json"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"strings"
)

var MS = MongoServer{
	URI:	"mongodb://mongoadmin:mongosecret@localhost:27017",
	Database: "ors",
}

func init() {

	mongoURI, exists := os.LookupEnv("MONGO_URI")
	if exists {
		MS.URI = mongoURI
	}


	mongoDB, exists := os.LookupEnv("MONGO_DB")
	if exists {
		MS.Database = mongoDB
	}

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

	err = CreateNamespace(payload, guid)
	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write([]byte(`{"created": "` + guid + `"}`))

	case ErrAlreadyExists:
		serveJSON(w, 400, map[string]interface{}{"error": err.Error()})

	default:
		serveJSON(w, 500, map[string]interface{}{"error": err.Error(), "message": "Error Creating Namespace"})

	}

	return

}

func GetArkNamespaceHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"]

	ns, err := GetNamespace(guid)

	switch err {

	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(ns)

	case ErrNilDocument:
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

	response, err := UpdateNamespace(update, guid)

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

	identifier, err := GetIdentifier(guid)

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

	var u User
	err = CreateIdentifier(bodyBytes, guid, u)

	switch err {
	case nil:
		serveJSON(w, 201, map[string]interface{}{"created": guid})

	case ErrNoNamespace:
		serveJSON(w, 404, map[string]interface{}{"error": "Namespace ark:" + namespace + " does not exist"})

	case ErrAlreadyExists:
		serveJSON(w, 400, map[string]interface{}{"error": "Identifier ark:" + guid + " already exists"})

	case ErrInvalidMetadata:
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
	identifierUUID, err := uuid.NewV4()
	if err != nil {
		serveJSON(w, 500, map[string]interface{}{"error": err.Error(), "message": "Failed to Generate UUID"})
		return
	}

	// append to identifier
	guid := "ark:" + vars["prefix"] + "/" + identifierUUID.String()

	// store identifier record
	var u User
	err = CreateIdentifier(bodyBytes, guid, u)

	switch err {

	case nil:
		serveJSON(w, 201, map[string]interface{}{"created": guid})

	case ErrNoNamespace:
		serveJSON(w, 404, map[string]interface{}{"error": "Namespace ark:" + vars["prefix"] + " does not exist"})

	case ErrInvalidMetadata:
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

	identifier, err := UpdateIdentifier(guid, update)

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

	identifier, err := DeleteIdentifier(guid)

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
