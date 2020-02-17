# masterdata-api

This Microservice provides the source of truth for master data.

* tenant master-data
  * name/desc
* tenant settings
  * cloud
    * limit max projects
    * limit max clusters
  * cluster
    * limit
  * authn/authz & iam
    * tenant-id
    * issuer-url
    * issuer-client-id
    * idm
      * type AD/UX - relevant for config of restrictions
    * idm-webhook
      * enabled true/false
      * webhook config
        * type AD/UX (duplicate, see above)
        * subsystem
        * ...
* tenant project
  * settings
    * limit max clusters

## Design

The services are exposed as grpc-services. The included client can be used
in other services by simply importing it.

The data is stored in a generic way using a postgres database
with tables consisting of id /json-document fields.

The main entities are generated from a `<type>.proto`-file
plus some additional mapping-code in a `<type>.go` file.
Using a go-generate-statement the db-schema and some boilerplating code
is generated using naming-conventions.

## Initial Data

It is possible to insert data on startup, this is done by placing one ore more yaml documents into the `initdb.d` directory.
Multi document yaml files are not supported at the moment. If the given version of the entity is lower or equal the entity version
stored in the database, no create or update happens. Otherwise a update is executed.
On every error happening during initdb is logged, but the affected entity is not processed.

## Build

```bash
make all
```

### Install protoc

```bash
* https://github.com/protocolbuffers/protobuf
* latest https://gist.github.com/sofyanhadia/37787e5ed098c97919b8c593f0ec44d8#gistcomment-2760267
```

### Install protoc-gen-go

```bash
go get -u github.com/golang/protobuf/protoc-gen-go
```

## Run

```bash
make postgres-up
```

Start client with extensive logging
```bash
make clean protoc client
GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info bin/client
```

Start server
```bash
make clean protoc server
bin/server
```

## Metrics

```bash
http://localhost:2112/metrics
```

## pprof

```
go tool pprof -http :8080 localhost:2113/debug/pprof/heap
go tool pprof -http :8080 localhost:2113/debug/pprof/goroutine
```
