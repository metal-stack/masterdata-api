package datastore

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
)

type memoryDatastore[E Entity] struct {
	lock     sync.RWMutex
	entities map[string]map[string]E
	table    string
	log      *slog.Logger
}

func NewMemory[E Entity](log *slog.Logger, e E) Storage[E] {
	table := e.TableName()
	return &memoryDatastore[E]{
		lock:     sync.RWMutex{},
		entities: map[string]map[string]E{},
		table:    table,
		log:      log,
	}
}

// Create implements Storage.
func (m *memoryDatastore[E]) Create(ctx context.Context, ve E) error {
	m.log.Debug("create", "entity", m.table, "value", ve)
	meta := ve.GetMeta()
	if meta == nil {
		return fmt.Errorf("create of type:%s failed, meta is nil", m.table)
	}

	id := ve.GetMeta().Id
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.entities[m.table][id]
	if ok {
		return NewDuplicateKeyError(fmt.Sprintf("an entity of type:%s with the id:%s already exists", m.table, id))
	}
	m.entities[m.table][id] = ve
	return nil
}

// Delete implements Storage.
func (m *memoryDatastore[E]) Delete(ctx context.Context, id string) error {
	m.log.Debug("delete", "entity", m.table, "id", id)

	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.entities[m.table][id]
	if !ok {
		return NewNotFoundError(fmt.Sprintf("not found: delete of %s with id %s", m.table, id))
	}
	delete(m.entities[m.table], id)

	return nil
}

// DeleteAll implements Storage.
func (m *memoryDatastore[E]) DeleteAll(ctx context.Context, ids ...string) error {
	panic("unimplemented")
}

// Find implements Storage.
func (m *memoryDatastore[E]) Find(ctx context.Context, filter map[string]any, paging *v1.Paging) ([]E, *uint64, error) {
	m.log.Debug("find", "entity", m.table, "filter", filter)

	m.lock.Lock()
	defer m.lock.Unlock()

	var result []E
	for _, e := range m.entities[m.table] {
		// FIXME implement filtering
		result = append(result, e)
	}

	return result, nil, nil
}

// Get implements Storage.
func (m *memoryDatastore[E]) Get(ctx context.Context, id string) (E, error) {
	m.log.Debug("get", "entity", m.table, "id", id)
	var zero E
	m.lock.Lock()
	defer m.lock.Unlock()

	e, ok := m.entities[m.table][id]
	if !ok {
		return zero, NewNotFoundError(fmt.Sprintf("not found: delete of %s with id %s", m.table, id))
	}

	return e, nil
}

// GetHistory implements Storage.
func (m *memoryDatastore[E]) GetHistory(ctx context.Context, id string, at time.Time, ve E) error {
	panic("unimplemented")
}

// GetHistoryCreated implements Storage.
func (m *memoryDatastore[E]) GetHistoryCreated(ctx context.Context, id string, ve E) error {
	panic("unimplemented")
}

// Update implements Storage.
func (m *memoryDatastore[E]) Update(ctx context.Context, ve E) error {
	m.log.Debug("update", "entity", m.table)
	meta := ve.GetMeta()
	if meta == nil {
		return fmt.Errorf("update of type:%s failed, meta is nil", m.table)
	}
	id := ve.GetMeta().Id
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.entities[m.table][id]
	if !ok {
		return NewNotFoundError(fmt.Sprintf("not found: delete of %s with id %s", m.table, id))
	}

	m.entities[m.table][id] = ve
	return nil
}
