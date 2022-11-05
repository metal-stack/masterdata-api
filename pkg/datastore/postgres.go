package datastore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"go.uber.org/zap"

	// import for sqlx to use postgres driver
	_ "github.com/lib/pq"
)

// exchangeable for testing
var Now = time.Now

// Storage is a interface to store objects.
type Storage[E Entity] interface {
	// generic
	Create(ctx context.Context, ve E) error
	Update(ctx context.Context, ve E) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (E, error)
	GetHistory(ctx context.Context, id string, at time.Time, ve E) error
	GetHistoryCreated(ctx context.Context, id string, ve E) error
	Find(ctx context.Context, filter map[string]any, paging *v1.Paging) ([]E, *uint64, error)
}

// Entity defines a database entity which is stored in jsonb format and with version information
type Entity interface {
	JSONField() string
	TableName() string
	Schema() string
	GetMeta() *v1.Meta
	Kind() string
	APIVersion() string
}

// datastore is the adapter to talk to the database
type datastore[E Entity] struct {
	log       *zap.Logger
	db        *sqlx.DB
	sb        squirrel.StatementBuilderType
	jsonField string
	tableName string
}

type Op string

const (
	opCreate Op = "C"
	opUpdate Op = "U"
	opDelete Op = "D"
)

// NewPostgresStorage creates a new Storage which uses postgres.
func NewPostgresDB(logger *zap.Logger, host, port, user, password, dbname, sslmode string, ves ...Entity) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	for _, ve := range ves {
		jsonField := ve.JSONField()
		logger.Info("creating schema", zap.String("entity", jsonField))
		_, err := db.Exec(ve.Schema())
		if err != nil {
			logger.Fatal("unable to create schema", zap.String("entity", jsonField), zap.Error(err))
			return nil, err
		}
	}
	return db, nil
}

// NewPostgresStorage creates a new Storage which uses postgres.
func NewPostgresStorage[E Entity](logger *zap.Logger, db *sqlx.DB, e E) (Storage[E], error) {
	ds := &datastore[E]{
		log:       logger,
		db:        db,
		sb:        squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
		jsonField: e.JSONField(),
		tableName: e.TableName(),
	}
	return ds, nil
}

