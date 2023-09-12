package domain

import (
	"context"
	"git.miem.hse.ru/1206/app/errs"
	"git.miem.hse.ru/1206/material-library/internal/models/tplibrary"
)

func (d *Usecase) ListMaterialTypes(ctx context.Context) ([]*tplibrary.MaterialType, error) {
	mTypes, err := tplibrary.MaterialTypes().All(ctx, d.st.DBSX())
	if err != nil {
		return nil, errs.NewInternal(err)
	}
	return mTypes, nil
}
