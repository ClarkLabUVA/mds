package main

import (
	bson "go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
)

func TestGeneralNested(t *testing.T) {
	inputUpdate := []byte(`{"hello":{"world": {"goodnight": "moon"} } }`)
	bsonUpdate := nestedUpdate(inputUpdate)

	var dotted map[string]interface{}
	err := bson.Unmarshal(bsonUpdate, &dotted)
	if err != nil {
		t.Fatal("Failed to unmarshal update in dot notation", err)
	}

	t.Logf("Unmarshaled Map: %+v", dotted)

	val, ok := dotted["$set"]

	if !ok {
		t.Fatal("dotted.$set is unset: ", val)
	}

	if reflect.ValueOf(val).Kind() != reflect.Map {
		t.Fatal("dotted.$set is not a map", val)
	}

	if moon := val.(map[string]interface{})["hello.world.goodnight"]; moon != "moon" {
		t.Fatal("Iincorrect Value Set: ", moon)
	}

}

func TestMongoUpdate(t *testing.T) {
	//TODO attempt to ping mongo
	guid := "ark:99999/test"
	namespace := []byte(`{"name": "test namespace", "@type": "namespace"}`)
	docBytes := []byte(`{"id": "t", "nested": {"doc": {"id": "init", "other": "props"}, "other": "props"} }`)
	update := []byte(`{"nested": {"doc": {"id": "test"}}}`)

	// create namespace
	CreateNamespace(namespace, "ark:99999")

	// create identifier
	CreateIdentifier(docBytes, guid, User{})

	// update identifier
	_, err := UpdateIdentifier(guid, update)

	if err != nil {
		DeleteIdentifier(guid)
		t.Fatal("Failed to Update Identifier", err)
	}

	res, err := GetIdentifier(guid)
	if err != nil {
		DeleteIdentifier(guid)
		t.Fatal("Failed to Get Identifier: ", err)
	}

	// delete the identifier
	DeleteIdentifier(guid)
	t.Logf("Found Identifier: %+v", string(res))
}
