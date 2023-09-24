package server

import (
	"context"
	"microservice/internal/errors"
	v1 "microservice/pkg/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateAndSend(ctx context.Context, req *v1.CreateJarklinRequest) (*emptypb.Empty, error) {
	if req.GetName() == "" {
		return &emptypb.Empty{}, s.handleError(errors.NewValidationError("name", "must not be empty"))
	}

	err := s.jarklin.CreateAndSend(ctx, req.GetName())
	return &emptypb.Empty{}, s.handleError(err)
}

func (s *Server) CreateAndSave(ctx context.Context, req *v1.CreateJarklinRequest) (*emptypb.Empty, error) {
	if req.GetName() == "" {
		return &emptypb.Empty{}, s.handleError(errors.NewValidationError("name", "must not be empty"))
	}

	err := s.jarklin.CreateAndSave(ctx, req.GetName())
	return &emptypb.Empty{}, s.handleError(err)
}