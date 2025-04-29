package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	apiv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"github.com/metal-stack/v"
	cli "github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:    "masterdata server",
		Usage:   "grpc server for masterdata",
		Version: v.V.String(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "grpc-server-endpoint",
				Value:   ":9090",
				Usage:   "gRPC server endpoint",
				EnvVars: []string{"MASTERDATA_GRPC_SERVER_ENDPOINT"},
			},
			&cli.StringFlag{
				Name:    "metrics-endpoint",
				Value:   ":2112",
				Usage:   "metrics endpoint",
				EnvVars: []string{"MASTERDATA_METRICS_ENDPOINT"},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "info",
				Usage:   "log-level can be one of error|warn|info|debug",
				EnvVars: []string{"MASTERDATA_LOG_LEVEL"},
			},
		},
		Commands: []*cli.Command{
			// {
			// 	Name:    "memory",
			// 	Aliases: []string{"m"},
			// 	Usage:   "start with memory backend",
			// 	Action: func(ctx *cli.Context) error {
			// 		c := getConfig(ctx)
			// 		c.Storage = goipam.NewMemory(ctx.Context)
			// 		s := newServer(c)
			// 		return s.Run()
			// 	},
			// },
			{
				Name:    "postgres",
				Aliases: []string{"pg"},
				Usage:   "start with postgres backend",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "host",
						Value:   "localhost",
						Usage:   "postgres db hostname",
						EnvVars: []string{"MASTERDATA_PG_HOST"},
					},
					&cli.StringFlag{
						Name:    "port",
						Value:   "5432",
						Usage:   "postgres db port",
						EnvVars: []string{"MASTERDATA_PG_PORT"},
					},
					&cli.StringFlag{
						Name:    "user",
						Value:   "go-ipam",
						Usage:   "postgres db user",
						EnvVars: []string{"MASTERDATA_PG_USER"},
					},
					&cli.StringFlag{
						Name:    "password",
						Value:   "secret",
						Usage:   "postgres db password",
						EnvVars: []string{"MASTERDATA_PG_PASSWORD"},
					},
					&cli.StringFlag{
						Name:    "dbname",
						Value:   "goipam",
						Usage:   "postgres db name",
						EnvVars: []string{"MASTERDATA_PG_DBNAME"},
					},
					&cli.StringFlag{
						Name:    "sslmode",
						Value:   "disable",
						Usage:   "postgres sslmode, possible values: disable|require|verify-ca|verify-full",
						EnvVars: []string{"MASTERDATA_PG_SSLMODE"},
					},
				},
				Action: func(ctx *cli.Context) error {
					c := getConfig(ctx)
					host := ctx.String("host")
					port := ctx.String("port")
					user := ctx.String("user")
					password := ctx.String("password")
					dbname := ctx.String("dbname")
					sslmode := ctx.String("sslmode")

					ves := []datastore.Entity{
						&apiv1.Project{},
						&apiv1.ProjectMember{},
						&apiv1.Tenant{},
						&apiv1.TenantMember{},
					}

					db, err := datastore.NewPostgresDB(c.Log, host, port, user, password, dbname, sslmode, ves...)
					if err != nil {
						return fmt.Errorf("failed to create postgres connection: %w", err)
					}
					ps := datastore.New(c.Log, db, &apiv1.Project{})
					pms := datastore.New(c.Log, db, &apiv1.ProjectMember{})
					ts := datastore.New(c.Log, db, &apiv1.Tenant{})
					tms := datastore.New(c.Log, db, &apiv1.TenantMember{})
					c.ProjectDataStore = ps
					c.ProjectMemberDataStore = pms
					c.TenantDataStore = ts
					c.TenantMemberDataStore = tms
					c.DB = db
					s := newServer(c)
					return s.Run()
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("unable to start masterdata service: %v", err)
	}
}

func getConfig(ctx *cli.Context) config {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	switch ctx.String("log-level") {
	case "debug":
		opts.Level = slog.LevelDebug
	case "error":
		opts.Level = slog.LevelError
	}

	return config{
		GrpcServerEndpoint: ctx.String("grpc-server-endpoint"),
		MetricsEndpoint:    ctx.String("metrics-endpoint"),
		Log:                slog.New(slog.NewJSONHandler(os.Stdout, opts)),
	}
}
