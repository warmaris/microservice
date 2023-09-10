package server

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"microservice/internal/errors"
	v1 "microservice/pkg/v1"
)

func (s *Server) CreateSimple(ctx context.Context, req *v1.CreateAcaerRequest) (*emptypb.Empty, error) {
	err := validateCreateAcaerRequest(req)
	if err != nil {
		return &emptypb.Empty{}, s.handleError(err)
	}

	err = s.acaer.CreateSimple(ctx, req.GetName(), req.GetVersion())
	return &emptypb.Empty{}, s.handleError(err)
}

func (s *Server) CreateTransaction(ctx context.Context, req *v1.CreateAcaerRequest) (*emptypb.Empty, error) {
	err := validateCreateAcaerRequest(req)
	if err != nil {
		return &emptypb.Empty{}, s.handleError(err)
	}

	err = s.acaer.CreateTransaction(ctx, req.GetName(), req.GetVersion())
	return &emptypb.Empty{}, s.handleError(err)
}

func (s *Server) CreateAggregate(ctx context.Context, req *v1.CreateAcaerRequest) (*emptypb.Empty, error) {
	err := validateCreateAcaerRequest(req)
	if err != nil {
		return &emptypb.Empty{}, s.handleError(err)
	}

	err = s.acaer.CreateAggregate(ctx, req.GetName(), req.GetVersion())
	return &emptypb.Empty{}, s.handleError(err)
}

// validateCreateAcaerRequest provides infra-layer validation for incoming request, part of two-layered validation.
func validateCreateAcaerRequest(req *v1.CreateAcaerRequest) error {
	var validationErr errors.ValidationError
	if req.GetName() == "" {
		validationErr.Add("name", "must not be empty")
	}
	if req.GetVersion() == "" {
		validationErr.Add("version", "must not be empty")
	}

	if validationErr.HasErrors() {
		return validationErr
	}
	return nil
}
