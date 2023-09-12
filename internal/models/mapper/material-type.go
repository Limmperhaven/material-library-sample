package mapper

import (
	"git.miem.hse.ru/1206/material-library/internal/models/tplibrary"
	"git.miem.hse.ru/1206/proto/pb"
)

func NewListMaterialTypesFromPb(in []*tplibrary.MaterialType) *pb.MaterialTypeArray {
	items := make([]*pb.MaterialType, len(in))
	for i, mt := range in {
		items[i] = &pb.MaterialType{
			Id:   mt.ID,
			Name: mt.Name,
			Url:  mt.ImageURL,
		}
	}
	return &pb.MaterialTypeArray{Items: items}
}
