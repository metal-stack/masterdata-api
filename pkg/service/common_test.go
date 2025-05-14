package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pgContainer testcontainers.Container
)

func StartPostgres(ctx context.Context, ves ...datastore.Entity) (testcontainers.Container, *sqlx.DB, error) {
	var err error
	pgContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env:          map[string]string{"POSTGRES_PASSWORD": "password"},
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort("5432/tcp"),
			),
			Cmd: []string{"postgres", "-c", "max_connections=200"},
		},
		Started: true,
	})
	if err != nil {
		panic(err.Error())
	}

	ip, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, nil, err
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, err
	}

	fmt.Println(port.Port())

	db, err := datastore.NewPostgresDB(log, ip, port.Port(), "postgres", "password", "postgres", "disable", ves...)
	if err != nil {
		return nil, nil, err
	}

	return pgContainer, db, err
}
