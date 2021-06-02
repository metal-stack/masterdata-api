module github.com/metal-stack/masterdata-api

go 1.16

require (
	github.com/Masterminds/squirrel v1.5.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/status v1.1.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/json-iterator/go v1.1.11
	github.com/lib/pq v1.10.2
	github.com/lopezator/migrator v0.3.0
	github.com/metal-stack/metal-lib v0.8.0
	github.com/metal-stack/security v0.5.3
	github.com/metal-stack/v v1.0.3
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/testcontainers/testcontainers-go v0.11.0
	go.uber.org/zap v1.17.0
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

// required because by default viper depends on etcd v3.3.10 which has a corrupt sum
replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.18+incompatible
