// Package acaer provides creation of entity with related one.
// It shows:
// - ways to work with related entities from separated packages;
// - two-layered validation.
package acaer

import (
	"context"
	"go.uber.org/zap"
	"microservice/internal/app/looncan"
	"microservice/internal/utils"
)

type Storage interface {
	begin(ctx context.Context) (Storage, error)
	commit() error
	rollback() error
	getVersions(ctx context.Context) ([]string, error)
	create(ctx context.Context, acaer Acaer) (uint64, error)
	createLooncans(ctx context.Context, list []looncanDTO) error
	createFromAggregate(ctx context.Context, acaer Acaer) error
}

type LooncanService interface {
	Create(ctx context.Context, entities []looncan.Looncan) error
}

type Service struct {
	storage Storage
	looncan LooncanService
}

func NewService(storage Storage, looncan LooncanService) *Service {
	return &Service{storage: storage, looncan: looncan}
}

// CreateSimple uses method from related package. It's simple DI, but we cannot control any transaction here.
// Additionally, there is implemented service-layer complex validation, based on data from storage.
func (s *Service) CreateSimple(ctx context.Context, name, version string) error {
	entity := Acaer{
		Name:    name,
		Version: version,
	}

	versions, err := s.storage.getVersions(ctx)
	if err != nil {
		zap.S().Errorw("get allowed versions", "error", err)

		return err
	}

	validator := NewValidator(versions)
	err = validator.Validate(entity)
	if err != nil {
		return err
	}

	id, err := s.storage.create(ctx, entity)
	if err != nil {
		zap.S().Errorw("storage create", "error", err)

		return err
	}

	looncans := make([]looncan.Looncan, 0, 3)
	for i := 0; i < 3; i++ {
		looncans = append(looncans, looncan.Looncan{
			Name:       utils.RandomString(10),
			Value:      utils.RandomString(6),
			ParentID:   id,
			ParentType: looncan.ParentTypeAcaer,
		})
	}

	err = s.looncan.Create(ctx, looncans)
	if err != nil {
		zap.S().Errorw("looncan service create", "error", err)
	}

	return err
}

// CreateTransaction uses internal DTO instead of imported entity to work with storage.
// Transaction control implemented in storage, so we can use any storage methods in transaction without leaking Tx into
// business code. Storage works with whole database structure, makes queries to any table, not only one for the main entity.
func (s *Service) CreateTransaction(ctx context.Context, name, version string) error {
	st, err := s.storage.begin(ctx)
	if err != nil {
		zap.S().Errorw("begin storage transaction", "error", err)

		return err
	}
	defer func(st Storage) {
		err := st.rollback()
		if err != nil {
			zap.S().Errorw("rollback storage transaction", "error", err)
		}
	}(st)

	entity := Acaer{
		Name:    name,
		Version: version,
	}

	id, err := st.create(ctx, entity)
	if err != nil {
		zap.S().Errorw("storage create acaer", "error", err)

		return err
	}

	looncans := make([]looncanDTO, 0, 3)
	for i := 0; i < 3; i++ {
		looncans = append(looncans, looncanDTO{
			name:       utils.RandomString(10),
			value:      utils.RandomString(6),
			parentID:   id,
			parentType: parentAcaer,
		})
	}

	err = st.createLooncans(ctx, looncans)
	if err != nil {
		zap.S().Errorw("storage create looncans", "error", err)

		return err
	}

	err = st.commit()
	if err != nil {
		zap.S().Errorw("commit storage transaction", "error", err)
	}

	return err
}

// CreateAggregate uses DDD-like root aggregation and just one method of storage.
// Transaction implemented inside the method. Storage also has access to any table.
func (s *Service) CreateAggregate(ctx context.Context, name, version string) error {
	entity := Acaer{
		Name:    name,
		Version: version,
	}

	entity.Entsian()

	err := s.storage.createFromAggregate(ctx, entity)
	if err != nil {
		zap.S().Errorw("storage create from aggregate", "error", err)
	}

	return err
}
