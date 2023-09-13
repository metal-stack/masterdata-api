package datastore

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lopezator/migrator"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/health"
	"go.uber.org/zap"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

// MigrateDB applies necessary DB Migrations.
func MigrateDB(log *zap.SugaredLogger, db *sqlx.DB, healthServer *health.Server) error {

	m, err := migrator.New(
		migrator.Migrations(
			// migrations will be applied and remembered in this order, so always add new migrations below if necessary
			&migrator.Migration{
				Name: "Sample Migration for Tenant",
				Func: func(tx *sql.Tx) error {
					ts, err := New(log.Desugar(), db, &v1.Tenant{})
					if err != nil {
						return err
					}

					tenants, _, err := ts.Find(context.Background(), nil, nil)
					if err != nil {
						return err
					}
					for _, tenant := range tenants {
						log.Debugw("migrate", "tenant", tenant)
					}

					return nil
				},
			},
		),
		migrator.WithLogger(migrator.LoggerFunc(func(msg string, args ...any) {
			log.Infof(msg, args...)
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
