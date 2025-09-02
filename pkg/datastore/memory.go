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
	entities map[string]E
	entity   string
	log      *slog.Logger
}

func NewMemory[E Entity](log *slog.Logger, e E) Storage[E] {
	entity := e.JSONField()
	return &memoryDatastore[E]{
		lock:     sync.RWMutex{},
		entities: make(map[string]E),
		entity:   entity,
		log:      log,
	}
}

// Create implements Storage.
func (m *memoryDatastore[E]) Create(ctx context.Context, ve E) error {
	m.log.Debug("create", "entity", m.entity, "value", ve)
	meta := ve.GetMeta()
	if meta == nil {
		return fmt.Errorf("create of type:%s failed, meta is nil", m.entity)
	}

	id := ve.GetMeta().Id
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.entities[id]
	if ok {
		return NewDuplicateKeyError(fmt.Sprintf("an entity of type:%s with the id:%s already exists", m.entity, id))
	}
	m.entities[id] = ve
	return nil
}

// Delete implements Storage.
func (m *memoryDatastore[E]) Delete(ctx context.Context, id string) error {
	m.log.Debug("delete", "entity", m.entity, "id", id)

	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.entities[id]
	if !ok {
		return NewNotFoundError(fmt.Sprintf("delete of %s with id %s", m.entity, id))
	}
	delete(m.entities, id)

	return nil
}

// DeleteAll implements Storage.
func (m *memoryDatastore[E]) DeleteAll(ctx context.Context, ids ...string) error {
	m.entities = make(map[string]E)
	return nil
}

// Find implements Storage.
func (m *memoryDatastore[E]) Find(ctx context.Context, paging *v1.Paging, filters ...any) ([]E, *uint64, error) {
	m.log.Debug("find", "entity", m.entity, "filter", filters)

	m.lock.Lock()
	defer m.lock.Unlock()

	var result []E
	for _, e := range m.entities {
		// FIXME implement filtering
		result = append(result, e)
	}

	return result, nil, nil
}

// Get implements Storage.
func (m *memoryDatastore[E]) Get(ctx context.Context, id string) (E, error) {
	m.log.Debug("get", "entity", m.entity, "id", id)
	var zero E
	m.lock.Lock()
	defer m.lock.Unlock()

	e, ok := m.entities[id]
	if !ok {
		return zero, NewNotFoundError(fmt.Sprintf("get of %s with id %s", m.entity, id))
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
	m.log.Debug("update", "entity", m.entity)
	meta := ve.GetMeta()
	if meta == nil {
		return fmt.Errorf("update of type:%s failed, meta is nil", m.entity)
	}
	id := ve.GetMeta().Id
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.entities[id]
	if !ok {
		return NewNotFoundError(fmt.Sprintf("update of %s with id %s", m.entity, id))
	}

	m.entities[id] = ve
	return nil
}
