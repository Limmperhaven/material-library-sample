package domain

import (
	"context"
	"crypto/sha256"
	"git.miem.hse.ru/1206/app/errs"
	"git.miem.hse.ru/1206/app/logger"
	"git.miem.hse.ru/1206/app/typ"
	"git.miem.hse.ru/1206/material-library/internal/models/tplibrary"
	"github.com/friendsofgo/errors"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (d *Usecase) UploadMaterial(ctx context.Context, req *tplibrary.MaterialRequest) error {
	userId, err := typ.ExtractUserId(ctx)
	if err != nil {
		return errs.NewUnauthorized(nil)
	}
	canUpload, err := d.perms.CheckCanUploadMaterials(ctx, userId)
	if err != nil {
		return errs.NewInternal(err)
	} else if !canUpload {
		return errs.NewForbidden(nil)
	}
	valid, err := d.validateMaterialRequest(ctx, req)
	if err != nil {
		return errs.NewInternal(err)
	} else if !valid {
		return errs.NewBadRequest(nil)
	}
	if req.FileLink == nil {
		fileAlreadyExists, err := tplibrary.Materials(
			tplibrary.MaterialWhere.ChecksumSha256.EQ(null.BytesFrom(d.getFileHashSum(req.File))),
		).Exists(ctx, d.st.DBSX())
		if err != nil {
			return errs.NewInternal(err)
		} else if fileAlreadyExists {
			return errs.NewBadRequest(errors.New("file already exists"))
		}
	}
	var fileInfo *tplibrary.FileInfo
	err = d.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		material := tplibrary.Material{
			Title:              req.Title,
			SubjectID:          req.SubjectId,
			DifficultcyLevelID: req.DifficultcyLevelId,
			TypeID:             req.MaterialTypeId,
		}

		if req.FileLink != nil {
			material.URL = *req.FileLink
		} else {
			fileInfo, err = d.s3.PutFile(ctx, req.Title, req.File)
			if err != nil {
				return err
			}
			material.StorageKey = null.StringFrom(fileInfo.Key)
			material.ChecksumSha256 = null.BytesFrom(d.getFileHashSum(req.File))
			material.Size = null.Int64From(fileInfo.Size)
			material.URL = fileInfo.Url
		}

		err = material.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}
		err = d.perms.SetNewMaterialRelations(ctx, userId, material.ID)
		return err
	})
	if err != nil {
		if fileInfo != nil {
			err := d.s3.RemoveFile(ctx, fileInfo.Key)
			logger.Get(ctx).WithError(err).Error("error removing file from s3")
		}
		return errs.NewInternal(err)
	}
	return nil
}

func (d *Usecase) DeleteMaterial(ctx context.Context, materialId int64) error {
	userId, err := typ.ExtractUserId(ctx)
	if err != nil {
		return errs.NewUnauthorized(err)
	}
	canDelete, err := d.perms.CheckCanEditMaterial(ctx, userId, materialId)
	if err != nil {
		return errs.NewInternal(err)
	} else if !canDelete {
		return errs.NewForbidden(nil)
	}
	err = d.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		material := tplibrary.Material{ID: materialId}
		_, err := material.Delete(ctx, tx, false)
		if err != nil {
			return err
		}
		err = d.perms.DeleteMaterialRelations(ctx, materialId)
		return err
	})
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

func (d *Usecase) GetMaterial(ctx context.Context, materialId int64) (tplibrary.MaterialResponse, error) {
	userId, err := typ.ExtractUserId(ctx)
	if err != nil {
		return tplibrary.MaterialResponse{}, errs.NewUnauthorized(err)
	}
	material, err := tplibrary.Materials(
		tplibrary.MaterialWhere.ID.EQ(materialId),
		qm.Load(tplibrary.MaterialRels.Type),
	).One(ctx, d.st.DBSX())
	if err != nil {
		return tplibrary.MaterialResponse{}, errs.NewDbError(err)
	}
	canView, err := d.perms.CheckCanViewMaterial(ctx, userId, materialId)
	if err != nil {
		return tplibrary.MaterialResponse{}, errs.NewInternal(err)
	} else if !canView {
		return tplibrary.MaterialResponse{}, errs.NewForbidden(nil)
	}
	difLevel, err := d.edu.GetDifficultcyLevelById(ctx, material.DifficultcyLevelID)
	if err != nil {
		return tplibrary.MaterialResponse{}, errs.NewInternal(err)
	}
	subject, err := d.edu.GetSubjectById(ctx, material.SubjectID)
	if err != nil {
		return tplibrary.MaterialResponse{}, errs.NewInternal(err)
	}
	var size *int64
	if material.Size.Valid {
		size = &material.Size.Int64
	}
	response := tplibrary.MaterialResponse{
		Id:    material.ID,
		Title: material.Title,
		Size:  size,
		MaterialType: tplibrary.IdName{
			Id:   material.R.Type.ID,
			Name: material.R.Type.Name,
		},
		DifficultcyLevel: difLevel,
		Subject:          subject,
		FileLink:         material.URL,
	}
	return response, nil
}

