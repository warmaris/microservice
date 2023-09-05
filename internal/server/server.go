// Package server represents mostly used incoming port.
// There is http server with router, but it could be a grpc server, or tcp with custom protocol.
package server

import (
	"context"
	"microservice/internal/app/exibillia"
	v1 "microservice/pkg/v1"
)

type ExibilliaService interface {
	Create(ctx context.Context, exibillia exibillia.Exibillia) (uint64, error)
	GetByID(ctx context.Context, id uint64) (exibillia.Exibillia, error)
	Update(ctx context.Context, req exibillia.UpdateRequest) error
	Delete(ctx context.Context, id uint64) error
}

type Server struct {
	exibillia ExibilliaService

	v1.UnimplementedExibilliaServiceServer
}

func NewServer(exibillia ExibilliaService) *Server {
	return &Server{exibillia: exibillia}
}
