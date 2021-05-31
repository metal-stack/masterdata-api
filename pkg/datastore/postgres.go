package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	// import for sqlx to use postgres driver
	_ "github.com/lib/pq"
)

// exchangable for testing
var Now = time.Now

// Storage is a interface to store objects.
type Storage interface {
	// generic
	Create(ctx context.Context, ve VersionedJSONEntity) error
	Update(ctx context.Context, ve VersionedJSONEntity) error
	Delete(ctx context.Context, ve VersionedJSONEntity) error
	Get(ctx context.Context, id string, ve VersionedJSONEntity) error
	GetHistory(ctx context.Context, id string, at time.Time, ve VersionedJSONEntity) error
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
	Kind() string
	APIVersion() string
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

type Op string

const (
	opCreate Op = "C"
	opUpdate Op = "U"
	opDelete Op = "D"
)

// NewPostgresStorage creates a new Storage which uses postgres.
func NewPostgresStorage(logger *zap.Logger, host, port, user, password, dbname, sslmode string, ves ...VersionedJSONEntity) (*Datastore, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	types := make(map[string]VersionedJSONEntity)
	for _, ve := range ves {
		jsonField := ve.JSONField()
		logger.Info("creating schema", zap.String("entity", jsonField))
		_, err = db.Exec(ve.Schema())
		if err != nil {
			logger.Fatal("unable to create schema", zap.String("entity", jsonField), zap.Error(err))
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
	kind := meta.GetKind()
	if kind == "" {
		meta.Kind = ve.Kind()
	} else if kind != ve.Kind() {
		return fmt.Errorf("create of type:%s failed, kind is set to:%s but must be:%s", jsonField, kind, ve.Kind())
	}
	apiVersion := meta.GetApiversion()
	if apiVersion == "" {
		meta.Apiversion = ve.APIVersion()
	} else if apiVersion != ve.APIVersion() {
		return fmt.Errorf("create of type:%s failed, apiversion must be set to:%s", jsonField, ve.APIVersion())
	}

	createdAtPb, createdAt := PbNow()
	meta.SetVersion(0)
	meta.SetCreatedTime(createdAtPb)

	q := ds.sb.Insert(
		tableName,
	).SetMap(map[string]interface{}{
		"id":      id,
		jsonField: ve,
	}).Suffix(
		"RETURNING " + jsonField,
	)
	sql, vals, _ := q.ToSql()
	ds.log.Info("create", zap.String("entity", tableName), zap.String("sql", sql), zap.Any("values", vals))

	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer ds.rollback(tx)

	err = q.RunWith(tx).QueryRowContext(ctx).Scan(ve)
	if err != nil {
		if IsErrorCode(err, UniqueViolationError) {
			return NewDuplicateKeyError(fmt.Sprintf("an entity of type:%s with the id:%s already exists", jsonField, meta.Id))
		}
		return err
	}
	err = ds.insertHistory(ve, opCreate, createdAt, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
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
	kind := meta.GetKind()
	if kind == "" {
		meta.Kind = ve.Kind()
	} else if kind != ve.Kind() {
		return fmt.Errorf("update of type:%s failed, kind is set to:%s but must be:%s", jsonField, kind, ve.Kind())
	}
	apiVersion := meta.GetApiversion()
	if apiVersion == "" {
		meta.Apiversion = ve.APIVersion()
	} else if apiVersion != ve.APIVersion() {
		return fmt.Errorf("update of type:%s failed, apiversion must be set to:%s", jsonField, ve.APIVersion())
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

	pbNow, now := PbNow()

	ve.GetMeta().SetVersion(ve.GetMeta().GetVersion() + 1)
	ve.GetMeta().SetUpdatedTime(pbNow)

	// handle non updateable fields like created_time
	// simple strategy: copy unmodifiable fields from existing before update
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
	ds.log.Info("update", zap.String("entity", tableName), zap.String("sql", sql), zap.Any("values", vals))

	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer ds.rollback(tx)

	err = q.RunWith(tx).QueryRowContext(ctx).Scan(ve)
	if err != nil {
		return err
	}

	// insert dataset in history table
	err = ds.insertHistory(ve, opUpdate, now, tx)
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
	ds.log.Info("get", zap.String("entity", jsonField), zap.String("sql", sql), zap.String("id", id))
	rows, err := q.QueryContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
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

	elemt := reflect.TypeOf(ve).Elem()
	existingVE, ok := reflect.New(elemt).Interface().(VersionedJSONEntity)
	if !ok {
		return fmt.Errorf("entity is not a VersionedJSONEntity: %v", ve)
	}
	err := ds.Get(ctx, ve.GetMeta().Id, existingVE)
	if err != nil {
		return err
	}

	// delete dataset in table
	q := ds.sb.Delete(tableName).
		Where(squirrel.Eq{"id": ve.GetMeta().GetId()})
	sql, _, err := q.ToSql()
	if err != nil {
		return err
	}
	ds.log.Debug("delete", zap.String("entity", jsonField), zap.String("sql", sql))

	// in tx
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer ds.rollback(tx)

	result, err := q.RunWith(tx).ExecContext(ctx)
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

	// insert dataset in history table
	err = ds.insertHistory(existingVE, opDelete, Now(), tx)
	if err != nil {
		return err
	}

	return tx.Commit()
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
	ds.log.Debug("find", zap.String("sql", sql), zap.Any("values", vals))

	rows, err := q.QueryContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
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

// Get the history entity for given id and latest before or equal the given point in time
// returns NotFoundError if no entity can be found
func (ds *Datastore) GetHistory(ctx context.Context, id string, at time.Time, ve VersionedJSONEntity) error {
	return ds.getHistoryWithPredicate(ctx, squirrel.And{
		squirrel.Eq{
			"id": id,
		},
		squirrel.LtOrEq{
			"created_at": at,
		},
	}, ve)
}

// Get the first history entity for given id, returns NotFoundError if no entity can be found
func (ds *Datastore) GetHistoryCreated(ctx context.Context, id string, ve VersionedJSONEntity) error {
	return ds.getHistoryWithPredicate(ctx, squirrel.And{
		squirrel.Eq{
			"id": id,
		},
		squirrel.Eq{
			"op": opCreate,
		},
	}, ve)
}

// Get the top matching history entity for given filter criteria,
// returns NotFoundError if no entity can be found
func (ds *Datastore) getHistoryWithPredicate(ctx context.Context, pred interface{}, ve VersionedJSONEntity) error {
	jsonField := ve.JSONField()
	tableName := historyTablename(ve.TableName())
	_, ok := ds.types[jsonField]
	if !ok {
		return fmt.Errorf("type:%s is not registered", jsonField)
	}
	q := ds.sb.Select(jsonField).From(tableName).Where(pred).OrderByClause("created_at DESC").Limit(1)

	sql, _, _ := q.ToSql()
	ds.log.Info("get", zap.String("entity", jsonField), zap.String("sql", sql), zap.Any("predicate", pred))
	rows, err := q.QueryContext(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()
	if rows.Next() {
		return rows.Scan(ve)
	}
	// we have no row
	return NewNotFoundError(fmt.Sprintf("entity of type:%s with predicate:%s not found", jsonField, pred))
}

// insertHistory inserts the given entity in the history table of the entity using the runner, which may be a Tx.
func (ds *Datastore) insertHistory(ve VersionedJSONEntity, op Op, createdAt time.Time, runner squirrel.BaseRunner) error {
	jsonField := ve.JSONField()
	tableName := ve.TableName()
	qh := ds.sb.Insert(historyTablename(tableName)).
		SetMap(map[string]interface{}{
			"id":         ve.GetMeta().Id,
			"op":         op,
			"created_at": createdAt,
			jsonField:    ve,
		})
	_, err := qh.RunWith(runner).Exec()
	if err != nil {
		return err
	}
	return nil
}

// historyTablename returns the tablename of the table-twin with historic data.
func historyTablename(table string) string {
	return fmt.Sprintf("%s_history", table)
}

// PbNow returns the current time as Protobuf and time
func PbNow() (*timestamppb.Timestamp, time.Time) {
	now := Now()
	nowPb := timestamppb.New(now)
	return nowPb, now
}

// rollback tries to rollback the given transaction and logs an eventual rollback error
func (ds *Datastore) rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		ds.log.Error("error rolling back", zap.Error(err))
	}
}
