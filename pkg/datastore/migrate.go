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
			&migrator.Migration{
				Name: "Migrate quota field",
				Func: func(tx *sql.Tx) error {
					ts := New(log, db, &v1.Tenant{})
					ctx := context.Background()

					tenants, _, err := ts.Find(ctx, nil, nil)
					if err != nil {
						return err
					}
					for _, tenant := range tenants {
						if tenant.Quotas == nil {
							continue
						}

						for _, q := range []*v1.Quota{
							tenant.Quotas.Cluster,
							tenant.Quotas.Ip,
							tenant.Quotas.Machine,
							tenant.Quotas.Project,
						} {
							if q == nil {
								continue
							}

							if q.Quota != nil { // nolint:staticcheck
								q.Max = &q.Quota.Value // nolint:staticcheck
								log.Info("migrating deprecated quota field of tenant", "tenant", tenant.Meta.Id, "value", *q.Max)
							}
						}

						err := ts.Update(ctx, tenant)
						if err != nil {
							return err
						}
					}

					return nil
				},
			},
			&migrator.Migration{
				Name: "Unique membership indexes",
				Func: func(tx *sql.Tx) error {
					var (
						ctx = context.Background()
						pms = New(log, db, &v1.ProjectMember{})
						tms = New(log, db, &v1.TenantMember{})
					)

					projectMembers, _, err := pms.Find(ctx, nil)
					if err != nil {
						return err
					}

					seen := map[string]bool{}

					for _, m := range projectMembers {
						key := fmt.Sprintf("%s|%s|%s", m.GetProjectId(), m.GetTenantId(), m.GetNamespace())
						if seen[key] {
							if err := pms.Delete(ctx, m.Meta.GetId()); err != nil {
								return err
							}
							continue
						}
						seen[key] = true
					}

					tenantMembers, _, err := tms.Find(ctx, nil)
					if err != nil {
						return err
					}

					seen = map[string]bool{}

					for _, m := range tenantMembers {
						key := fmt.Sprintf("%s|%s|%s", m.GetTenantId(), m.GetMemberId(), m.GetNamespace())
						if seen[key] {
							if err := tms.Delete(ctx, m.Meta.GetId()); err != nil {
								return err
							}
							continue
						}
						seen[key] = true
					}

					_, err = tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS projectmembers_membership_idx ON projectmembers ((projectmember ->> 'project_id'), (projectmember ->> 'tenant_id'), (COALESCE(projectmember ->> 'namespace', '')))`)
					if err != nil {
						return err
					}

					_, err = tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS tenantmembers_membership_idx ON tenantmembers ((tenantmember ->> 'tenant_id'), (tenantmember ->> 'member_id'), (COALESCE(tenantmember ->> 'namespace', '')))`)
					if err != nil {
						return err
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
