package main

import (
	"testing"
)



func TestStardog(t *testing.T) {


	var s = StardogServer{
		URI: "http://localhost:5820",
		Password: "admin",
		Username: "admin",
		Database: "testing",
	}

	t.Run("Database", func(t *testing.T){

		t.Run("Create", func(t *testing.T){

			response, statusCode, err := s.createDatabase(s.Database)

			if err != nil {
				t.Fatalf("Failed to Create Database\nStatusCode: %d\nResponse: %s", statusCode, string(response))
			}

			t.Logf("Success Created Database\nStatusCode: %d\nResponse: %s", statusCode, string(response))

		})

		t.Run("Delete", func(t *testing.T){
			response, err := s.dropDatabase(s.Database)

			if err != nil {
				t.Fatalf("Failed to Drop Database\nResponse: %s", string(response))
			}

			t.Logf("Successfully Dropped Database\nResponse: %s", string(response))

		})
	})

	s.createDatabase(s.Database)
	identifier := []byte(`{"@id": "ark:/99999/identifier-test", "@context": {"@vocab": "http://schema.org/"}, "name": "identifier-test"}`)
	t.Run("Identifier", func(t *testing.T){
		t.Run("Transaction", func(t *testing.T){
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

		t.Run("Create", func(t *testing.T){
			err := s.AddIdentifier(identifier)

			if err != nil {
				t.Fatalf("Failed to Add Identifier: %s", err.Error())
			}
		})

		t.Run("Delete", func(t *testing.T){
				err := s.RemoveIdentifier(identifier)

				if err != nil {
					t.Fatalf("Failed to Add Identifier: %s", err.Error())
				}
			})
	})

	// var namedGraph = "ark:/99999/test-named-graph"
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
