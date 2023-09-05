// Package exibillia provides CRUD operations for one simple entity.
// Also, it shows how to:
// - scaffold basic service with DI via constructor
// - work with errors (logging, typed errors)
// - return updated entity for API
package exibillia

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	errt "microservice/internal/errors"
)

//go:generate mockgen -source=./service.go -package=exibillia -destination=service_mock.go

type Storage interface {
	create(ctx context.Context, exibillia Exibillia) (uint64, error)
	getByID(ctx context.Context, id uint64) (Exibillia, error)
	update(ctx context.Context, exibillia *Exibillia) error
	delete(ctx context.Context, id uint64) error
}
type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Create(ctx context.Context, exibillia Exibillia) (uint64, error) {
	id, err := s.storage.create(ctx, exibillia)
	if err != nil {
		zap.S().Errorw("storage create", "error", err)

		return 0, err
	}

	return id, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (Exibillia, error) {
	ex, err := s.storage.getByID(ctx, id)
	if errors.Is(err, errNoRows) {
		return Exibillia{}, errt.NewNotFoundError(fmt.Sprintf("exibillia with id %d", id))
	}

	if err != nil {
		zap.S().Errorw("get from storage", "error", err, "id", id)

		return Exibillia{}, err
	}

	return ex, nil
}

func (s *Service) Update(ctx context.Context, req UpdateRequest) error {
	ex, err := s.storage.getByID(ctx, req.ID)
	if errors.Is(err, errNoRows) {
		return errt.NewNotFoundError(fmt.Sprintf("exibillia with id %d", req.ID))
	}
	if err != nil {
		zap.S().Errorw("get from storage", "error", err, "id", req.ID)

		return err
	}

	ex.Description = req.Description
	ex.Tags = req.Tags

	if err = s.storage.update(ctx, &ex); err != nil {
		zap.S().Errorw("storage update", "error", err)

		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	err := s.storage.delete(ctx, id)
	if errors.Is(err, errNoRows) {
		return errt.NewNotFoundError(fmt.Sprintf("exibillia with id %d", id))
	}
	if err != nil {
		zap.S().Errorw("storage delete", "error", err)

		return err
	}

	return nil
}
