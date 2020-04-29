# Notes for 


## Deployment

get a license

```
docker run -it -v ~/stardog-home/:/var/opt/stardog -p 5820:5820 stardog/stardog

# get the license locally
docker run -d -v /Users/mal8ch-admin/Dev/ORS+/mds/bin/:/var/opt/stardog -p 5820:5820 stardog/stardog

# create a config map for 

```

# Todo list

On Project Creation
  - update Org -> subOrganization points to proj

On Identifier Creation
  - update Proj -> identifiers appended {"@id": ,"@type": "", "name": "", "author": "",}

Calls to Metadata Validation Service

Landing Service
  - Table view
  - Full Text Search
  - Advanced Search, (by @type, project, organization)
  - Graph View of JSON-LD


Inserts into Stardog


- SPARQL Search



fields to modify
√ @id
√ @context
- identifiers
- url
- sdPublisher
- sdPublicationDate

default
