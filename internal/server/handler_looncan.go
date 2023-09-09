package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"microservice/internal/app/looncan"
	v1 "microservice/pkg/v1"
)

func (s *Server) List(ctx context.Context, _ *empty.Empty) (*v1.ListLooncanResponse, error) {
	list, err := s.looncan.GetAll(ctx)
	if err != nil {
		return &v1.ListLooncanResponse{}, status.Error(codes.Internal, err.Error())
	}

	res := make([]*v1.Looncan, 0, len(list))
	for _, l := range list {
		res = append(res, mapLooncan2Pb(l))
	}

	return &v1.ListLooncanResponse{Items: res}, nil
}

func (s *Server) ListForParent(ctx context.Context, req *v1.ListLooncanForParentRequest) (*v1.ListLooncanResponse, error) {
	if req.GetParentId() == 0 {
		return &v1.ListLooncanResponse{}, status.Error(codes.FailedPrecondition, "invalid parent_id")
	}
	if req.GetParentType() == v1.ListLooncanForParentRequest_ParentUnspecified {
		return &v1.ListLooncanResponse{}, status.Error(codes.FailedPrecondition, "parent_type must be specified")
	}

	var parentType looncan.ParentType
	switch req.GetParentType() {
	case v1.ListLooncanForParentRequest_ParentAcaer:
		parentType = looncan.ParentTypeAcaer
	case v1.ListLooncanForParentRequest_ParentExibillia:
		parentType = looncan.ParentTypeExibillia
	}

	list, err := s.looncan.GetByParent(ctx, req.GetParentId(), parentType)
	if err != nil {
		return &v1.ListLooncanResponse{}, status.Error(codes.Internal, err.Error())
	}

	res := make([]*v1.Looncan, 0, len(list))
	for _, l := range list {
		res = append(res, mapLooncan2Pb(l))
	}

	return &v1.ListLooncanResponse{Items: res}, nil
}

func mapLooncan2Pb(l looncan.Looncan) *v1.Looncan {
	var updatedAt *timestamppb.Timestamp
	if l.UpdatedAt != nil {
		updatedAt = timestamppb.New(*l.UpdatedAt)
	}
	return &v1.Looncan{
		Id:        l.ID,
		Name:      l.Name,
		Value:     l.Value,
		CreatedAt: timestamppb.New(l.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