// Create a entity
func (ds *datastore[E]) Create(ctx context.Context, ve E) error {
	ds.log.Debug("create", zap.String("entity", ds.jsonField), zap.Any("id", ve))
	meta := ve.GetMeta()
	if meta == nil {
		return fmt.Errorf("create of type:%s failed, meta is nil", ds.jsonField)
	}
	id := meta.GetId()
	if id == "" {
		id = uuid.NewString()
		meta.SetId(id)
	}
	kind := meta.GetKind()
	if kind == "" {
		meta.Kind = ve.Kind()
	} else if kind != ve.Kind() {
		return fmt.Errorf("create of type:%s failed, kind is set to:%s but must be:%s", ds.jsonField, kind, ve.Kind())
	}
	apiVersion := meta.GetApiversion()
	if apiVersion == "" {
		meta.Apiversion = ve.APIVersion()
	} else if apiVersion != ve.APIVersion() {
		return fmt.Errorf("create of type:%s failed, apiversion must be set to:%s", ds.jsonField, ve.APIVersion())
	}

	createdAtPb, createdAt := pbNow()
	meta.SetVersion(0)
	meta.SetCreatedTime(createdAtPb)

	q := ds.sb.Insert(
		ds.tableName,
	).SetMap(map[string]any{
		"id":         id,
		ds.jsonField: ve,
	}).Suffix(
		"RETURNING " + ds.jsonField,
	)

	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer ds.rollback(tx)

	err = q.RunWith(tx).QueryRowContext(ctx).Scan(ve)
	if err != nil {
		if IsErrorCode(err, UniqueViolationError) {
			return NewDuplicateKeyError(fmt.Sprintf("an entity of type:%s with the id:%s already exists", ds.jsonField, meta.Id))
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
func (ds *datastore[E]) Update(ctx context.Context, ve E) error {
	ds.log.Info("update", zap.String("entity", ds.jsonField))
	meta := ve.GetMeta()
	if meta == nil {
		return fmt.Errorf("update of type:%s failed, meta is nil", ds.jsonField)
	}
	id := meta.GetId()
	if id == "" {
		return fmt.Errorf("entity of type:%s has no id, cannot update: %v", ds.jsonField, ve)
	}
	kind := meta.GetKind()
	if kind == "" {
		meta.Kind = ve.Kind()
	} else if kind != ve.Kind() {
		return fmt.Errorf("update of type:%s failed, kind is set to:%s but must be:%s", ds.jsonField, kind, ve.Kind())
	}
	apiVersion := meta.GetApiversion()
	if apiVersion == "" {
		meta.Apiversion = ve.APIVersion()
	} else if apiVersion != ve.APIVersion() {
		return fmt.Errorf("update of type:%s failed, apiversion must be set to:%s", ds.jsonField, ve.APIVersion())
	}

	existingVE, err := ds.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("update - no entity of type:%s with id:%s found", ds.jsonField, id)
	}

	if ve.GetMeta().GetVersion() < existingVE.GetMeta().GetVersion() {
		return NewOptimisticLockError(
			fmt.Sprintf("optimistic lock error updating %s with id %s, existing version %d mismatches entity version %d",
				ds.jsonField, id, existingVE.GetMeta().GetVersion(), ve.GetMeta().GetVersion(),
			),
		)
	}

	pbNow, now := pbNow()

	ve.GetMeta().SetVersion(ve.GetMeta().GetVersion() + 1)
	ve.GetMeta().SetUpdatedTime(pbNow)

	// handle non updatable fields like created_time
	// simple strategy: copy unmodifiable fields from existing before update
	ve.GetMeta().SetCreatedTime(existingVE.GetMeta().GetCreatedTime())

	q := ds.sb.Update(ds.tableName).
		SetMap(map[string]any{
			ds.jsonField: ve,
		}).
		Where(squirrel.Eq{
			"id": id,
		}).
		Suffix(
			"RETURNING " + ds.jsonField,
		)

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
func (ds *datastore[E]) Get(ctx context.Context, id string) (E, error) {
	ds.log.Debug("create", zap.String("entity", ds.jsonField), zap.String("id", id))
	var zero E
	q := ds.sb.Select(
		ds.jsonField,
	).From(
		ds.tableName,
	).Where(squirrel.Eq{
		"id": id,
	})

	row := q.QueryRowContext(ctx)
	e := new(E)
	err := row.Scan(e)
	if err != nil {
		return zero, NewNotFoundError(fmt.Sprintf("entity of type:%s with id:%s not found %v", ds.jsonField, id, err))
	}
	return *e, nil
}

// Delete deletes the entity
func (ds *datastore[E]) Delete(ctx context.Context, id string) error {
	ds.log.Debug("delete", zap.String("entity", ds.jsonField), zap.String("id", id))
	ve, err := ds.Get(ctx, id)
	if err != nil {
		return err
	}

	// delete dataset in table
	q := ds.sb.Delete(ds.tableName).
		Where(squirrel.Eq{"id": id})
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
		return NewDataCorruptionError(fmt.Sprintf("data corruption: delete of %s with id %s affected %d rows", ds.jsonField, id, rowsAffected))
	}
	if rowsAffected < 1 {
		return NewNotFoundError(fmt.Sprintf("not found: delete of %s with id %s affected %d rows", ds.jsonField, id, rowsAffected))
	}

	// insert dataset in history table
	err = ds.insertHistory(ve, opDelete, Now(), tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Find returns matching elements from the database
func (ds *datastore[E]) Find(ctx context.Context, filter map[string]any, paging *v1.Paging) ([]E, *uint64, error) {
	ds.log.Debug("find", zap.String("entity", ds.jsonField), zap.Any("filter", filter))
	q := ds.sb.Select(ds.jsonField).
		From(ds.tableName)

	if len(filter) > 0 {
		q = q.Where(filter)
	}
	q = q.OrderBy("id")

	// Add paging query if paging is defined
	q, nextPage := addPaging(q, paging)

	rows, err := q.QueryContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	ves := []E{}
	for rows.Next() {
		e := new(E)
		err = rows.Scan(e)
		if err != nil {
			return nil, nil, err
		}
		ves = append(ves, *e)
	}

	err = rows.Err()
	if err != nil {
		return nil, nil, err
	}
	if paging != nil && paging.Count != nil && uint64(len(ves)) == *paging.Count {
		return ves, nextPage, err
	}
	return ves, nil, nil
}

// Get the history entity for given id and latest before or equal the given point in time
// returns NotFoundError if no entity can be found
func (ds *datastore[E]) GetHistory(ctx context.Context, id string, at time.Time, ve E) error {
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
func (ds *datastore[E]) GetHistoryCreated(ctx context.Context, id string, ve E) error {
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
func (ds *datastore[E]) getHistoryWithPredicate(ctx context.Context, pred any, ve E) error {
	tablename := historyTablename(ds.tableName)
	q := ds.sb.Select(ds.jsonField).From(tablename).Where(pred).OrderByClause("created_at DESC").Limit(1)

	sql, _, _ := q.ToSql()
	ds.log.Info("get", zap.String("entity", ds.jsonField), zap.String("sql", sql), zap.Any("predicate", pred))
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
	return NewNotFoundError(fmt.Sprintf("entity of type:%s with predicate:%s not found", ds.jsonField, pred))
}

// insertHistory inserts the given entity in the history table of the entity using the runner, which may be a Tx.
func (ds *datastore[E]) insertHistory(ve E, op Op, createdAt time.Time, runner squirrel.BaseRunner) error {
	qh := ds.sb.Insert(historyTablename(ds.tableName)).
		SetMap(map[string]any{
			"id":         ve.GetMeta().Id,
			"op":         op,
			"created_at": createdAt,
			ds.jsonField: ve,
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

// pbNow returns the current time as Protobuf and time
func pbNow() (*timestamppb.Timestamp, time.Time) {
	now := Now()
	nowPb := timestamppb.New(now)
	return nowPb, now
}

// rollback tries to rollback the given transaction and logs an eventual rollback error
func (ds *datastore[E]) rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		ds.log.Error("error rolling back", zap.Error(err))
	}
}

const defaultPagingLimit = uint64(100)

func addPaging(q squirrel.SelectBuilder, paging *v1.Paging) (squirrel.SelectBuilder, *uint64) {
	if paging == nil {
		return q, nil
	}

	limit := defaultPagingLimit
	if paging.Count != nil {
		limit = *paging.Count
	}
	offset := uint64(0)
	nextpage := uint64(1)
	if paging.Page != nil {
		offset = *paging.Page * limit
		nextpage = *paging.Page + 1
	}
	q = q.Limit(limit).Offset(offset)
	return q, &nextpage
}
