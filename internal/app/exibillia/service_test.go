package exibillia

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ServiceTestSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	storage *MockStorage
	service *Service
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupSuite() {
	s.ctrl = gomock.NewController(s.T())
	s.storage = NewMockStorage(s.ctrl)
	s.service = NewService(s.storage)
}

func (s *ServiceTestSuite) TestCreate() {
	entity := Exibillia{
		Name:        "Test Name",
		Description: "Test Description",
		Tags:        []string{"test", "create"},
	}
	testErr := errors.New("test error")

	s.Run("success", func() {
		s.storage.EXPECT().create(gomock.Any(), entity).Return(nil)

		assert.NoError(s.T(), s.service.Create(context.Background(), entity))
	})

	s.Run("storage error - return it", func() {
		s.storage.EXPECT().create(gomock.Any(), entity).Return(testErr)

		assert.ErrorIs(s.T(), s.service.Create(context.Background(), entity), testErr)
	})
}

func (s *ServiceTestSuite) TestGetByID() {
	id := uint64(123)
	entity := Exibillia{
		ID:          123,
		Name:        "Test Name",
		Description: "Test Description",
		Tags:        []string{"test", "get_by"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	testErr := errors.New("test error")

	s.Run("success", func() {
		s.storage.EXPECT().getByID(gomock.Any(), id).Return(entity, nil)

		res, err := s.service.GetByID(context.Background(), id)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), entity, res)
	})

	s.Run("storage error - return it", func() {
		s.storage.EXPECT().getByID(gomock.Any(), id).Return(Exibillia{}, testErr)

		res, err := s.service.GetByID(context.Background(), id)

		assert.ErrorIs(s.T(), err, testErr)
		assert.Empty(s.T(), res)
	})
}

func (s *ServiceTestSuite) TestUpdate() {
	req := UpdateRequest{
		ID:          123,
		Description: "New Description",
		Tags:        []string{"old tag", "new tag"},
	}
	entityOld := Exibillia{
		ID:          123,
		Name:        "Test Name",
		Description: "Test Description",
		Tags:        []string{"old tag"},
		CreatedAt:   time.Date(2023, 2, 7, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, 2, 7, 0, 0, 0, 0, time.UTC),
	}
	entityNew := Exibillia{
		ID:          123,
		Name:        "Test Name",
		Description: "New Description",
		Tags:        []string{"old tag", "new tag"},
		CreatedAt:   time.Date(2023, 2, 7, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, 2, 7, 0, 0, 0, 0, time.UTC),
	}
	testErr := errors.New("test err")

	s.Run("success", func() {
		s.storage.EXPECT().getByID(gomock.Any(), uint64(123)).Return(entityOld, nil)
		s.storage.EXPECT().update(gomock.Any(), gomock.Eq(&entityNew)).Return(nil)

		res, err := s.service.Update(context.Background(), req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), entityNew, res)
	})

	s.Run("update error - return it", func() {
		s.storage.EXPECT().getByID(gomock.Any(), uint64(123)).Return(entityOld, nil)
		s.storage.EXPECT().update(gomock.Any(), gomock.Eq(&entityNew)).Return(testErr)

		res, err := s.service.Update(context.Background(), req)
		assert.ErrorIs(s.T(), err, testErr)
		assert.Empty(s.T(), res)
	})

	s.Run("get by id error - return it", func() {
		s.storage.EXPECT().getByID(gomock.Any(), uint64(123)).Return(Exibillia{}, testErr)

		res, err := s.service.Update(context.Background(), req)
		assert.ErrorIs(s.T(), err, testErr)
		assert.Empty(s.T(), res)
	})
}

func (s *ServiceTestSuite) TestDelete() {
	id := uint64(123)
	testErr := errors.New("test error")

	s.Run("success", func() {
		s.storage.EXPECT().delete(gomock.Any(), id).Return(nil)

		assert.NoError(s.T(), s.service.Delete(context.Background(), id))
	})

	s.Run("storage error - return it", func() {
		s.storage.EXPECT().delete(gomock.Any(), id).Return(testErr)

		assert.ErrorIs(s.T(), s.service.Delete(context.Background(), id), testErr)
	})
}
