module github.com/metal-stack/masterdata-api

go 1.13

require (
	github.com/Masterminds/squirrel v1.1.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/status v1.1.0
	github.com/golang/protobuf v1.3.3
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.9
	github.com/lib/pq v1.2.0
	github.com/metal-stack/metal-lib v0.0.0-20200214140254-ee6c31afaf28
	github.com/metal-stack/security v0.0.0-20200213101721-0aed751953ac
	github.com/metal-stack/v v1.0.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.2.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
	github.com/testcontainers/testcontainers-go v0.0.9
	github.com/vektra/mockery v0.0.0-20181123154057-e78b021dcbb5 // indirect
	go.uber.org/zap v1.13.0
	golang.org/x/sys v0.0.0-20200212091648-12a6c2dcc1e4 // indirect
	google.golang.org/genproto v0.0.0-20200212174721-66ed5ce911ce // indirect
	google.golang.org/grpc v1.27.1
)

// required because by default viper depends on etcd v3.3.10 which has a corrupt sum
replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.18+incompatible
