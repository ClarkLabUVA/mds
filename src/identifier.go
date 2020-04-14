package main

import (
	"encoding/json"
	"errors"
	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
	"fmt"
	"github.com/buger/jsonparser"
)

var ErrInvalidMetadata = errors.New("Metadata Document is Invalid")
var ErrNilDocument = errors.New("No Document was Found")
var ErrAlreadyExists = errors.New("Document Already Exists")
var ErrNoNamespace = errors.New("No Namespace Record Found")
var ErrMissingProp = errors.New("Instance is missing required properties")

var ErrJSONUnmarshal = errors.New("Failed to Unmarshal JSON")

type Backend struct {
	Mongo	MongoServer
	Stardog	StardogServer
	useStardog bool
}


func (b *Backend) CreateNamespace(guid string, payload []byte) (err error) {

	ns := make(map[string]interface{})
	err = json.Unmarshal(payload, &ns)

	if err != nil {
		return fmt.Errorf(`{"message": "%q", "error": "%s"}`, ErrJSONUnmarshal, err.Error())
	}

	ns["@id"] = guid
	ns["_id"] = guid

	bsonRecord, err := bson.Marshal(ns)

	if err != nil {
		return
	}

	err = b.Mongo.InsertOne(bsonRecord)

	if err != nil {
		_, foundErr := b.Mongo.FindOne(bson.D{{"_id", guid}})
		if foundErr == nil {
			err = ErrAlreadyExists
		}
	}

	return
}


func (b *Backend) GetNamespace(guid string) (response []byte, err error) {

	ns, err := b.Mongo.FindOne(bson.D{{"_id", guid}})

	if err != nil {
		return
	}

	response, err = json.Marshal(ns)

	return
}


func (b *Backend) UpdateNamespace(payload []byte, guid string) (response []byte, err error) { return }


func (b *Backend) DeleteNamespace(guid string) (response []byte, err error) { return }


func (b *Backend) CreateIdentifier(guid string, payload []byte,  author User) (err error) {

	guidSplit := strings.Split(guid, "/")
	_, err = b.GetNamespace(guidSplit[0])

	if err == mongo.ErrNoDocuments {
		return ErrNoNamespace
	}

	metadata, err := processMetadataWrite(payload, guid, author)
	if err != nil {
		return
	}

	// TODO validate identifier metadata



	// add to stardog
	err = b.Stardog.AddIdentifier(payload)
	if err != nil {
		return fmt.Errorf("Stardog Failed to Create Identifier: %s", err.Error())
	}

	// store identifier in Mongo
	var bsonRecord bson.D
	err = bson.UnmarshalExtJSON(metadata, true, &bsonRecord)

	if err != nil {
		return fmt.Errorf("Failed to Unmarshal JSON to BSON\tError: %s", err.Error())
	}

	err = b.Mongo.InsertOne(bsonRecord)

	// if insert fails check that identifier doesn't already exist
	if err != nil {
		_, foundErr := b.Mongo.FindOne(bson.D{{"_id", guid}})
		if foundErr == nil {
			err = ErrAlreadyExists
		}
	}

	return
}


func (b *Backend) GetIdentifier(guid string) (response []byte, err error) {

	record, err := b.Mongo.FindOne(bson.D{{"_id", guid}})

	if err != nil {
		return
	}

	response, err = json.Marshal(record)

	response = processMetadataRead(response)

	return
}


func (b *Backend) DeleteIdentifier(guid string) (response []byte, err error) {

	record, err := b.Mongo.DeleteOne(bson.D{{"_id", guid}})

	if err != nil {
		return
	}

	response, err = json.Marshal(record)

	// remove identifier from stardog
	err = b.Stardog.RemoveIdentifier(response)

	//response = processMetadataRead(response)

	return

}


func (b *Backend) UpdateIdentifier(guid string, update []byte) (response []byte, err error) {

	return
}


func processMetadataWrite(inputMetadata []byte, guid string, author User) (metadata []byte, err error) {

	// set @id
	metadata, err = jsonparser.Set(inputMetadata, []byte(`"`+ guid +`"`), "@id")
	if err != nil {
		return
	}

	metadata, err = jsonparser.Set(metadata, []byte(`"`+ guid +`"`), "_id")
	if err != nil {
		return
	}

	// set the default context
	// TODO: if object add property "@vocab": "http://schema.org"
	metadata, err = jsonparser.Set(metadata, []byte(`{"@vocab": "http://schema.org"}`), "@context")
	if err != nil {
		return
	}

	// set namespace
	guidSplit := strings.Split(guid, "/")
	metadata, err = jsonparser.Set(metadata, []byte(`"`+guidSplit[0]+`"`), "namespace")
	if err != nil {
		return
	}


	// set url
	metadata, err = jsonparser.Set(metadata, []byte(`"http://ors.uvadcos.io/` + guid + `"`), "url")
	if err != nil {
		return
	}


	// fill in author
	if author.ID != "" {
		metadata, err = jsonparser.Set(metadata, []byte(`"`+ author.ID +`"`), "sdPublisher", "@id")
		if err != nil {
			return
		}
	}


	if author.Name != "" {
		metadata, err = jsonparser.Set(metadata, []byte(`"`+ author.Name +`"`), "sdPublisher", "name")
		if err != nil {
			return
		}
	}


	// set sdPublicationDate
	now, err := time.Now().MarshalJSON()
	metadata, err = jsonparser.Set(metadata, now, "sdPublicationDate")


	// TODO if not set default to "version": 1
	// metadata["version"] = 1


	return
}

