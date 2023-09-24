// Package server represents mostly used incoming port.
// There is http server with router, but it could be a grpc server, or tcp with custom protocol.
package server

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservice/internal/app/exibillia"
	"microservice/internal/app/looncan"
	errt "microservice/internal/errors"
	v1 "microservice/pkg/v1"
)

type ExibilliaService interface {
	Create(ctx context.Context, exibillia exibillia.Exibillia) (uint64, error)
	GetByID(ctx context.Context, id uint64) (exibillia.Exibillia, error)
	Update(ctx context.Context, req exibillia.UpdateRequest) error
	Delete(ctx context.Context, id uint64) error
}

type LooncanService interface {
	GetByParent(ctx context.Context, parentID uint64, parentType looncan.ParentType) ([]looncan.Looncan, error)
	GetAll(ctx context.Context) ([]looncan.Looncan, error)
}

type AcaerService interface {
	CreateSimple(ctx context.Context, name, version string) error
	CreateTransaction(ctx context.Context, name, version string) error
	CreateAggregate(ctx context.Context, name, version string) error
}

type JarklinService interface {
	CreateAndSend(ctx context.Context, name string) error
	CreateAndSave(ctx context.Context, name string) error
}

type Server struct {
	exibillia ExibilliaService
	looncan   LooncanService
	acaer     AcaerService
	jarklin JarklinService

	v1.UnimplementedExibilliaServiceServer
	v1.UnimplementedLooncanServiceServer
	v1.UnimplementedAcaerServiceServer
	v1.UnimplementedJarklinServiceServer
}

func NewServer(exibillia ExibilliaService, looncan LooncanService, acaer AcaerService, jarklin JarklinService) *Server {
	return &Server{exibillia: exibillia, looncan: looncan, acaer: acaer, jarklin: jarklin}
}

func (s *Server) handleError(err error) error {
	if errors.As(err, &errt.NotFoundError{}) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.As(err, &errt.ValidationError{}) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}
