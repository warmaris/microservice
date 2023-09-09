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

func (s *Service) CreateSimple(ctx context.Context, name, version string) error {
	// todo: add version validation

	entity := Acaer{
		Name:    name,
		Version: version,
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
