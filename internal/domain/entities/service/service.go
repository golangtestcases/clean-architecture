package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/golangtestcases/clean-architecture/internal/domain/model"
	"github.com/google/uuid"
)

type EntityRepository interface {
	CreateEntity(_ context.Context, entity model.Entity) (model.Entity, error)
	GetEntitiesByID(_ context.Context, id model.EntityID) ([]model.Entity, error)
}

type EntityService struct {
	entityRepository EntityRepository
}

func NewEntityService(entityRepository EntityRepository) *EntityService {
	return &EntityService{entityRepository: entityRepository}
}

func (s *EntityService) CreateEntity(ctx context.Context, entity model.Entity) (model.Entity, error) {
	if entity.UserID == uuid.Nil {
		return model.Entity{}, errors.New("user_id must be provided")
	}

	newEntity, err := s.entityRepository.CreateEntity(ctx, entity)
	if err != nil {
		return model.Entity{}, fmt.Errorf("entityRepository.CreateEntity: %w", err)
	}

	return newEntity, nil
}

func (s *EntityService) GetEntitiesByID(ctx context.Context, id model.EntityID) ([]model.Entity, error) {
	if id < 1 {
		return nil, errors.New("id must be provided")
	}

	entities, err := s.entityRepository.GetEntitiesByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("entityRepository.GetEntitiesByID: %w", err)
	}

	return entities, nil
}
