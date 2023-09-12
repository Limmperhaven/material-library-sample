package domain

import (
	"context"
	"git.miem.hse.ru/1206/material-library/internal/models/tplibrary"
)

type Permissions interface {
	CheckCanUploadMaterials(ctx context.Context, userId int64) (bool, error)
	CheckCanEditMaterial(ctx context.Context, userId int64, materialId int64) (bool, error)
	CheckCanViewMaterial(ctx context.Context, userId int64, materialId int64) (bool, error)
	SetNewMaterialRelations(ctx context.Context, ownerId int64, materialId int64) error
	ListAvailableMaterialIds(ctx context.Context, userId int64) ([]int64, error)
	DeleteMaterialRelations(ctx context.Context, materialId int64) error
}

type S3 interface {
	PutFile(ctx context.Context, fName string, body []byte) (info *tplibrary.FileInfo, err error)
	RemoveFile(ctx context.Context, key string) error
}

type Education interface {
	ValidateSubjectId(ctx context.Context, subjectId int64) (bool, error)
	ValidateDifficultcyLevelId(ctx context.Context, difficultcyLevelId int64) (bool, error)
	GetSubjectById(ctx context.Context, subjectId int64) (tplibrary.IdName, error)
	GetDifficultcyLevelById(ctx context.Context, difficultcyLevelId int64) (tplibrary.IdName, error)
	GetSubjectIdToName(ctx context.Context) (map[int64]string, error)
	GetDifficultcyLevelIdToName(ctx context.Context) (map[int64]string, error)
}
