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
* tenant project
  * settings
    * limit max clusters
* tenant members
* tenant project members

## Design

The services are exposed as grpc-services. The included client can be used
in other services by simply importing it.

The data is stored in a generic way using a postgres database
with tables consisting of id and json-document fields.

Changes to the data are reflected in a history table-twin per entity. When data
is created, updated or deleted, the change is also written to the history table.

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

```bash
go tool pprof -http :8080 localhost:2113/debug/pprof/heap
go tool pprof -http :8080 localhost:2113/debug/pprof/goroutine
```

## Generics migration

In order to get rid of all the reflection based logic in `postgres.go`, we decided to migrate to generics which are available since go 1.18.
This leads to much nicer code to read and also brings some benefits regarding allocations. Performance is at the same level as the reflection based approach.

To measure the impact, a bunch of benchmarks have been implemented for all CRUD operations provided by `postgres.go`.

Results comparing old(reflection based) vs. new(generics based):

```plain

name             old time/op    new time/op    delta
GetTenant-16       92.9µs ±11%    94.4µs ± 5%     ~     (p=0.421 n=5+5)
CreateTenant-16    3.06ms ± 9%    3.40ms ± 4%  +10.95%  (p=0.008 n=5+5)
UpdateTenant-16    3.59ms ± 9%    3.81ms ±19%     ~     (p=0.548 n=5+5)
FindTenant-16       259µs ±12%     224µs ± 3%  -13.75%  (p=0.008 n=5+5)

name             old alloc/op   new alloc/op   delta
GetTenant-16       5.68kB ± 0%    4.40kB ± 0%  -22.55%  (p=0.029 n=4+4)
CreateTenant-16    10.8kB ± 0%     9.6kB ± 0%  -11.22%  (p=0.008 n=5+5)
UpdateTenant-16    22.7kB ± 0%    19.0kB ± 0%  -16.26%  (p=0.008 n=5+5)
FindTenant-16      7.15kB ± 0%    5.19kB ± 0%  -27.38%  (p=0.016 n=4+5)

name             old allocs/op  new allocs/op  delta
GetTenant-16          118 ± 0%        92 ± 0%  -22.03%  (p=0.008 n=5+5)
CreateTenant-16       238 ± 0%       204 ± 0%  -14.29%  (p=0.008 n=5+5)
UpdateTenant-16       500 ± 0%       408 ± 0%  -18.40%  (p=0.008 n=5+5)
FindTenant-16         146 ± 0%       108 ± 0%  -26.03%  (p=0.008 n=5+5)

```

As shown, performance is about the same, but allocations in terms of bytes and count have been reduced quite significant.
