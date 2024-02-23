export GO111MODULE := on
export CGO_ENABLED := 0
DOCKER_TAG := $(or ${GITHUB_TAG_NAME}, latest)

SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
# gnu date format iso-8601 is parsable with Go RFC3339
BUILDDATE := $(shell date --iso-8601=seconds)
VERSION := $(or ${VERSION},$(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD))

.PHONY: all
all: generate test server client

.PHONY: release
release: generate test server client

.PHONY: clean
clean:
	rm -f api/v1/*pb.go bin/*

.PHONY: protoc
protoc:
	make -C proto protoc

.PHONY: test
test:
	CGO_ENABLED=1 go test -cover -race -timeout 30s ./...

.PHONY: bench
bench:
	cd pkg/datastore && CGO_ENABLED=1 go test -bench=. -run=- -benchmem -count 5 && cd -

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate ./...

.PHONY: server
server: protoc generate
	go build -tags netgo -ldflags "-X 'github.com/metal-stack/v.Version=$(VERSION)' \
								   -X 'github.com/metal-stack/v.Revision=$(GITVERSION)' \
								   -X 'github.com/metal-stack/v.GitSHA1=$(SHA)' \
								   -X 'github.com/metal-stack/v.BuildDate=$(BUILDDATE)'" \
						 -o bin/server server/main.go
	strip bin/server

.PHONY: client
client:
	go build -tags netgo -o bin/client client/main.go
	strip bin/client

.PHONY: mocks
mocks:
	rm -rf pkg/datastore/mocks api/v1/mocks
	docker run --rm \
		--user $$(id -u):$$(id -g) \
		-w /work \
		-v ${PWD}:/work \
		vektra/mockery:v2.42.0 -r --keeptree --dir api/v1 --output api/v1/mocks --all

	docker run --rm \
		--user $$(id -u):$$(id -g) \
		-w /work \
		-v ${PWD}:/work \
		vektra/mockery:v2.42.0 -r --keeptree --dir pkg/datastore --output pkg/datastore/mocks --all

.PHONY: postgres-up
postgres-up: postgres-rm
	docker run -d --name masterdatadb -p 5432:5432 -e POSTGRES_PASSWORD="password" -e POSTGRES_USER="masterdata" -e POSTGRES_DB="masterdata" postgres:16-alpine

.PHONY: postgres-rm
postgres-rm:
	docker rm -f masterdatadb || true

.PHONY: certs
certs:
	cd certs && cfssl gencert -initca ca-csr.json | cfssljson -bare ca -
	cd certs && cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile client-server server.json | cfssljson -bare server -
	cd certs && cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile client client.json | cfssljson -bare client -

.PHONY: mini-lab-push
mini-lab-push:
	docker build -t metalstack/masterdata-api:latest .
	kind --name metal-control-plane load docker-image metalstack/masterdata-api:latest
	kubectl --kubeconfig=$(MINI_LAB_KUBECONFIG) patch deployments.apps -n metal-control-plane masterdata-api --patch='{"spec":{"template":{"spec":{"containers":[{"name": "masterdata-api","imagePullPolicy":"IfNotPresent","image":"metalstack/masterdata-api:latest"}]}}}}'
	kubectl --kubeconfig=$(MINI_LAB_KUBECONFIG) delete pod -n metal-control-plane -l app=masterdata-api