// Replace with buger/jsonparser
func processMetadataRead(metadata []byte) []byte {
	// delete _id property
	metadata = jsonparser.Delete(metadata, "_id")

	// delete namespace properties namespace
	metadata = jsonparser.Delete(metadata, "namespace")


	return metadata
}

var backend = Backend{
	Stardog: StardogServer{
		URI:      "http://stardog.uvadcos.io",
		Password: "admin",
		Username: "admin",
		Database: "testing",
	},
	Mongo: MongoServer{
		URI:      "mongodb://mongoadmin:mongosecret@localhost:27017",
		Database: "ors",
		Collection: "ids",
	},
}


/*
func CreateNamespace(payload []byte, guid string) (err error) {
	ns := make(map[string]interface{})
	err = json.Unmarshal(payload, &ns)

	if err != nil {
		return
	}

	ns["@id"] = guid
	ns["_id"] = guid

	bsonRecord, err := bson.Marshal(ns)

	if err != nil {
		return
	}

	err = MS.InsertOne(bsonRecord, COL)

	if err != nil {
		_, foundErr := MS.FindOne(bson.D{{"_id", guid}}, COL)
		if foundErr == nil {
			err = ErrAlreadyExists
		}
	}

	return
}

func GetNamespace(guid string) (response []byte, err error) {
	ns, err := MS.FindOne(bson.D{{"_id", guid}}, "ids")

	if err != nil {
		return
	}

	response, err = json.Marshal(ns)

	return
}

func UpdateNamespace(payload []byte, guid string) (response []byte, err error) {
	updateBSON, err := bson.MarshalExtJSON(payload, false, false)

	if err != nil {
		return
	}

	raw, err := MS.UpdateOne(bson.D{{"_id", guid}}, bson.D{{"$set", updateBSON}}, COL)
	if err != nil {
		return
	}

	response, err = json.Marshal(raw)
	return
}

func DeleteNamespace(guid string) (response []byte, err error) {
	raw, err := MS.DeleteOne(bson.D{{"_id", guid}}, "ids")
	response, err = json.Marshal(raw)
	return
}


// unmarshal BSON record to JSON
// and format metadata
// - full urls for identifiers
// - pops _id

// Q: TRANSACTION ATOMICITY FOR MULTIPLE SERVICES
// if fails
// - Mongo.deleteOne
// - Stardog.deleteOne
func CreateIdentifier(payload []byte, guid string, author User) (err error) {
	guidSplit := strings.Split(guid, "/")
	_, err = GetNamespace(guidSplit[0])

	if err == mongo.ErrNoDocuments {
		return ErrNoNamespace
	}

	metadata := make(map[string]interface{})

	err = json.Unmarshal(payload, &metadata)
	if err != nil {
		return
	}

	metadata = processMetadataWrite(metadata, guid, author)

	// validate identifier metadata

	// add to stardog
	err = Stardog.AddIdentifier(payload)
	if err != nil {
		return
	}

	// store identifier in Mongo
	bsonRecord, err := bson.Marshal(metadata)

	if err != nil {
		return
	}

	err = MS.InsertOne(bsonRecord, COL)

	// if insert fails check that identifier doesn't already exist
	if err != nil {
		_, foundErr := MS.FindOne(bson.D{{"_id", guid}}, COL)
		if foundErr == nil {
			err = ErrAlreadyExists
		}
	}

	return
}

func GetIdentifier(guid string) (response []byte, err error) {
	record, err := MS.FindOne(bson.D{{"_id", guid}}, COL)

	if err != nil {
		return
	}

	response, err = json.Marshal(processMetadataRead(record))

	return
}

// TODO: Delete Behavior based on creativeWorkStatus
//DRAFT -> remove the document
//PUBLIC|PRIVATE -> status becomes WITHDRAWN
func DeleteIdentifier(guid string) (response []byte, err error) {
	record, err := MS.DeleteOne(bson.D{{"_id", guid}}, COL)

	if err != nil {
		return
	}

	response, err = json.Marshal(record)

	// remove identifier from stardog
	err = Stardog.RemoveIdentifier(response)

	//response = processMetadataRead(response)

	return
}

// Handle updates to creativeWorkStatus
// Allowed (add to resolver rules)
// DRAFT -> PUBLIC
// DRAFT -> PRIVATE
// PUBLIC|PRIVATE -> WITHDRAWN
// Not Allowed
// PUBLIC -> PRIVATE|DRAFT
// PRIVATE -> DRAFT
func UpdateIdentifier(guid string, update []byte) (response []byte, err error) {
	// unmarshal bson in to Bson.D
	var updateD bson.D
	err = bson.Unmarshal(nestedUpdate(update), &updateD)

	raw, err := MS.UpdateOne(bson.D{{"_id", guid}},
		updateD,
		COL)

	if err != nil {
		return
	}

	response, err = bson.MarshalExtJSON(raw, false, false)
	return
}
*/

type User struct {
	ID    string `json:"@id" bson:"_id"`
	Type  string `json:"@type" bson:"@type"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}
