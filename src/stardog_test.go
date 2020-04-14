package main

import (
	"testing"
)

var namedGraph = "ark:/99999/test-named-graph"

var s = StardogServer{
	URI: "http://localhost:5820",
	Password: "admin",
	Username: "admin",
	Database: "testing",
}

func TestStardog(t *testing.T) {


	t.Run("CreateIdentifier", func(t *testing.T){
		txId, err := s.NewTransaction()
		if err != nil {
			t.Fatalf("Failed To Start Transaction: %s", err.Error())
		}

		t.Logf("Started Transaction: %s", txId)

		data := []byte(`{"@id": "ark:/99999/test-data", "@context": {"@vocab": "http://schema.org/"}, "name": "test-data"}`)
		err = s.AddData(txId, data, "")

		if err != nil {
			t.Fatalf("Transaction Failed to Add Data: %s", err.Error())
		}

		err = s.Commit(txId)

		if err != nil {
			t.Fatalf("Failed to Commit Transaction: %s", err.Error())
		}

		})

	t.Run("IdentifierOperations", func(t *testing.T){
		identifier := []byte(`{"@id": "ark:/99999/identifier-test", "@context": {"@vocab": "http://schema.org/"}, "name": "identifier-test"}`)
		err := s.AddIdentifier(identifier)

		if err != nil {
			t.Fatalf("Failed to Add Identifier: %s", err.Error())
		}

		err = s.RemoveIdentifier(identifier)

		if err != nil {
			t.Fatalf("Failed to Add Identifier: %s", err.Error())
		}
	})



	/*
	t.Run("NamedGraph", func(t *testing.T){

		txId, err := s.NewTransaction()
		if err != nil {
			t.Fatalf("Failed To Start Transaction: %s", err.Error())
		}

		t.Logf("Started Transaction: %s", txId)

		data := []byte(`{"@id": "ark:/99999/test-data", "@context": {"@vocab": "http://schema.org/"}, "name": "test-data"}`)
		err = s.AddData(txId, data, namedGraph)

		if err != nil {
			t.Fatalf("Transaction Failed to Add Data: %s", err.Error())
		}

		err = s.Commit(txId)

		if err != nil {
			t.Fatalf("Failed to Commit Transaction: %s", err.Error())
		}

	})
	*/

}


/*
func TestStardogTransactionRemoveData(t *testing.T) {

	txId, err := s.NewTransaction()
	if err != nil {
		t.Fatalf("Failed To Start Transaction: %s", err.Error())
	}

	t.Logf("Started Transaction: %s", txId)

	data := []byte(`{"@id": "ark:/99999/test-data", "@context": {"@vocab": "http://schema.org/"}, "name": "test-data"}`)
	err = s.RemoveData(txId, data, "")

	if err != nil {
		t.Fatalf("Transaction Failed to Add Data: %s", err.Error())
	}

	err = s.Commit(txId)

	if err != nil {
		t.Fatalf("Failed to Commit Transaction: %s", err.Error())
	}

}

func TestStardogTransactionRemoveDataGraphURI(t *testing.T) {

	txId, err := s.NewTransaction()
	if err != nil {
		t.Fatalf("Failed To Start Transaction: %s", err.Error())
	}

	t.Logf("Started Transaction: %s", txId)

	data := []byte(`{"@id": "ark:/99999/test-data", "@context": {"@vocab": "http://schema.org/"}, "name": "test-data"}`)
	err = s.RemoveData(txId, data, namedGraph)

	if err != nil {
		t.Fatalf("Transaction Failed to Add Data: %s", err.Error())
	}

	err = s.Commit(txId)

	if err != nil {
		t.Fatalf("Failed to Commit Transaction: %s", err.Error())
	}

}
*/
