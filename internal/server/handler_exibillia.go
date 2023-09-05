package server

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"microservice/internal/app/exibillia"
	errt "microservice/internal/errors"
	v1 "microservice/pkg/v1"
)

func (s *Server) Get(ctx context.Context, req *v1.GetExibilliaRequest) (*v1.GetExibilliaResponse, error) {
	obj, err := s.exibillia.GetByID(ctx, req.GetId())
	if errors.As(err, &errt.NotFoundError{}) {
		return &v1.GetExibilliaResponse{}, status.Error(codes.NotFound, "not found")
	}
	if err != nil {
		return &v1.GetExibilliaResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &v1.GetExibilliaResponse{Exibillia: &v1.Exibillia{
		Id:          obj.ID,
		Name:        obj.Name,
		Description: obj.Description,
		Tags:        obj.Tags,
		CreatedAt:   timestamppb.New(obj.CreatedAt),
		UpdatedAt:   timestamppb.New(obj.UpdatedAt),
	}}, nil
}

func (s *Server) Create(ctx context.Context, req *v1.CreateExibilliaRequest) (*v1.CreateExibilliaResponse, error) {
	model := exibillia.Exibillia{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Tags:        req.GetTags(),
	}

	id, err := s.exibillia.Create(ctx, model)
	if err != nil {
		return &v1.CreateExibilliaResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &v1.CreateExibilliaResponse{Id: id}, nil
}

func (s *Server) Update(ctx context.Context, req *v1.UpdateExibilliaRequest) (*emptypb.Empty, error) {
	updateReq := exibillia.UpdateRequest{
		ID:          req.GetId(),
		Description: req.GetDescription(),
		Tags:        req.GetTags(),
	}
	err := s.exibillia.Update(ctx, updateReq)
	if errors.As(err, &errt.NotFoundError{}) {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "not found")
	}
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, req *v1.DeleteExibilliaRequest) (*emptypb.Empty, error) {
	err := s.exibillia.Delete(ctx, req.GetId())
	if errors.As(err, &errt.NotFoundError{}) {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "not found")
	}
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
