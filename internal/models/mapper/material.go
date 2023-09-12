package mapper

import (
	"git.miem.hse.ru/1206/material-library/internal/models/tplibrary"
	"git.miem.hse.ru/1206/proto/pb"
)

func NewIdFromPb(in *pb.Id) int64 {
	return in.Id
}

func NewIdNameToPb(in *tplibrary.IdName) *pb.IdName {
	return &pb.IdName{
		Id:   in.Id,
		Name: in.Name,
	}
}

func NewMaterialRequestFromPb(in *pb.MaterialRequest) *tplibrary.MaterialRequest {
	return &tplibrary.MaterialRequest{
		Id:                 in.Id,
		Title:              in.Title,
		SubjectId:          in.SubjectId,
		DifficultcyLevelId: in.DifficultcyLevelId,
		MaterialTypeId:     in.MaterialTypeId,
		FileLink:           in.FileLink,
		File:               in.File,
	}
}

func NewMaterialResponseToPb(in *tplibrary.MaterialResponse) *pb.MaterialResponse {
	return &pb.MaterialResponse{
		Id:               in.Id,
		Title:            in.Title,
		Size:             in.Size,
		MaterialType:     NewIdNameToPb(&in.MaterialType),
		DifficultcyLevel: NewIdNameToPb(&in.DifficultcyLevel),
		Subject:          NewIdNameToPb(&in.Subject),
		FileLink:         in.FileLink,
	}
}

func NewListMaterialsResponseToPb(in []*tplibrary.MaterialResponse) *pb.ListMaterialsResponse {
	items := make([]*pb.MaterialResponse, len(in))
	for i, mr := range in {
		items[i] = &pb.MaterialResponse{
			Id:               mr.Id,
			Title:            mr.Title,
			Size:             mr.Size,
			MaterialType:     NewIdNameToPb(&mr.MaterialType),
			DifficultcyLevel: NewIdNameToPb(&mr.DifficultcyLevel),
			Subject:          NewIdNameToPb(&mr.Subject),
			FileLink:         mr.FileLink,
		}
	}
	return &pb.ListMaterialsResponse{Items: items}
}
