package main

import (
	"testing"
)


func TestUpdateguid(t *testing.T) {

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
