package datastore

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/lib/pq"
	"reflect"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	// import for sqlx to use postgres driver
	_ "github.com/lib/pq"
)

// Storage is a interface to store objects.
type Storage interface {
	// generic
	Create(ctx context.Context, ve VersionedJSONEntity) error
	Update(ctx context.Context, ve VersionedJSONEntity) error
	Delete(ctx context.Context, ve VersionedJSONEntity) error
	Get(ctx context.Context, id string, ve VersionedJSONEntity) error
	Find(ctx context.Context, filter map[string]interface{}, result interface{}) error
}

// JSONEntity is storable in json format
type JSONEntity interface {
	JSONField() string
	TableName() string
	Schema() string
}

// VersionedEntity defines a database entity which is stored with version information
type VersionedEntity interface {
	GetMeta() *v1.Meta
}

// VersionedJSONEntity defines a database entity which is stored in jsonb format and with version information
type VersionedJSONEntity interface {
	JSONEntity
	VersionedEntity
}

// Datastore is the adapter to talk to the database
type Datastore struct {
	log   *zap.Logger
	db    *sqlx.DB
	sb    squirrel.StatementBuilderType
	types map[string]VersionedJSONEntity
}

// NewPostgresStorage creates a new Storage which uses postgres.
func NewPostgresStorage(logger *zap.Logger, host, port, user, password, dbname, sslmode string, ves ...VersionedJSONEntity) (*Datastore, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	types := make(map[string]VersionedJSONEntity)
	for _, ve := range ves {
		jsonField := ve.JSONField()
		logger.Sugar().Infow("creating schema", "entity", jsonField)
		_, err = db.Exec(ve.Schema())
		if err != nil {
			logger.Sugar().Fatalw("unable to create schema", "entity", jsonField, "err", err)
			return nil, err
		}
		types[jsonField] = ve
	}
	ds := &Datastore{
		log:   logger,
		db:    db,
		sb:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
		types: types,
	}
	return ds, nil
}

// Create a entity
func (ds *Datastore) Create(ctx context.Context, ve VersionedJSONEntity) error {
	jsonField := ve.JSONField()
	tableName := ve.TableName()
	_, ok := ds.types[jsonField]
	if !ok {
		return fmt.Errorf("type:%s is not registered", jsonField)
	}
	meta := ve.GetMeta()
	if meta == nil {
		return fmt.Errorf("create of type:%s failed, meta is nil", jsonField)
	}
	id := meta.GetId()
	if id == "" {
		id = uuid.Must(uuid.NewRandom()).String()
		meta.SetId(id)
	}

	meta.SetVersion(0)
	meta.SetCreatedTime(ptypes.TimestampNow())

	q := ds.sb.Insert(
		tableName,
	).SetMap(map[string]interface{}{
		"id":      id,
		jsonField: ve,
	}).Suffix(
		"RETURNING " + jsonField,
	)
	sql, vals, _ := q.ToSql()
	ds.log.Sugar().Infow("create", "entity", tableName, "sql", sql, "values", vals)

	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = q.QueryRowContext(ctx).Scan(ve)
	if err != nil {
		switch pqe := err.(type) {
		case *pq.Error:
			if pqe.Code == "23505" {
				return NewDuplicateKeyError(fmt.Sprintf("an entity of type:%s with the id:%s already exists", jsonField, meta.Id))
			}
		}
		return err
	}

	err = tx.Commit()

	return err
}

