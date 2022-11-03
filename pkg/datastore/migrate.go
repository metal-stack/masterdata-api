package datastore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/lopezator/migrator"
	healthv1 "github.com/metal-stack/masterdata-api/api/grpc/health/v1"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/health"
	"go.uber.org/zap"
)

// MigrateDB applies necessary DB Migrations.
func (ds *Datastore) MigrateDB(healthServer *health.Server) error {

	m, err := migrator.New(
		migrator.Migrations(
			// migrations will be applied and remembered in this order, so always add new migrations below if necessary
			&migrator.Migration{
				Name: "Consolidate History for Tenant and Project",
				Func: func(tx *sql.Tx) error {
					entities := []any{&[]v1.Project{}, &[]v1.Tenant{}}
					for _, e := range entities {
						err := ds.consolidateHistory(tx, e)
						if err != nil {
							ds.log.Error("error consolidate history", zap.Error(err))
							return err
						}
					}
					return nil
				},
			},
		),
		migrator.WithLogger(migrator.LoggerFunc(func(msg string, args ...any) {
			ds.log.Sugar().Infof(msg, args...)
		})),
	)
	if err != nil {
		return err
	}

	// Migrate up
	if err := m.Migrate(ds.db.DB); err != nil {
		return err
	}

	healthServer.SetServingStatus("migratedb", healthv1.HealthCheckResponse_SERVING)
	return nil
}

// consolidateHistory ensures, that for each VersionedJSONEntity there is at least one "created"-row in the history table.
// The type of entities to consolidate is specified by the given pointer to a slice of entities.
func (ds *Datastore) consolidateHistory(tx *sql.Tx, entitySlicePtr any) error {

	entitySliceV := reflect.ValueOf(entitySlicePtr)
	if entitySliceV.Kind() != reflect.Ptr || entitySliceV.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("entity argument must be a slice address")
	}
	entitySliceElem := entitySliceV.Elem()
	entitySliceElementType := entitySliceElem.Type().Elem()

	filter := make(map[string]any)
	_, err := ds.Find(context.Background(), filter, nil, entitySlicePtr)
	if err != nil {
		return err
	}

	for i := 0; i < entitySliceElem.Len(); i++ {
		vpi := entitySliceElem.Index(i).Addr().Interface()
		enityVe, ok := vpi.(Entity)
		if !ok {
			return fmt.Errorf("element type must implement VersionedJSONEntity-Interface, was %T", vpi)
		}

		historyVe, ok := reflect.New(entitySliceElementType).Interface().(Entity)
		if !ok {
			return fmt.Errorf("element type must implement VersionedJSONEntity-Interface")
		}

		// check if we already have a "created" row for this entity in history
		err = ds.GetHistoryCreated(context.Background(), enityVe.GetMeta().Id, historyVe)
		if err == nil {
			continue
		}

		if !errors.As(err, &NotFoundError{}) {
			return err // some sort of technical error stops us
		}

		// consolidate history by inserting the "created" row in history at the correct date
		entityCreatedPbTs := enityVe.GetMeta().CreatedTime
		entityCreated := entityCreatedPbTs.AsTime()
		err := ds.insertHistory(enityVe, opCreate, entityCreated, tx)
		if err != nil {
			return err
		}
	}
	return nil
}