func (d *Usecase) UpdateMaterial(ctx context.Context, req *tplibrary.MaterialRequest) error {
	userId, err := typ.ExtractUserId(ctx)
	if err != nil {
		return errs.NewUnauthorized(nil)
	}
	canEdit, err := d.perms.CheckCanEditMaterial(ctx, userId, req.Id)
	if err != nil {
		return errs.NewInternal(err)
	} else if !canEdit {
		return errs.NewForbidden(nil)
	}
	valid, err := d.validateMaterialRequest(ctx, req)
	if err != nil {
		return errs.NewInternal(err)
	} else if !valid {
		return errs.NewBadRequest(nil)
	}
	if req.FileLink == nil {
		fileAlreadyExists, err := tplibrary.Materials(
			tplibrary.MaterialWhere.ChecksumSha256.EQ(null.BytesFrom(d.getFileHashSum(req.File))),
			tplibrary.MaterialWhere.ID.NEQ(req.Id),
		).Exists(ctx, d.st.DBSX())
		if err != nil {
			return errs.NewInternal(err)
		} else if fileAlreadyExists {
			return errs.NewBadRequest(errors.New("file already exists"))
		}
	}
	var fileInfo *tplibrary.FileInfo
	err = d.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		material := tplibrary.Material{
			ID:                 req.Id,
			Title:              req.Title,
			SubjectID:          req.SubjectId,
			DifficultcyLevelID: req.DifficultcyLevelId,
			TypeID:             req.MaterialTypeId,
		}

		if req.FileLink != nil {
			material.URL = *req.FileLink
		} else {
			fileInfo, err = d.s3.PutFile(ctx, req.Title, req.File)
			if err != nil {
				return err
			}
			material.StorageKey = null.StringFrom(fileInfo.Key)
			material.ChecksumSha256 = null.BytesFrom(d.getFileHashSum(req.File))
			material.Size = null.Int64From(fileInfo.Size)
			material.URL = fileInfo.Url
		}

		_, err = material.Update(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		if fileInfo != nil {
			err := d.s3.RemoveFile(ctx, fileInfo.Key)
			logger.Get(ctx).WithError(err).Error("error removing file from s3")
		}
		return errs.NewInternal(err)
	}
	return nil
}

func (d *Usecase) ListMaterials(ctx context.Context) ([]*tplibrary.MaterialResponse, error) {
	userId, err := typ.ExtractUserId(ctx)
	if err != nil {
		return nil, errs.NewUnauthorized(err)
	}
	availableIds, err := d.perms.ListAvailableMaterialIds(ctx, userId)
	materials, err := tplibrary.Materials(
		tplibrary.MaterialWhere.ID.IN(availableIds),
		qm.Load(tplibrary.MaterialRels.Type),
	).All(ctx, d.st.DBSX())
	if err != nil {
		return nil, errs.NewInternal(err)
	}
	subIdToName, err := d.edu.GetSubjectIdToName(ctx)
	if err != nil {
		return nil, errs.NewInternal(err)
	}
	difLvlIdToName, err := d.edu.GetDifficultcyLevelIdToName(ctx)
	if err != nil {
		return nil, errs.NewInternal(err)
	}
	resp := make([]*tplibrary.MaterialResponse, len(materials))
	for i, m := range materials {
		var size *int64
		if m.Size.Valid {
			size = &m.Size.Int64
		}
		resp[i] = &tplibrary.MaterialResponse{
			Id:    m.ID,
			Title: m.Title,
			Size:  size,
			MaterialType: tplibrary.IdName{
				Id:   m.R.Type.ID,
				Name: m.R.Type.Name,
			},
			DifficultcyLevel: tplibrary.IdName{
				Id:   m.DifficultcyLevelID,
				Name: difLvlIdToName[m.DifficultcyLevelID],
			},
			Subject: tplibrary.IdName{
				Id:   m.SubjectID,
				Name: subIdToName[m.SubjectID],
			},
			FileLink: m.URL,
		}
	}
	return resp, nil
}

func (d *Usecase) getFileHashSum(file []byte) []byte {
	hash := sha256.Sum256(file)
	return hash[:]
}

func (d *Usecase) validateMaterialRequest(ctx context.Context, req *tplibrary.MaterialRequest) (bool, error) {
	if req.FileLink == nil && req.File == nil {
		return false, nil
	}
	valid, err := d.edu.ValidateSubjectId(ctx, req.SubjectId)
	if err != nil {
		return false, err
	} else if !valid {
		return false, nil
	}
	valid, err = d.edu.ValidateDifficultcyLevelId(ctx, req.DifficultcyLevelId)
	if err != nil {
		return false, err
	} else if !valid {
		return false, nil
	}
	return true, nil
}
