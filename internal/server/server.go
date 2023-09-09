// Package server represents mostly used incoming port.
// There is http server with router, but it could be a grpc server, or tcp with custom protocol.
package server

import (
	"context"
	"microservice/internal/app/exibillia"
	"microservice/internal/app/looncan"
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
}

type Server struct {
	exibillia ExibilliaService
	looncan   LooncanService
	acaer     AcaerService

	v1.UnimplementedExibilliaServiceServer
	v1.UnimplementedLooncanServiceServer
	v1.UnimplementedAcaerServiceServer
}

func NewServer(exibillia ExibilliaService, looncan LooncanService, acaer AcaerService) *Server {
	return &Server{exibillia: exibillia, looncan: looncan, acaer: acaer}
}
