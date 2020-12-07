# MDS API Documentation 

#### The metadata service handles minting and resolving of identifiers in FAIRSCAPE. The provided metadata is stored in [Mongo](https://www.mongodb.com/cloud/atlas) and [Stardog](https://www.stardog.com/). Stardog is optional, but required for the evidence graph features. 

# Endpoints
 - **/ark:{prefix}**
 - **/shoulder/ark:{namespace}**
 - **/ark:{namespace}/{Identifier}**

# /ark:{prefix}

Handles creating and editing namespaces to store identifiers. 
{prefix} is 5 digit ark namespace prefix

## GET

Get namespace description. Returns metadata describing requested namespace.

```console
$ curl http://clarklab.uvarc.io/ark:99999 
```

## POST

Create a namespace. 

### Parameters 

 - **Name**
 - **Description**
 - owner


```bash
$ curl --request POST \
  --url https://clarklab.uvarc.io/mds/ark:99999 \
  --header 'Authorization: Bearer YOUR_JWT' \
  --header 'Content-Type: application/json' \
  --data '{"name":"Test","description":"Test Namespace"}'
```
## PUT

Update a namespace

### Parameters 

 - **Name**
 - **Description**
 - owner


```bash
$ curl --request PUT \
  --url https://clarklab.uvarc.io/mds/ark:99999 \
  --header 'Authorization: Bearer YOUR_JWT' \
  --header 'Content-Type: application/json' \
  --data '{"name":"New Name","description":"Updated Namespace"}'
```
  
  # /shoulder/ark:{prefix}
  
  ## POST
Mint an Identifier with associated json-ld metadata. Returns minted PID
### Parameters 
None required post json-ld metadata. 
 - Name
 - Description
 - Type


```bash
$ curl --request POST \
  --url https://clarklab.uvarc.io/mds/shoulder/ark:99999 \
  --header 'Authorization: Bearer YOUR_JWT' \
  --header 'Content-Type: application/json' \
  --data '{"name":"Example Dataset", "@type":"Datatset", "description":"Example made up data"}'
``` 
# /ark:{prefix}/{suffix}

## GET

Resolve ARK

```console
$ curl http://clarklab.uvarc.io/ark:99999/ra1-ndom-32-ark 
```

  ## POST
Mint an Identifier with specified PID and json-ld metadata.
 
### Parameters 
None required post json-ld metadata. 
 - Name
 - Description
 - Type


```bash
$ curl --request POST \
  --url https://clarklab.uvarc.io/mds/ark:99999/test-id \
  --header 'Authorization: Bearer YOUR_JWT' \
  --header 'Content-Type: application/json' \
  --data '{"name":"Example Dataset", "@type":"Datatset", "description":"Example made up data"}'
```
  ## PUT
Update metadata of a previously minted identifier.
 
### Parameters 
None required put json-ld metadata. 
 - Name
 - Description
 - Type


```bash
$ curl --request PUT \
  --url https://clarklab.uvarc.io/mds/ark:99999/test-id \
  --header 'Authorization: Bearer YOUR_JWT' \
  --header 'Content-Type: application/json' \
  --data '{"name":"Example Dataset", "@type":"Datatset", "description":"Example made up data"}'
```
