package service

import (
	"context"
	"git.miem.hse.ru/1206/material-library/internal/models/mapper"
	"git.miem.hse.ru/1206/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *LibraryServer) UploadMaterial(ctx context.Context, request *pb.MaterialRequest) (*emptypb.Empty, error) {
	in := mapper.NewMaterialRequestFromPb(request)
	err := s.uc.UploadMaterial(ctx, in)
	return &emptypb.Empty{}, err
}

func (s *LibraryServer) DeleteMaterial(ctx context.Context, id *pb.Id) (*emptypb.Empty, error) {
	in := mapper.NewIdFromPb(id)
	err := s.uc.DeleteMaterial(ctx, in)
	return &emptypb.Empty{}, err
}

func (s *LibraryServer) GetMaterial(ctx context.Context, id *pb.Id) (*pb.MaterialResponse, error) {
	in := mapper.NewIdFromPb(id)
	out, err := s.uc.GetMaterial(ctx, in)
	return mapper.NewMaterialResponseToPb(&out), err
}

func (s *LibraryServer) UpdateMaterial(ctx context.Context, request *pb.MaterialRequest) (*emptypb.Empty, error) {
	in := mapper.NewMaterialRequestFromPb(request)
	err := s.uc.UpdateMaterial(ctx, in)
	return &emptypb.Empty{}, err
}

func (s *LibraryServer) ListMaterials(ctx context.Context, _ *emptypb.Empty) (*pb.ListMaterialsResponse, error) {
	out, err := s.uc.ListMaterials(ctx)
	return mapper.NewListMaterialsResponseToPb(out), err
}

func (s *LibraryServer) ListMaterialTypes(ctx context.Context, _ *emptypb.Empty) (*pb.MaterialTypeArray, error) {
	out, err := s.uc.ListMaterialTypes(ctx)
	return mapper.NewListMaterialTypesFromPb(out), err
}
