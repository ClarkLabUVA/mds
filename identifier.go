package main

import (
	"encoding/json"
	"errors"
	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"strings"
	"time"
)

var ErrInvalidMetadata = errors.New("Metadata Document is Invalid")
var ErrNilDocument = errors.New("No Document was Found")
var ErrAlreadyExists = errors.New("Document Already Exists")
var ErrNoNamespace = errors.New("No Namespace Record Found")
var ErrMissingProp = errors.New("Instance is missing required properties")

var COL = "ids"

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

func processMetadataWrite(metadata map[string]interface{}, guid string, author User) map[string]interface{} {
	// set @id
	metadata["@id"] = guid
	metadata["_id"] = guid

	// set @context
	if _, ok := metadata["@context"]; !ok {
		metadata["@context"] = map[string]string{"@base": "http://schema.org/"}
	}

	// set namespace
	guidSplit := strings.Split(guid, "/")
	metadata["namespace"] = guidSplit[0]

	// set url
	metadata["url"] = "http://ors.uvadcos.io/" + guid

	// set sdPublisher
	if author.ID != "" && author.Name != "" {
		metadata["sdPublisher"] = author
	}

	// set sdPublicationDate
	metadata["sdPublicationDate"] = time.Now()

	// set version
	metadata["version"] = 1

	// set identifierStatus
	if _, ok := metadata["identifierStatus"]; !ok {
		metadata["identifierStatus"] = "DRAFT"
	}

	return metadata
}

// unmarshal BSON record to JSON
// and format metadata
// - full urls for identifiers
// - pops _id
func processMetadataRead(metadata map[string]interface{}) map[string]interface{} {
	// del _id
	delete(metadata, "_id")

	// del namespace
	delete(metadata, "namespace")

	return metadata
}

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

type tuple struct {
	Key   string
	Value interface{}
}

func nestedUpdate(update []byte) (bsonUpdate []byte) {
	processedMap := make(map[string]interface{})
	updateMap := make(map[string]interface{})
	json.Unmarshal(update, &updateMap)

	resChan := make(chan tuple, 50)
	dotConvert("", updateMap, resChan)

	for {
		select {
		case elem := <-resChan:
			processedMap[elem.Key] = elem.Value
		default:
			close(resChan)
			bsonUpdate, _ = bson.Marshal(
				map[string]interface{}{
					"$set": processedMap,
				})
			return
		}
	}
}

func dotConvert(base string, input map[string]interface{}, res chan tuple) {
	for key, val := range input {
		var newBase string
		if base == "" {
			newBase = key
		} else {
			newBase = base + "." + key
		}
		if valType := reflect.ValueOf(val); valType.Kind() == reflect.Map {
			dotConvert(newBase, val.(map[string]interface{}), res)
		} else {
			//log.Println("$set: ", newBase, " ", val)
			res <- tuple{Key: newBase, Value: val}
		}
	}

	return

}

type User struct {
	ID    string `json:"@id" bson:"_id"`
	Type  string `json:"@type" bson:"@type"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}