// Update the entity
func (ds *Datastore) Update(ctx context.Context, ve VersionedJSONEntity) error {
	jsonField := ve.JSONField()
	tableName := ve.TableName()
	_, ok := ds.types[jsonField]
	if !ok {
		return fmt.Errorf("type:%s is not registered", jsonField)
	}
	meta := ve.GetMeta()
	if meta == nil {
		return fmt.Errorf("update of type:%s failed, meta is nil", jsonField)
	}
	id := meta.GetId()
	if id == "" {
		return fmt.Errorf("entity of type:%s has no id, cannot update: %v", jsonField, ve)
	}

	elemt := reflect.TypeOf(ve).Elem()
	existingVE, ok := reflect.New(elemt).Interface().(VersionedJSONEntity)
	if !ok {
		return fmt.Errorf("entity is not a VersionedJSONEntity: %v", ve)
	}

	err := ds.Get(ctx, id, existingVE)
	if err != nil {
		return errors.Errorf("update - no entity of type:%s with id:%s found", jsonField, id)
	}

	if ve.GetMeta().GetVersion() < existingVE.GetMeta().GetVersion() {
		return NewOptimisticLockError(
			fmt.Sprintf("optimistic lock error updating %s with id %s, existing version %d mismatches entity version %d",
				jsonField, id, existingVE.GetMeta().GetVersion(), ve.GetMeta().GetVersion(),
			),
		)
	}

	ve.GetMeta().SetVersion(ve.GetMeta().GetVersion() + 1)
	ve.GetMeta().SetUpdatedTime(ptypes.TimestampNow())

	// FIXME how to handle non updateable fields like created_time
	// simple strategy copy unmodifiable fields from existing before update
	ve.GetMeta().SetCreatedTime(existingVE.GetMeta().GetCreatedTime())

	q := ds.sb.Update(tableName).
		SetMap(map[string]interface{}{
			jsonField: ve,
		}).
		Where(squirrel.Eq{
			"id": id,
		}).
		Suffix(
			"RETURNING " + jsonField,
		)
	sql, vals, _ := q.ToSql()
	ds.log.Sugar().Infow("update", "entity", tableName, "sql", sql, "values", vals)

	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = q.QueryRowContext(ctx).Scan(ve)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// Get the entity for given id
// returns NotFoundError if no entity can be found
func (ds *Datastore) Get(ctx context.Context, id string, ve VersionedJSONEntity) error {
	jsonField := ve.JSONField()
	tableName := ve.TableName()
	_, ok := ds.types[jsonField]
	if !ok {
		return fmt.Errorf("type:%s is not registered", jsonField)
	}
	q := ds.sb.Select(
		jsonField,
	).From(
		tableName,
	).Where(squirrel.Eq{
		"id": id,
	})

	sql, _, _ := q.ToSql()
	ds.log.Sugar().Infow("get", "entity", jsonField, "sql", sql, "id", id)
	rows, err := q.QueryContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		cerr := rows.Close()
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	if rows.Next() {
		return rows.Scan(ve)
	}
	// we have no row
	return NewNotFoundError(fmt.Sprintf("entity of type:%s with id:%s not found", jsonField, id))
}

// Delete deletes the entity
func (ds *Datastore) Delete(ctx context.Context, ve VersionedJSONEntity) error {
	jsonField := ve.JSONField()
	tableName := ve.TableName()
	_, ok := ds.types[jsonField]
	if !ok {
		return fmt.Errorf("type:%s is not registered", jsonField)
	}
	q := ds.sb.Delete(tableName).
		Where(squirrel.Eq{"id": ve.GetMeta().GetId()})

	sql, _, _ := q.ToSql()
	ds.log.Sugar().Debugw("delete", "entity", jsonField, "sql", sql)
	result, err := q.ExecContext(ctx)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected > 1 {
		return NewDataCorruptionError(fmt.Sprintf("datacorruption: delete of %s with id %s affected %d rows", jsonField, ve.GetMeta().Id, rowsAffected))
	}
	if rowsAffected < 1 {
		return NewNotFoundError(fmt.Sprintf("not found: delete of %s with id %s affected %d rows", jsonField, ve.GetMeta().Id, rowsAffected))
	}

	return err
}

// Find returns matching elements from the database
func (ds *Datastore) Find(ctx context.Context, filter map[string]interface{}, result interface{}) error {
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("result argument must be a slice address")
	}

	slicev := resultv.Elem()
	elemt := slicev.Type().Elem()

	ve, ok := reflect.New(elemt).Interface().(VersionedJSONEntity)
	if !ok {
		return fmt.Errorf("result slice element type must implement VersionedJSONEntity-Interface")
	}
	jsonField := ve.JSONField()
	tableName := ve.TableName()
	_, ok = ds.types[jsonField]
	if !ok {
		return fmt.Errorf("type:%s is not registered", jsonField)
	}

	q := ds.sb.Select(jsonField).
		From(tableName)

	if len(filter) > 0 {
		q = q.Where(filter)
	}
	q = q.OrderBy("id")

	sql, vals, _ := q.ToSql()
	ds.log.Sugar().Debugw("find", "sql", sql, "values", vals)

	rows, err := q.QueryContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		cerr := rows.Close()
		if err == nil && cerr != nil {
			err = cerr
		}
	}()

	for rows.Next() {
		elemp := reflect.New(elemt)
		err = rows.Scan(elemp.Interface())
		if err != nil {
			return err
		}
		slicev = reflect.Append(slicev, elemp.Elem())
	}
	resultv.Elem().Set(slicev)

	err = rows.Err()
	if err != nil {
		return err
	}

	return err
}
