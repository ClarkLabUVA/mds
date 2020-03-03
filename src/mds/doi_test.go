package main

import (
	"testing"
)

func TestDOIConstructor(t *testing.T) {

	// test DOI constructor
	t.Run("Success", func(te *testing.T) {
		content := []byte(
			`{  
    "@context":"http://schema.org",
    "@type":"SoftwareSourceCode",
    "@id": "https://doi.org/10.5438/qeg0-3gm3",
    "url":"https://github.com/datacite/maremma",
    "name":"Maremma: a Ruby library for simplified network calls",
    "author":{  
    "@type":"person",
    "@id":"http://orcid.org/0000-0003-0077-4738",
    "name":"Martin Fenner"
    },
    "description":"Simplifies network calls, including json/xml parsing and error handling. Based on Faraday.",
    "keywords":"faraday, excon, net/http",
    "dateCreated":"2015-11-28",
    "datePublished":"2017-02-24",
    "dateModified":"2017-02-24",
    "publisher":{  
      "@type":"Organization",
      "name":"DataCite"
    }
}`)

		url := "https://example.org"
		identifier := "10.5072/1234"
		doi, err := NewDOI(identifier, content, url)

		if err != nil {
			te.Fatalf("Constructor Failed to Produce DOI\n\tError: %s", err.Error())
		}

		if doi.Identifier != identifier || doi.URL != url {
			te.Fatalf("Constructor Failed")
		}
		te.Logf("XML: %s", string(doi.DataciteXML))
	})

	t.Run("NoMetadata", func(te *testing.T) {
		content := []byte(``)
		url := "https://example.org"
		identifier := "10.5072/1234"
		doi, err := NewDOI(identifier, content, url)

		if err != nil {
			te.Fatalf("Failed to Convert Metadata")
		}
		te.Logf("XML: %s", string(doi.DataciteXML))
	})

}

func TestDOIDatacite(t *testing.T) {

	content := []byte(`{  
    "@context":"http://schema.org",
    "@type":"SoftwareSourceCode",
    "@id": "https://doi.org/10.5438/qeg0-3gm3",
    "url":"https://github.com/datacite/maremma",
    "name":"Maremma: a Ruby library for simplified network calls",
    "author":{  
    "@type":"person",
    "@id":"http://orcid.org/0000-0003-0077-4738",
    "name":"Martin Fenner"
    },
    "description":"Simplifies network calls, including json/xml parsing and error handling. Based on Faraday.",
    "keywords":"faraday, excon, net/http",
    "dateCreated":"2015-11-28",
    "datePublished":"2017-02-24",
    "dateModified":"2017-02-24",
    "publisher":{  
      "@type":"Organization",
      "name":"DataCite"
    }
    }`)

	url := "https://example.org"
	identifier := "10.70020/90820"
	doi, err := NewDOI(identifier, content, url)

	if err != nil {
		t.Fatalf("Failed to Convert Metadata: %s", err.Error())
	}

	t.Run("PutMetadata", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			err = doi.datacitePutMetadata()
			if err != nil {
				t.Fatalf("Failed to Add Metadata to Datacite: %s", err.Error())
			}
		})
		t.Run("Failure", func(t *testing.T) {
			content := []byte(`{  
	    "@context":"http://schema.org",
	    "@type":"SoftwareSourceCode",
	    "@id": "https://doi.org/10.5438/qeg0-3gm3",
	    "url":"https://github.com/datacite/maremma",
	    "name":"Maremma: a Ruby library for simplified network calls",
	    "description":"Simplifies network calls, including json/xml parsing and error handling. Based on Faraday.",
	    "keywords":"faraday, excon, net/http",
	    "dateCreated":"2015-11-28",
	    "dateModified":"2017-02-24",
	    "publisher":{  
	      "@type":"Organization",
	      "name":"DataCite"
	    }
	}`)

			url := "https://example.org"
			identifier := "10.70020/90820"
			failDOI, err := NewDOI(identifier, content, url)

			err = failDOI.datacitePutMetadata()
			if err == nil {
				t.Fatalf("Updated Metadata with ineadequate fields")
			}

			t.Logf("MissingMetadata Error: %s", err.Error())
		})

	})

	t.Run("PutResolver", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			err := doi.datacitePutResolver()

			if err != nil {
				t.Fatalf("Failed to PutResolver: %s", err.Error())
			}
		})
		t.Run("Failure", func(t *testing.T) {

			t.Run("BeforeMetadata", func(t *testing.T) {
				fail := DOI{URL: "http://example.org", Identifier: "10.70020/837410"}

				err := fail.datacitePutResolver()

				if err == nil {
					t.Fatalf("Put Resolver Succeeded without Metadata")
				}

				t.Logf("Error: %s", err.Error())
			})
			t.Run("MalformedURL", func(t *testing.T) {
				fail := DOI{URL: "h://no.", Identifier: doi.Identifier}
				err := fail.datacitePutResolver()

				if err == nil {
					t.Fatalf("Put Resolver Succeeded with malformed link")
				}
				t.Logf("Error: %s", err.Error())
			})

		})

	})

	// func TestDOIUpdate(t *testing.T) {}

	// func TestDOICreate(t *testing.T) {}

	// func TestDOIDelete(t *testing.T) {}

}
