// Package looncan provides sub-related entity and some methods:
// - common methods for usage from related packages
// - own methods for request handling
package looncan

import (
	"context"
	"go.uber.org/zap"
)

type Storage interface {
	getAllForParent(ctx context.Context, parentID uint64, parentType ParentType) ([]Looncan, error)
	list(ctx context.Context) ([]Looncan, error)
	create(ctx context.Context, entities []Looncan) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Create(ctx context.Context, entities []Looncan) error {
	err := s.storage.create(ctx, entities)
	if err != nil {
		zap.S().Errorw("storage create", "error", err)
	}

	return err
}

func (s *Service) GetByParent(ctx context.Context, parentID uint64, parentType ParentType) ([]Looncan, error) {
	list, err := s.storage.getAllForParent(ctx, parentID, parentType)
	if err != nil {
		zap.S().Errorw("storage get for parent", "error", err, "parent_id", parentID, "parent_type", parentType)
	}

	return list, err
}

func (s *Service) GetAll(ctx context.Context) ([]Looncan, error) {
	list, err := s.storage.list(ctx)
	if err != nil {
		zap.S().Errorw("storage list", "error", err)
	}

	return list, err
}
