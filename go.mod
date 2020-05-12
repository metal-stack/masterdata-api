module github.com/metal-stack/masterdata-api

go 1.13

require (
	github.com/Masterminds/squirrel v1.3.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-pg/migrations/v7 v7.1.10
	github.com/gogo/status v1.1.0
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.9
	github.com/lib/pq v1.5.0
	github.com/lopezator/migrator v0.3.0
	github.com/metal-stack/metal-lib v0.3.5
	github.com/metal-stack/security v0.3.0
	github.com/metal-stack/v v1.0.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/testcontainers/testcontainers-go v0.5.1
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20200501053045-e0ff5e5a1de5
	google.golang.org/genproto v0.0.0-20200430143042-b979b6f78d84 // indirect
	google.golang.org/grpc v1.29.1
)

// required because by default viper depends on etcd v3.3.10 which has a corrupt sum
replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.18+incompatible
