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

func TestNamespace(t *testing.T) {

	t.Run("Create", func(t *testing.T) {
		namespace := []byte(`{
			"@id": "ark:99999",
			"@context": {"@vocab": "http://schema.org/"},
			"name": "test namespace"
		}`)

		err := CreateNamespace(namespace, "ark:99999")

		if err != nil {
			t.Fatalf("Create Namespace Failed: %s", err.Error())
		}
	})

	//t.Run("Update", func(t *testing.T){})

	t.Run("Get", func(t *testing.T) {
		response, err := GetNamespace("ark:99999")

		if err != nil {
			t.Fatalf("Failed to Get Namespace: %s", err.Error())
		}

		t.Logf("Got Namespace: %s", string(response))
	})

	t.Run("Delete", func(t *testing.T) {
		response, err := DeleteNamespace("ark:99999")

		if err != nil {
			t.Fatalf("Failed to Delete Namespace: %s", err.Error())
		}

		t.Logf("Delete Namespace: %s", string(response))
	})

}

func TestIdentifier(t *testing.T) {

	namespace := "ark:90909"
	namespace_payload := []byte(`{
		"@context": {"@vocab": "http://schema.org/"},
		"name": "test namespace"
	}`)

	err := CreateNamespace(namespace_payload, namespace)

	if err != nil {
		t.Fatalf("Create Namespace Failed: %s", err.Error())
	}

	guid := "ark:90909/test"
	payload := []byte(`{
		"@context": {"@vocab": "http://schema.org/"},
		"name": "TestID",
		"@type": "Dataset"
	}`)

	t.Run("Create", func(t *testing.T) {
		var u User
		err := CreateIdentifier(payload, guid, u)
		if err != nil {
			t.Fatalf("Failed to Create Identifier: %s", err.Error())
		}
	})

	t.Run("Update", func(t *testing.T) {

		update := []byte(`{"name": "UpdatedName", "newprop": "newval"}`)
		response, err := UpdateIdentifier(guid, update)
		if err != nil {
			t.Fatalf("Failed to Update Identifier: %s", err.Error())
		}

		t.Logf("Updated Identifier %s: %s", guid, string(response))

	})

	t.Run("Get", func(t *testing.T) {
		response, err := GetIdentifier(guid)
		if err != nil {
			t.Fatalf("Failed to Get Identifier: %s", err.Error())
		}

		t.Logf("Retrieved Identifier %s: %s", guid, string(response))
	})

	t.Run("Delete", func(t *testing.T) {
		response, err := DeleteIdentifier(guid)
		if err != nil {
			t.Fatalf("Failed to Delete Identifier: %s", err.Error())
		}

		t.Logf("Deleted Identifier %s: %s", guid, string(response))

	})

	_, err = DeleteNamespace(namespace)
	if err != nil {
		t.Fatalf("Failed to Delete Namespace %s: %s", namespace, err.Error())
	}

}
