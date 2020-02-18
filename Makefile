export GO111MODULE := on
export CGO_ENABLED := 0
DOCKER_TAG := $(or ${GITHUB_TAG_NAME}, latest)

SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
BUILDDATE := $(shell date -Iseconds)
VERSION := $(or ${VERSION},devel)

.PHONY: all
all: protoc generate mocks test server client

.PHONY: release
release: generate test server client

.PHONY: clean
clean: 
	rm -f api/v1/*pb.go bin/*

.PHONY: protoc
protoc:
	# docker pull metalstack/builder
	docker run -it --rm -v ${PWD}/api:/work/api metalstack/builder protoc -I api/ api/v1/*.proto --go_out=plugins=grpc:api
	docker run -it --rm -v ${PWD}/api:/work/api metalstack/builder protoc -I api/ api/grpc/health/v1/*.proto --go_out=plugins=grpc:api

.PHONY: test
test:
	CGO_ENABLED=1 go test -cover -race -timeout 30s ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate ./...

.PHONY: server
server:
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
	go get github.com/vektra/mockery/.../
	cd api/v1 && mockery -name ProjectServiceClient && mockery -name TenantServiceClient && cd -
	cd pkg/datastore && mockery -name Storage && cd -

.PHONY: postgres-up
postgres-up: postgres-rm
	docker run -d --name masterdatadb -p 5432:5432 -e POSTGRES_PASSWORD="password" -e POSTGRES_USER="masterdata" -e POSTGRES_DB="masterdata" postgres:12-alpine

.PHONY: postgres-rm
postgres-rm:
	docker rm -f masterdatadb || true

.PHONY: certs
certs:
	cd certs && cfssl gencert -initca ca-csr.json | cfssljson -bare ca -
	cd certs && cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile client-server server.json | cfssljson -bare server -
	cd certs && cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile client client.json | cfssljson -bare client -
	
