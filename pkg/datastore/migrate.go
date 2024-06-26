package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/lopezator/migrator"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

// MigrateDB applies necessary DB Migrations.
func MigrateDB(log *slog.Logger, db *sqlx.DB, healthServer *health.Server) error {

	m, err := migrator.New(
		migrator.Migrations(
			// migrations will be applied and remembered in this order, so always add new migrations below if necessary
			&migrator.Migration{
				Name: "Sample Migration for Tenant",
				Func: func(tx *sql.Tx) error {
					ts := New(log, db, &v1.Tenant{})

					tenants, _, err := ts.Find(context.Background(), nil, nil)
					if err != nil {
						return err
					}
					for _, tenant := range tenants {
						log.Debug("migrate", "tenant", tenant)
					}

					return nil
				},
			},
		),
		migrator.WithLogger(migrator.LoggerFunc(func(msg string, args ...any) {
			log.Info(fmt.Sprintf(msg, args...))
		})),
	)
	if err != nil {
		return err
	}

	// Migrate up
	if err := m.Migrate(db.DB); err != nil {
		return err
	}

	healthServer.SetServingStatus("migratedb", healthv1.HealthCheckResponse_SERVING)
	return nil
}
