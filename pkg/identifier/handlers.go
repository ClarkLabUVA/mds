package identifier

import (
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"strings"
	"github.com/google/uuid"
	"encoding/json"
)


// CreateArkNamespaceHandler is the http handler for creating identifier namespaces
func (b *Backend) CreateArkNamespaceHandler(w http.ResponseWriter, r *http.Request) {

    /*
	// extract user from request context
	var u User
	contextUser := r.Context().Value("user")
	u = contextUser.(User)

	// if user is not an admin return a 403 error
	if u.Role != "admin" {
		serveJSON(w, 403, map[string]interface{}{"error": "action not permitted", "message": "only admins may create ark namespaces"})
		return 
	}
    */

	// read in response from request
	payload, err := ioutil.ReadAll(r.Body)

	if err != nil {
		serveJSON(w, 400, map[string]interface{}{"error": err.Error(), "message": "Error reading in payload"})
		return
	}

	// get vars from path
	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"]

	err = b.CreateNamespace(guid, payload)
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


//GetArkNamespaceHandler is the http handler for getting identifier namespaces
func (b *Backend) GetArkNamespaceHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"]

	ns, err := b.GetNamespace(guid)

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


//UpdateArkNamespaceHandler is the http handler for update identifier namespaces
func (b *Backend) UpdateArkNamespaceHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

    /*
	// extract user from request context
	var u User
	contextUser := r.Context().Value("user")
	u = contextUser.(User)

	// if user is not an admin return a 403 error
	if u.Role != "admin" {
		serveJSON(w, 403, map[string]interface{}{"error": "action not permitted", "message": "only admins may create ark namespaces"})
		return 
	}
    */


	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"]

	update, err := ioutil.ReadAll(r.Body)

	response, err := b.UpdateNamespace(guid, update)

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


//ArkResolveHandler 
func (b *Backend) ArkResolveHandler(w http.ResponseWriter, r *http.Request) {

	guid := strings.TrimPrefix(r.RequestURI, "/")

	identifier, err := b.GetIdentifier(guid)

	if err != nil {
		serveJSON(w, 500, map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(identifier)
	return

}


//ArkCreateHandler
func (b *Backend) ArkCreateHandler(w http.ResponseWriter, r *http.Request) {

	var u User
    /*
	// extract user from request context
	var err error
	contextUser := r.Context().Value("user")
	u = contextUser.(User)
	
	if u.Role != "admin" && u.Role != "user" {
		serveJSON(w, 403, map[string]interface{}{"error": "action not permitted", "message": "must be a user or  may create ark identifiers"})
		return 
	}

    */

	guid := strings.TrimPrefix(r.RequestURI, "/")

	splitPath := strings.Split(guid, "/")
	namespace := splitPath[0]

    /*
	// create resource in auth service
	err = AuthCreateACL(guid, u)

	// if error is found in the auth service
	if err != nil {
		serveJSON(w, 500, map[string]interface{}{"error": err.Error(), "message": "Error registering ACL for identifier"})
		return
	}
    */

	// read in response from request
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		serveJSON(w, 400, map[string]interface{}{"error": err.Error(), "message": "Error reading in payload"})
		return
	}

	err = b.CreateIdentifier(guid, bodyBytes, u)

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


//ArkMintHandler
func (b *Backend) ArkMintHandler(w http.ResponseWriter, r *http.Request) {

	var u User
    /*
	// extract user from request context
	contextUser := r.Context().Value("user")
	u = contextUser.(User)

	if u.Role != "admin" && u.Role != "user" {
		serveJSON(w, 403, map[string]interface{}{"error": "action not permitted", "message": "must be a user or  may create ark identifiers"})
		return 
	}
    */

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

    /*
	// create resource in auth service
	err = AuthCreateACL(guid, u)
    */

	// if error is found in the auth service
	if err != nil {
		serveJSON(w, 500, map[string]interface{}{"error": err.Error(), "message": "Error registering ACL for identifier"})
		return
	}
	
	// store identifier record
	err = b.CreateIdentifier(guid, bodyBytes, u)

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


//ArkUpdateHandler
func (b *Backend) ArkUpdateHandler(w http.ResponseWriter, r *http.Request) {

	// get vars from path
	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"] + "/" + vars["suffix"]

    /*
	// extract user from request context
	var u User
	contextUser := r.Context().Value("user")
	u = contextUser.(User)

	// get permissions on identifier from auth service
	acl, err := AuthGetACL(guid)

	if err != nil {
		serveJSON(w, 500, map[string]interface{}{"error": err.Error(), "message": "Error retrieving permissions from auth service"})
		return
	}


	// check if user is in allowed
	if u.AllowedAccess(acl) != true {
		serveJSON(w, 403, map[string]interface{}{"error": "user is unauthorized to preform this action"})
		return
	}
    */


	// read in response from request
	update, err := ioutil.ReadAll(r.Body)

	if err != nil {
		serveJSON(w, 400, map[string]interface{}{"error": err.Error(), "message": "Error reading in payload"})
		return
	}



	identifier, err := b.UpdateIdentifier(guid, update)

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


//ArkDeleteHandler
func (b *Backend) ArkDeleteHandler(w http.ResponseWriter, r *http.Request) {

	// get vars from path
	vars := mux.Vars(r)
	guid := "ark:" + vars["prefix"] + "/" + vars["suffix"]
	
    /*
	// extract user from request context
	var u User
	contextUser := r.Context().Value("user")
	u = contextUser.(User)

	// get permissions on identifier from auth service
	acl, err := AuthGetACL(guid)

	if err != nil {
		serveJSON(w, 500, map[string]interface{}{"error": err.Error(), "message": "Error retrieving permissions from auth service"})
		return
	}

	// check if user is in allowed
	if u.AllowedAccess(acl) != true {
		serveJSON(w, 403, map[string]interface{}{"error": "user is unauthorized to preform this action"})
		return
	}
    */

	identifier, err := b.DeleteIdentifier(guid)

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
