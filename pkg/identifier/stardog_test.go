//© 2020 By The Rector And Visitors Of The University Of Virginia

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package identifier

import (
	"testing"
)

func TestStardog(t *testing.T) {

	var s = StardogServer{
		URI:      "http://localhost:5820",
		Password: "admin",
		Username: "admin",
		Database: "testing",
	}

	t.Run("Database", func(t *testing.T) {

		t.Run("Create", func(t *testing.T) {

			response, statusCode, err := s.CreateDatabase(s.Database)

			if err != nil {
				t.Fatalf("Failed to Create Database\nStatusCode: %d\nResponse: %s", statusCode, string(response))
			}

			t.Logf("Success Created Database\nStatusCode: %d\nResponse: %s", statusCode, string(response))

		})

		t.Run("Delete", func(t *testing.T) {
			response, err := s.DropDatabase(s.Database)

			if err != nil {
				t.Fatalf("Failed to Drop Database\nResponse: %s", string(response))
			}

			t.Logf("Successfully Dropped Database\nResponse: %s", string(response))

		})
	})

	s.createDatabase(s.Database)
	identifier := []byte(`{"@id": "ark:/99999/identifier-test", "@context": {"@vocab": "http://schema.org/"}, "name": "identifier-test"}`)
	t.Run("Identifier", func(t *testing.T) {
		t.Run("Transaction", func(t *testing.T) {
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

		t.Run("Create", func(t *testing.T) {
			err := s.AddIdentifier(identifier)

			if err != nil {
				t.Fatalf("Failed to Add Identifier: %s", err.Error())
			}
		})

		t.Run("Delete", func(t *testing.T) {
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
