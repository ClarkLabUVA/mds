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

func TestBackend(t *testing.T) {

	var backend = Backend{
		Stardog: StardogServer{
			URI:      "http://localhost:5820",
			Password: "admin",
			Username: "admin",
			Database: "testing",
		},
		Mongo: MongoServer{
			URI:      "mongodb://mongoadmin:mongosecret@localhost:27017",
			Database: "ors",
			Collection: "test",
		},
	}


	// drop the test database from mongo
	ctx, cancel, client, _ := backend.Mongo.connect()
	defer cancel()

	col := client.Database(backend.Mongo.Database).Collection(backend.Mongo.Collection)
	if cleanUpErr := col.Drop(ctx); cleanUpErr != nil {
		t.Fatalf("Failed CleanUp\tError: %s", cleanUpErr.Error())
	}


	namespaceGUID := "ark:9999"
	namespacePayload := []byte(`{"name": "test namespace"}`)

	identifierGUID := "ark:9999/test"
	identifierPayload := []byte(`{"@context": "https://schema.org/", "name": "test identifier", "@type": "Dataset"}`)

	t.Run("CreateNamespace", func(t *testing.T) {

		err := backend.CreateNamespace(namespaceGUID, namespacePayload)
		if err != nil {
			t.Fatalf("Failed to Create Namespace\tError: %s", err.Error())
		}

	})
	t.Run("GetNamespace", func(t *testing.T) {
		namespace, err := backend.GetNamespace(namespaceGUID)
		if err != nil {
			t.Fatalf("Failed to Get Namespace\tError: %s", err.Error())
		}

		t.Logf("Found Namespace with Content\t%s", string(namespace))
	})

	//t.Run("UpdateNamespace", func(t *testing.T) {})
	//t.Run("DeleteNamespace", func(t *testing.T) {})
	t.Run("CreateIdentifier", func(t *testing.T) {

		err := backend.CreateIdentifier(identifierGUID, identifierPayload, User{})
		if err != nil {
			t.Fatalf("Failed to Create Identifier\tError: %s", err.Error())
		}

	})
	t.Run("GetIdentifier", func(t *testing.T){

		payload, err := backend.GetIdentifier(identifierGUID)
		if err != nil {
			t.Fatalf("Failed to Get Identifier: %s", err.Error())
		}

		t.Logf("Found Identifier: %s", string(payload))

	})
	t.Run("UpdateIdentifier", func(t *testing.T){})




}

/*
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
*/
