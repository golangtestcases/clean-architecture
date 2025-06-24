package repository

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/golangtestcases/clean-architecture/internal/domain/model"
)

type InMemoryRepository struct {
	storage map[model.EntityID][]model.Entity
	mx      sync.RWMutex

	idFactory atomic.Uint64
}

func NewInMemoryRepository(cap int) *InMemoryRepository {
	return &InMemoryRepository{
		storage: make(map[model.EntityID][]model.Entity, cap),
	}
}

func (r *InMemoryRepository) CreateEntity(_ context.Context, entity model.Entity) (model.Entity, error) {
	entityID := r.idFactory.Add(1)
	entity.ID = model.EntityID(entityID)

	r.mx.Lock()
	defer r.mx.Unlock()
	r.storage[entity.ID] = append(r.storage[entity.ID], entity)

	return entity, nil
}

func (r *InMemoryRepository) GetEntitiesByID(_ context.Context, id model.EntityID) ([]model.Entity, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	return r.storage[id], nil
}
